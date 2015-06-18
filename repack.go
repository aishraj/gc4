package main

import (
	"log"
)

// A repacker repacks trucks.
type repacker struct {
}

type Fit int

type shelf struct {
	minHeight, maxHeight, width uint8
}

const (
	_               = iota
	VerticalFit Fit = 1 << iota
	HorizontalFit
	NewShelfFit
	UnFit
)

func shelfNF(t *truck) (out *truck) {
	out = &truck{id: t.id}
	var boxes []box
	//collect all boxes
	for _, p := range t.pallets {
		for _, b := range p.boxes {
			boxes = append(boxes, b)
		}
	}
	var outPallets []pallet
	var uno pallet
	var shelves []shelf
	//for every box
	for _, item := range boxes {
		if len(shelves) == 0 {
			uno, shelves = makePallet(item)
			continue
		}
		fitness := findFit(item, &shelves)
		switch fitness {
		case HorizontalFit:
			item = sideWays(item)
			uno = addToPallet(uno, item, &shelves)
			break
		case VerticalFit:
			item = upRight(item)
			uno = addToPallet(uno, item, &shelves)
			break
		case NewShelfFit:
			topShelf := shelves[len(shelves)-1]
			item = sideWays(item)
			newTopShelf := shelf{minHeight: topShelf.maxHeight, maxHeight: (item.l + topShelf.maxHeight), width: 0}
			shelves = append(shelves, newTopShelf)
			uno = addToPallet(uno, item, &shelves)
			break
		case UnFit:
			//this means that uno is out of space.
			//we add it to the staging area and get a new pallet
			dummyTruck := truck{id: 0}
			dummyTruck.pallets = outPallets
			outPallets = append(outPallets, uno)
			uno, shelves = makePallet(item)
			break
		default:
			panic("Does not fit anywhere. Its an error")
		}
	}
	if len(uno.boxes) > 0 {
		outPallets = append(outPallets, uno)
	}
	//put it in the truck
	out.pallets = outPallets
	return
}

func findFit(item box, shelves *[]shelf) Fit {
	topShelf := (*shelves)[len(*shelves)-1]
	upBox := upRight(item)
	sideBox := sideWays(item)
	if upBox.l+topShelf.minHeight <= topShelf.maxHeight && (upBox.w+topShelf.width <= palletWidth) && (upBox.l+topShelf.minHeight <= palletLength) {
		return VerticalFit
	}
	if sideBox.l+topShelf.minHeight <= topShelf.maxHeight && (sideBox.w+topShelf.width <= palletWidth) && (sideBox.l+topShelf.minHeight <= palletLength) {
		return HorizontalFit
	}
	if (sideBox.l + topShelf.maxHeight) <= palletLength {
		return NewShelfFit
	}
	return UnFit
}

func addToPallet(uno pallet, item box, shelves *[]shelf) pallet {
	currentTop := (*shelves)[len(*shelves)-1]
	item.x = currentTop.minHeight
	item.y = currentTop.width
	currentTop.width += item.w
	(*shelves)[len(*shelves)-1] = currentTop
	uno.boxes = append(uno.boxes, item)
	return uno
}

func makePallet(item box) (packet pallet, shelves []shelf) {
	item = item.canon()

	bottomShelf := shelf{minHeight: 0, maxHeight: item.l, width: item.w}
	shelves = append(shelves, bottomShelf)
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
				log.Println("Last truck...")
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
