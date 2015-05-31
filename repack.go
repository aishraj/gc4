package main

import "fmt"

// A repacker repacks trucks.
type repacker struct {
}

// This repacker is the worst possible, since it uses a new pallet for
// every box. Your job is to replace it with something better.
func oneBoxPerPallet(t *truck) (out *truck) {
	out = &truck{id: t.id}
	for _, p := range t.pallets {
		for _, b := range p.boxes {
			b.x, b.y = 0, 0
			out.pallets = append(out.pallets, pallet{boxes: []box{b}})
		}
	}
	return
}

func shelfNF(t *truck) (out *truck) {
	out = &truck{id: t.id}
	var boxes []box
	//collect all boxes
	for _, p := range t.pallets {
		for _, b := range p.boxes {
			boxes = append(boxes, b)
		}
	}
	//the worst we can do is have one pallet per box
	outPallets := make([]pallet, 0, len(boxes))
	var uno pallet
	var shelf []uint8
	//for every box
	for _, item := range boxes {
		//if its the first pallet, then create a new shelf and add the box to it
		if len(outPallets) == 0 {
			uno, shelf = makePallet(item)
		}
		//if we can add to the current pallet add the box, if requied update the shelf
		//todo figure out if this is pass by reference or value, if now passing a pointer so that i can modify it and is reflected here
		if fitsCurrentShelf(uno, item, &shelf) {
			uno = addToPallet(uno, item, &shelf)
		} else if fitsNextShelf(uno, item, &shelf) {
			uno = addToNewShelf(uno, item, &shelf)
		} else {
			outPallets = append(outPallets, uno)
			uno, shelf = makePallet(item)
		}
	}
	//put it in the truck
	out.pallets = outPallets
	//steps
	// for every pallet we get, iterate over the boxes
	// for every box, try to see if this is the first box. if yes put it sideWays
	// if this is not the first box, try to put vertically
	// if we can't put it vertically, try putting it sideWays
	// if that does not work either, move on to the next pallet (brand new)
	// do this for all boxes irresepctive of which pallet they came in
	return
}

func fitsCurrentShelf(uno pallet, item box, shelf *[]uint8) bool {
	//TODO: add
	//first check if the current shelf has space. all we need to do is check if the current box fits
	// first in a vertical way.
	//if not in a horizontal way in the current shelf (we need to check the width)
	if len(*shelf) == 0 {
		horizontalBox := item.canon()
		*shelf = append(*shelf, horizontalBox.l)
		return true
	}
	shelfHeignt := (*shelf)[len(*shelf)-1]
	verticalBox := upRight(item)
	if verticalBox.l < shelfHeignt {
		//now we check the width
		
	}
	horizontalBox := sideWays(item)
	if horizontalBox.
}

func fitsNextShelf(uno pallet, item box, shelf *[]uint8) bool {
	//TODO: add
	//first check if the current shelf has space
	// here we need to check the remaining height and then decide
	return false
}

func addToPallet(uno pallet, item box, shelf *[]uint8) pallet {
	//it should keep track of the shelves
	//probably a map is a good way to keep track of this
	//map[pallet][]int ie map of a pallet to a slice of int
	//or maybe not even a map, just have a slice of int, and pass it over here. each time we add a box to the pallet, we add the current top shelf
	//if the top slef i
	//each element of the slice represents the current peak in the pallet
	//the peak cannot exceed height of the pallet ie (peak < max(palletWidth, palletLength)
	return uno
}

func addToNewShelf(uno pallet, item box, shelf *[]uint8) pallet {
	return uno
}

func makePallet(item box) (packet pallet, shelf []uint8) {
	item = sideWays(item)
	item.x = 0
	item.y = 0

	shelf = append(shelf, item.l)
	packet = pallet{boxes: []box{item}}
	return
}

func sideWays(inbox box) (outbox box) {
	outbox = inbox
	if outbox.w < outbox.l {
		outbox.l, outbox.w = outbox.w, outbox.l
	}
	return
}

func upRight(inbox box) (outbox box) {
	outbox = inbox
	if outbox.l < outbox.w {
		outbox.l, outbox.w = outbox.w, outbox.l
	}
	return
}

func newRepacker(in <-chan *truck, out chan<- *truck) *repacker {
	go func() {
		for t := range in {
			// The last truck is indicated by its id. You might
			// need to do something special here to make sure you
			// send all the boxes.
			if t.id == idLastTruck {
				fmt.Println("Last truck crap...")
			}

			//t = oneBoxPerPallet(t)
			t = shelfNF(t)
			out <- t
		}
		// The repacker must close channel out after it detects that
		// channel in is closed so that the driver program will finish
		// and print the stats.
		close(out)
	}()
	return &repacker{}
}
