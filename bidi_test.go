package sdbidi

import (
	"log"
	"testing"
)

type runInformation struct {
	str   string
	dir   Direction
	start int
	end   int
}

func TestSimple(t *testing.T) {
	str := "Hellö"
	p := Paragraph{}
	p.SetString(str)
	order, err := p.Order()
	if err != nil {
		log.Fatal(err)
	}
	expectedRuns := []runInformation{
		{"Hellö", LeftToRight, 0, 4},
	}

	if !p.IsLeftToRight() {
		t.Error("Paragraph should return LeftToRight() == true")
	}
	if nr, expected := order.NumRuns(), len(expectedRuns); nr != expected {
		t.Errorf("Number of runs must be %d but got %d", expected, nr)
	}
	for i, er := range expectedRuns {
		r := order.Run(i)
		if str := r.String(); str != er.str {
			t.Errorf("Run %d should have string %q but has %q", i, er.str, str)
		}
		if s, e := r.Pos(); s != er.start || e != er.end {
			t.Errorf("Run %d should go from %d to %d but got %d to %d", i, er.start, er.end, s, e)
		}
		if d := r.Direction(); d != er.dir {
			t.Errorf("Run %d direction should be %d but got %d", i, er.dir, d)
		}
	}
}

func TestMixed(t *testing.T) {
	str := `العاشر ليونيكود (Unicode Conference)، الذي سيعقد في 10-12 آذار 1997 مبدينة`
	p := Paragraph{}
	p.SetString(str)
	order, err := p.Order()
	if err != nil {
		log.Fatal(err)
	}
	if p.IsLeftToRight() {
		t.Error("Paragraph should return LeftToRight() == false")
	}

	expectedRuns := []runInformation{
		{"العاشر ليونيكود (", RightToLeft, 0, 16},
		{"Unicode Conference", LeftToRight, 17, 34},
		{")، الذي سيعقد في ", RightToLeft, 35, 51},
		{"10", LeftToRight, 52, 53},
		{"-", RightToLeft, 54, 54},
		{"12", LeftToRight, 55, 56},
		{" آذار ", RightToLeft, 57, 62},
		{"1997", LeftToRight, 63, 66},
		{" مبدينة", RightToLeft, 67, 73},
	}

	if nr, expected := order.NumRuns(), len(expectedRuns); nr != expected {
		t.Errorf("Number of runs must be %d but got %d", expected, nr)
	}

	for i, er := range expectedRuns {
		r := order.Run(i)
		if str := r.String(); str != er.str {
			t.Errorf("Run %d should have string %q but has %q", i, er.str, str)
		}
		if s, e := r.Pos(); s != er.start || e != er.end {
			t.Errorf("Run %d should go from %d to %d but got %d to %d", i, er.start, er.end, s, e)
		}
		if d := r.Direction(); d != er.dir {
			t.Errorf("Run %d direction should be %d but got %d", i, er.dir, d)
		}
	}
}

func TestExplicitIsolate(t *testing.T) {
	// https://www.w3.org/International/articles/inline-bidi-markup/uba-basics.en#beyond
	str := "The names of these states in Arabic are \u2067مصر\u2069, \u2067البحرين\u2069 and \u2067الكويت\u2069 respectively."
	p := Paragraph{}
	p.SetString(str)
	order, err := p.Order()
	if err != nil {
		log.Fatal(err)
	}
	if !p.IsLeftToRight() {
		t.Error("Paragraph should return LeftToRight() == true")
	}

	expectedRuns := []runInformation{
		{"The names of these states in Arabic are \u2067", LeftToRight, 0, 40},
		{"مصر", RightToLeft, 41, 43},
		{"\u2069, \u2067", LeftToRight, 44, 47},
		{"البحرين", RightToLeft, 48, 54},
		{"\u2069 and \u2067", LeftToRight, 55, 61},
		{"الكويت", RightToLeft, 62, 67},
		{"\u2069 respectively.", LeftToRight, 68, 82},
	}

	if nr, expected := order.NumRuns(), len(expectedRuns); nr != expected {
		t.Errorf("Number of runs must be %d but got %d", expected, nr)
	}

	for i, er := range expectedRuns {
		r := order.Run(i)
		if str := r.String(); str != er.str {
			t.Errorf("Run %d should have string %q but has %q", i, er.str, str)
		}
		if s, e := r.Pos(); s != er.start || e != er.end {
			t.Errorf("Run %d should go from %d to %d but got %d to %d", i, er.start, er.end, s, e)
		}
		if d := r.Direction(); d != er.dir {
			t.Errorf("Run %d direction should be %d but got %d", i, er.dir, d)
		}
	}
}

func TestWithoutExplicitIsolate(t *testing.T) {
	str := "The names of these states in Arabic are مصر, البحرين and الكويت respectively."
	p := Paragraph{}
	p.SetString(str)
	order, err := p.Order()
	if err != nil {
		log.Fatal(err)
	}
	if !p.IsLeftToRight() {
		t.Error("Paragraph should return LeftToRight() == true")
	}

	expectedRuns := []runInformation{
		{"The names of these states in Arabic are ", LeftToRight, 0, 39},
		{"مصر, البحرين", RightToLeft, 40, 51},
		{" and ", LeftToRight, 52, 56},
		{"الكويت", RightToLeft, 57, 62},
		{" respectively.", LeftToRight, 63, 76},
	}

	if nr, expected := order.NumRuns(), len(expectedRuns); nr != expected {
		t.Errorf("Number of runs must be %d but got %d", expected, nr)
	}

	for i, er := range expectedRuns {
		r := order.Run(i)
		if str := r.String(); str != er.str {
			t.Errorf("Run %d should have string %q but has %q", i, er.str, str)
		}
		if s, e := r.Pos(); s != er.start || e != er.end {
			t.Errorf("Run %d should go from %d to %d but got %d to %d", i, er.start, er.end, s, e)
		}
		if d := r.Direction(); d != er.dir {
			t.Errorf("Run %d direction should be %d but got %d", i, er.dir, d)
		}
	}
}

