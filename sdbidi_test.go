package sdbidi

import (
	"log"
	"testing"
)

func TestMixed(t *testing.T) {
	str := "العاشر ليونيكود (Unicode Conference)،"
	p := Paragraph{}
	p.SetString(str)
	order, err := p.Order()
	if err != nil {
		log.Fatal(err)
	}
	if nr, expected := order.NumRuns(), 3; nr != expected {
		t.Errorf("Number of runs must be %d but got %d", expected, nr)
	}
	if p.IsLeftToRight() {
		t.Error("Paragraph should return LeftToRight() == false")
	}
}

func TestSimple(t *testing.T) {
	str := "Hello"
	p := Paragraph{}
	p.SetString(str)
	order, err := p.Order()
	if err != nil {
		log.Fatal(err)
	}

	if nr, expected := order.NumRuns(), 1; nr != expected {
		t.Errorf("Number of runs must be %d but got %d", expected, nr)
	}
	if !p.IsLeftToRight() {
		t.Error("Paragraph should return LeftToRight() == true")
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
