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

// https://www.w3.org/International/articles/inline-bidi-markup/uba-basics
// https://www.w3.org/International/articles/inline-bidi-markup/uba-basics-data/directional_runs_rtl
//	str := "The names of these states in Arabic are \u2067مصر\u2069, \u2067البحرين\u2069 and \u2067الكويت\u2069 respectively."