func TestLongUTF8(t *testing.T) {
	str := `𠀀`
	p := Paragraph{}
	p.SetString(str)
	order, err := p.Order()
	if err != nil {
		log.Fatal(err)
	}
	if !p.IsLeftToRight() {
		t.Error("Paragraph should return LeftToRight() == true")
	}

	expectedRuns := []runInformation{
		{"𠀀", LeftToRight, 0, 0},
	}

	if nr, expected := order.NumRuns(), len(expectedRuns); nr != expected {
		t.Errorf("Number of runs must be %d but got %d", expected, nr)
	}

	for i, er := range expectedRuns {
		r := order.Run(i)
		if str := r.String(); str != er.str {
			t.Errorf("Run %d should have string %q but has %q", i, er.str, str)
		}
		if s, e := r.Pos(); s != er.start || e != er.end {
			t.Errorf("Run %d should go from %d to %d but got %d to %d", i, er.start, er.end, s, e)
		}
		if d := r.Direction(); d != er.dir {
			t.Errorf("Run %d direction should be %d but got %d", i, er.dir, d)
		}
	}
}

func TestLLongUTF8(t *testing.T) {
	strTester := []struct {
		str string
		l   int
	}{
		{"ö", 2},
		{"ॡ", 3},
		{`𠀀`, 4},
	}
	for _, st := range strTester {
		str := st.str
		expectedLen := st.l
		if _, l := LookupString(str); l != expectedLen {
			t.Errorf("LookupString(%q) length should return %d but got %d", str, expectedLen, l)
		}

	}

}

func TestMixedSimple(t *testing.T) {
	str := `Uا`
	p := Paragraph{}
	p.SetString(str)
	order, err := p.Order()
	if err != nil {
		log.Fatal(err)
	}
	if !p.IsLeftToRight() {
		t.Error("Paragraph should return LeftToRight() == true")
	}

	expectedRuns := []runInformation{
		{"U", LeftToRight, 0, 0},
		{"ا", RightToLeft, 1, 1},
	}

	if nr, expected := order.NumRuns(), len(expectedRuns); nr != expected {
		t.Errorf("Number of runs must be %d but got %d", expected, nr)
	}

	for i, er := range expectedRuns {
		r := order.Run(i)
		if str := r.String(); str != er.str {
			t.Errorf("Run %d should have string %q but has %q", i, er.str, str)
		}
		if s, e := r.Pos(); s != er.start || e != er.end {
			t.Errorf("Run %d should go from %d to %d but got %d to %d", i, er.start, er.end, s, e)
		}
		if d := r.Direction(); d != er.dir {
			t.Errorf("Run %d direction should be %d but got %d", i, er.dir, d)
		}
	}
}

func TestDefaultDirection(t *testing.T) {
	str := "+"
	p := Paragraph{}
	p.SetString(str, DefaultDirection(RightToLeft))
	_, err := p.Order()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if expected, dir := false, p.IsLeftToRight(); expected != dir {
		t.Errorf("Paragraph isLeftToRight should return %t but got %t", expected, dir)
	}
	p.SetString(str, DefaultDirection(LeftToRight))
	_, err = p.Order()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if expected, dir := true, p.IsLeftToRight(); expected != dir {
		t.Errorf("Paragraph isLeftToRight should return %t but got %t", expected, dir)
	}

}

func TestEmpty(t *testing.T) {
	p := Paragraph{}
	p.SetBytes([]byte{})
	_, err := p.Order()
	if err == nil {
		t.Error("p.Order must return an error on empty input")
	}
}

func TestNewline(t *testing.T) {
	str := "Hello\nworld"
	p := Paragraph{}
	n, err := p.SetString(str)
	if err != nil {
		t.Error(err)
	}
	if expected := 6; n != expected {
		t.Errorf("Length of SetString: expected %d but got %d", expected, n)
	}
}

func TestDoubleSetString(t *testing.T) {
	str := "العاشر ليونيكود (Unicode Conference)،"
	p := Paragraph{}
	_, err := p.SetString(str)
	if err != nil {
		t.Error(err)
	}
	_, err = p.SetString(str)
	if err != nil {
		t.Error(err)
	}
	_, err = p.Order()
	if err != nil {
		t.Error(err)
	}

}

func TestReverseString(t *testing.T) {
	input := "(Hello)"
	expected := "(olleH)"
	if str := ReverseString(input); str != expected {
		t.Errorf("ReverseString expected %q but got %q", expected, str)
	}
}

func TestAppendReverse(t *testing.T) {
	outString := "Hëllo"
	inString := "nice (wörld)"

	// empty in
	expected := "Hëllo"
	if r := AppendReverse([]byte(outString), []byte{}); string(r) != expected {
		t.Errorf("AppendReverse expected %q but got %q", expected, string(r))
	}

	// empty out
	expected = "(dlröw) ecin"
	if r := AppendReverse([]byte{}, []byte(inString)); string(r) != expected {
		t.Errorf("AppendReverse expected %q but got %q", expected, string(r))
	}

	// both given
	expected = "Hëllo(dlröw) ecin"
	if r := AppendReverse([]byte(outString), []byte(inString)); string(r) != expected {
		t.Errorf("AppendReverse expected %q but got %q", expected, string(r))
	}

}
