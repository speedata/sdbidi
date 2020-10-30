package sdbidi

import (
	"bytes"
	"fmt"
)

// I took the file bidi.go from the above and tried to fill out the "unimplemented" methods.

// A Direction indicates the overall flow of text.
type Direction int

const (
	// LeftToRight indicates the text contains no right-to-left characters and
	// that either there are some left-to-right characters or the option
	// DefaultDirection(LeftToRight) was passed.
	LeftToRight Direction = iota

	// RightToLeft indicates the text contains no left-to-right characters and
	// that either there are some right-to-left characters or the option
	// DefaultDirection(RightToLeft) was passed.
	RightToLeft

	// Mixed indicates text contains both left-to-right and right-to-left
	// characters.
	Mixed

	// Neutral means that text contains no left-to-right and right-to-left
	// characters and that no default direction has been set.
	Neutral
)

type options struct{}

// An Option is an option for Bidi processing.
type Option func(*options)

// ICU allows the user to define embedding levels. This may be used, for example,
// to use hierarchical structure of markup languages to define embeddings.
// The following option may be a way to expose this functionality in this API.
// // LevelFunc sets a function that associates nesting levels with the given text.
// // The levels function will be called with monotonically increasing values for p.
// func LevelFunc(levels func(p int) int) Option {
// 	panic("unimplemented")
// }

// DefaultDirection sets the default direction for a Paragraph. The direction is
// overridden if the text contains directional characters.
func DefaultDirection(d Direction) Option {
	panic("unimplemented")
}

// A Paragraph holds a single Paragraph for Bidi processing.
type Paragraph struct {
	p []byte
	o Ordering
	// buffers
}

// SetBytes configures p for the given paragraph text. It replaces text
// previously set by SetBytes or SetString. If b contains a paragraph separator
// it will only process the first paragraph and report the number of bytes
// consumed from b including this separator. Error may be non-nil if options are
// given.
func (p *Paragraph) SetBytes(b []byte, opts ...Option) (n int, err error) {
	p.p = b
	return len(b), nil
}

// SetString configures p for the given paragraph text. It replaces text
// previously set by SetBytes or SetString. If s contains a paragraph separator
// it will only process the first paragraph and report the number of bytes
// consumed from s including this separator. Error may be non-nil if options are
// given.
func (p *Paragraph) SetString(s string, opts ...Option) (n int, err error) {
	p.p = []byte(s)
	return len(p.p), nil
}

// IsLeftToRight reports whether the principle direction of rendering for this
// paragraphs is left-to-right. If this returns false, the principle direction
// of rendering is right-to-left.
func (p *Paragraph) IsLeftToRight() bool {
	return p.Direction() == LeftToRight
}

// Direction returns the direction of the text of this paragraph.
//
// The direction may be LeftToRight, RightToLeft, Mixed, or Neutral.
func (p *Paragraph) Direction() Direction {
	return p.o.Direction()
}

// RunAt reports the Run at the given position of the input text.
//
// This method can be used for computing line breaks on paragraphs.
func (p *Paragraph) RunAt(pos int) Run {
	panic("unimplemented")
}

// Order computes the visual ordering of all the runs in a Paragraph.
func (p *Paragraph) Order() (Ordering, error) {
	types := []Class{}
	pairTypes := []bracketType{}
	pairValues := []rune{}
	runes := bytes.Runes(p.p)

	for _, r := range runes {
		props, _ := LookupRune(r)
		types = append(types, props.Class())
		switch r {
		case '(', '[', '{':
			pairTypes = append(pairTypes, bpOpen)
			pairValues = append(pairValues, r)
		case ')', ']', '}':
			pairTypes = append(pairTypes, bpClose)
			pairValues = append(pairValues, r)
		default:
			pairTypes = append(pairTypes, bpNone)
			pairValues = append(pairValues, 0)
		}
	}

	if len(types) == 0 {
		return Ordering{}, fmt.Errorf("Cannot order empty string")
	}
	para := newParagraph(types, pairTypes, pairValues, -1)
	levels := para.getLevels([]int{len(types)})

	var curDir Direction

	prevDir := Neutral
	prevI := 0

	o := Ordering{}
	for i, lvl := range levels {
		if lvl%2 == 0 {
			curDir = LeftToRight
		} else {
			curDir = RightToLeft
		}
		if curDir != prevDir {
			if i > 0 {
				o.runes = append(o.runes, runes[prevI:i])
				o.directions = append(o.directions, prevDir)
			}
			prevI = i
			prevDir = curDir
		}
	}
	o.runes = append(o.runes, runes[prevI:])
	o.directions = append(o.directions, prevDir)
	p.o = o
	return o, nil
}

// Line computes the visual ordering of runs for a single line starting and
// ending at the given positions in the original text.
func (p *Paragraph) Line(start, end int) (Ordering, error) {
	panic("unimplemented")
}

// An Ordering holds the computed visual order of runs of a Paragraph. Calling
// SetBytes or SetString on the originating Paragraph invalidates an Ordering.
// The methods of an Ordering should only be called by one goroutine at a time.
type Ordering struct {
	runes      [][]rune
	directions []Direction
}

// Direction reports the directionality of the runs.
//
// The direction may be LeftToRight, RightToLeft, Mixed, or Neutral.
func (o *Ordering) Direction() Direction {
	return o.directions[0]
}

// NumRuns returns the number of runs.
func (o *Ordering) NumRuns() int {
	return len(o.runes)
}

// Run returns the ith run within the ordering.
func (o *Ordering) Run(i int) Run {
	r := Run{
		runes:     o.runes[i],
		direction: o.directions[i],
	}
	return r
}

// TODO: perhaps with options.
// // Reorder creates a reader that reads the runes in visual order per character.
// // Modifiers remain after the runes they modify.
// func (l *Runs) Reorder() io.Reader {
// 	panic("unimplemented")
// }

// A Run is a continuous sequence of characters of a single direction.
type Run struct {
	runes     []rune
	direction Direction
}

// String returns the text of the run in its original order.
func (r *Run) String() string {
	return string(r.runes)
}

// Bytes returns the text of the run in its original order.
func (r *Run) Bytes() []byte {
	return []byte(r.String())
}

// TODO: methods for
// - Display order
// - headers and footers
// - bracket replacement.

// Direction reports the direction of the run.
func (r *Run) Direction() Direction {
	return r.direction
}

// Position of the Run within the text passed to SetBytes or SetString of the
// originating Paragraph value.
func (r *Run) Pos() (start, end int) {
	panic("unimplemented")
}

// AppendReverse reverses the order of characters of in, appends them to out,
// and returns the result. Modifiers will still follow the runes they modify.
// Brackets are replaced with their counterparts.
func AppendReverse(out, in []byte) []byte {
	panic("unimplemented")
}

// ReverseString reverses the order of characters in s and returns a new string.
// Modifiers will still follow the runes they modify. Brackets are replaced with
// their counterparts.
func ReverseString(s string) string {
	panic("unimplemented")
}
