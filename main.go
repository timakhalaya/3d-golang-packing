package main

import (
	"fmt"
	. "packages/packTypes"
	"sync"
)

/**
- идеальный слой формируем самой большой стороной вниз
- Brute force

-----

*/

/*
*
Compare two int
*/
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

/*
*
Function acts a single source for palette and item validations
*/
func validateInputValues(palette *Palette, item *Item) {
	checkItemSideIsZero(item)
	checkPaletteIsValid(palette)
}

func checkItemSideIsZero(item *Item) {
	if item.Width == 0 || item.Depth == 0 || item.Height == 0 {
		panic("Can not assign Zero Value")
	}
}

func checkPaletteIsValid(palette *Palette) {
	if palette.Width == 0 || palette.Depth == 0 || palette.Height == 0 {
		panic("Can not assign Zero Value to palette")
	}
}

func placeItemsOnPalette(palette Palette, item Item, chanel chan<- ResultPlacement) {
	/**
	Check if values provided are ok, e.g no zero values provided
	*/
	validateInputValues(&palette, &item)
	/**
	We can then check if volume of item and palette is close to each other , eg
	volume of palette 1600 and item is 1300, then of course we can only place one item
	*/
	/**
	We can move through palette width
	*/
	for i := 0; i <= palette.Width; i += item.Width {
		fmt.Println("Im at", i)
	}
	chanel <- ResultPlacement{
		Layers: nil,
		Volume: 0,
	}
}

/*
Run program here
*/
func main() {
	/**
	Check for errors using defer - when all is done...
	*/
	defer func() {
		args := recover()
		if args != nil {
			fmt.Println(args)
		}
	}()
	/*
		Start palette placement in 3 different go routines and wait for all to complete
		Capture results in Channel
	*/
	resultChan := make(chan ResultPlacement, 3)
	palettePlans := GetPalettePlans()
	wgMain := sync.WaitGroup{}
	wgMain.Add(len(palettePlans))
	for _, palette := range palettePlans {
		// create intermediary variable 'palette' , otherwise we can end up with unexpected results!!!
		palette := palette
		go func(channel chan<- ResultPlacement) {
			placeItemsOnPalette(palette, Item{
				Width:  2,
				Height: 6,
				Depth:  2,
				X:      0,
				Y:      0,
				Z:      0,
			}, resultChan)
			defer wgMain.Done()
		}(resultChan)
	}
	wgMain.Wait()
	close(resultChan)

chanSelect:
	for {
		select {
		case data, ok := <-resultChan:
			if ok {
				fmt.Println(data)
			} else {
				break chanSelect
			}
		default:
			//do nothing..
		}
	}
}
