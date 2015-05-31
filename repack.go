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
			//log.Println(b.id)
			boxes = append(boxes, b)
		}
	}
	//the worst we can do is have one pallet per box
	outPallets := make([]pallet, 0, len(boxes))
	var uno pallet
	var shelves []shelf
	//for every box
	for _, item := range boxes {
		//log.Println("The state of UNO is:", uno.OneLine())
		if len(shelves) == 0 {
			//log.Println("Box %v is the first one\n", item.id)
			uno, shelves = makePallet(item)
			continue
		}
		fitness := findFit(item, &shelves)
		switch fitness {
		case HorizontalFit:
			//log.Printf("Box %v fits horizontal.\n", item.id)
			item = sideWays(item)
			uno = addToPallet(uno, item, &shelves)
			break
		case VerticalFit:
			//log.Printf("Box %v fits vertical.\n", item.id)
			item = upRight(item)
			uno = addToPallet(uno, item, &shelves)
			break
		case NewShelfFit:
			//log.Printf("Box %v fits in a new shelf \n", item.id)
			topShelf := shelves[len(shelves)-1]
			//log.Println("Topshelf is :", topShelf)
			item = sideWays(item)
			newTopShelf := shelf{minHeight: topShelf.maxHeight, maxHeight: (item.l + topShelf.maxHeight), width: 0}
			//log.Println("New Top shelf is:", newTopShelf)
			shelves = append(shelves, newTopShelf)
			//log.Println("checking if uno is valid BEFORE: ", uno.IsValid())
			uno = addToPallet(uno, item, &shelves)
			//log.Println("checking if uno is valid AFTER: ", uno.IsValid())
			break
		case UnFit:
			//log.Printf("Box %v does not fit in this pallet. Need a new pallet", item.id)
			//this means that uno is out of space.
			//we add it to the staging area and get a new pallet
			dummyTruck := truck{id: 0}
			dummyTruck.pallets = outPallets
			//log.Println("Current state of pallets is:", dummyTruck)
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
	//log.Println("Outgoing truck is :", out)
	return
}

func findFit(item box, shelves *[]shelf) Fit {
	//log.Println("shelves are:", shelves)
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
