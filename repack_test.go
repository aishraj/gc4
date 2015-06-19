package main

import "testing"

func TestFitness(t *testing.T) {
	testBox := box{x: 0, y: 0, w: 1, l: 1, id: 1}
	_, shelves := makePallet(testBox)
	if len(shelves) > 1 {
		t.Errorf("The length of the shelves cannot exceed 1 after creating a single pallet. Current length is %v", len(shelves))
	}
	if len(shelves) <= 0 {
		t.Error("The length of shelves cannot be empty or less than 0.")
	}

	anotherBox := box{x: 0, y: 3, w: 2, l: 2, id: 2}
	boxFitness := findFit(anotherBox, &shelves)
	exepctedFitness := newShelfFit
	if boxFitness != exepctedFitness {
		t.Errorf("The box does not fit the shelf. Expected %v got %v", exepctedFitness, boxFitness)
	}

}
