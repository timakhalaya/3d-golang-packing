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

/*
*
The main function that is used in order to place items on palette
*/
func placeItemsOnPalette(palette Palette, item *Item, chanel chan<- ResultChanLayers) {
	var layers []Layer
	var totalItems int
	totalVolume := palette.Width * palette.Height * palette.Depth
	hasBeenRotated := false
	/**
	Check if one of the item sides is close to pallete size in proximity of 15%
	*/
	RotationIn3d(item, palette.Width, palette.Height, palette.Depth)
	/**
	Check if values provided are ok, e.g no zero values provided
	*/
	validateInputValues(&palette, item)
	/**
	We can then check if volume of item and palette is close to each other , eg
	volume of palette 1600 and item is 1300, then of course we can only place one item
	*/
	/**
	We can move through palette depth,width and height
	*/
	for z := 0; z <= palette.Depth; z += item.Depth {
		remainingHeight := palette.Height
		var layerItems []Coordinate
		for y := 0; y < palette.Height; y += item.Height {
			if remainingHeight < item.Height {
				RotateItem(item)
				hasBeenRotated = true
				if remainingHeight < item.Height {
					break
				}
			}
			remainingHeight -= item.Height
			remainingWidth := palette.Width
			for x := 0; x < palette.Width; x += item.Width {
				// test also if height is left check volume of space left and try to fit
				if (remainingWidth < item.Width) && remainingHeight == 0 {
					break
				} else if (remainingWidth < item.Width) && remainingHeight > 0 {
					if !hasBeenRotated {
						RotateItem(item)
						hasBeenRotated = true
					}
					if remainingHeight*item.Width < item.Width*item.Height {
						break
					}
				}
				remainingWidth -= item.Width

				if x+item.Width <= palette.Width && y+item.Height <= palette.Height && z+item.Depth <= palette.Depth {
					layerItems = append(layerItems, Coordinate{X: x, Y: y, Z: z, ItemDepth: item.Depth, ItemWidth: item.Width, ItemHeight: item.Height})
					totalItems++
				}
			}
			if hasBeenRotated {
				RotateItem(item)
				hasBeenRotated = false
			}
		}
		if hasBeenRotated {
			RotateItem(item)
			hasBeenRotated = false
		}
		if len(layerItems) > 0 {
			itemOneVolume := item.Width * item.Height
			itemAllVolume := len(layerItems) * itemOneVolume
			paletteVolume := palette.Width * palette.Height
			if paletteVolume-itemAllVolume > itemOneVolume {
				rearrangeLayer(layerItems, item, palette.Width, palette.Height)
			}
		}
		layers = append(layers, Layer{Items: layerItems})

	}
	volumeUsage := float64(totalItems*item.Width*item.Height*item.Depth) / float64(totalVolume) * 100.0
	/**
	As we are working with gorotuines we need to push data to channels
	*/
	chanel <- struct {
		Layers      []Layer
		VolumeUsed  float64
		ItemsPlaced int
		PaletteUsed Palette
	}{Layers: layers, VolumeUsed: volumeUsage, ItemsPlaced: totalItems, PaletteUsed: palette}
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
	resultChan := make(chan ResultChanLayers, 18)
	palettePlans := GetPalettePlans()
	/**
	WaitGroup is used in order to sync separate goroutines
	*/
	wgMain := sync.WaitGroup{}
	wgMain.Add(len(palettePlans))
	for _, palette := range palettePlans {
		// create intermediary variable 'palette' , otherwise we can end up with unexpected results!!!
		palette := palette
		go func(channel chan<- ResultChanLayers) {
			/**
			W D H
			Original - 1 2 3
			- Change width and depth  - 2 1 3
			- Change width and height  - 3 2 1
			- Change depth and height - 1 3 2
			- Change width and depth and depth and height 2 3 1
			- Change depth and width and width and height  3 1 2
			Original + 5 possible rotations
			*/
			item := Item{
				Width:  550,
				Height: 310,
				Depth:  550,
			}
			/**
			we need syncronization on this level as well
			We are rotating each item and finding the best placement
			*/
			childWg := sync.WaitGroup{}
			for i := 0; i < 6; i++ {
				childWg.Add(1)
				i := i
				go func() {
					roatatedItem := item.Rotate(i)
					defer childWg.Done()
					placeItemsOnPalette(palette, roatatedItem, resultChan)
				}()
			}
			childWg.Wait()
			defer wgMain.Done()
		}(resultChan)
	}
	wgMain.Wait()
	close(resultChan)

	var resSlice []ResultChanLayers
	/**
	The chanSelect below iterates over result channel

	resultChan - channel that hold results from each goroutine
	goroutine operates in async nature and places result in channel, so we can act on it
	*/
chanSelect:
	for {
		select {
		case data, ok := <-resultChan:
			if ok {
				resSlice = append(resSlice, data)
			} else {
				break chanSelect
			}
			/**
			default cases needed when no data added
			*/
		default:
			// no item placed in here
		}
	}
	/**
	This hold all placements available for different palettes
	e.g.
	euro palette placement
	industrial palette placement
	half palette placement

	Based on information in each slice we can find best placement for each palette , we can use items placed and volume used
	*/

	palettePlansMap := GetPalettePlans()
	var sliceEuroPalette []ResultChanLayers
	var sliceIndustrialPalette []ResultChanLayers
	var sliceHalfPalette []ResultChanLayers

	var bestEuroPalettePlacement ResultChanLayers
	var sliceIndustrialPalettePlacement ResultChanLayers
	var bestSilceHalfPalettePlacement ResultChanLayers

	/**
	Print placement here...
	*/
	for _, placement := range resSlice {
		if placement.PaletteUsed == palettePlansMap["EuroPalette"] {
			sliceEuroPalette = append(sliceEuroPalette, placement)
			if placement.VolumeUsed > bestEuroPalettePlacement.VolumeUsed && placement.ItemsPlaced > bestEuroPalettePlacement.ItemsPlaced {
				bestEuroPalettePlacement = placement
			}
		}
		if placement.PaletteUsed == palettePlansMap["IndustrialPalette"] {
			sliceIndustrialPalette = append(sliceIndustrialPalette, placement)
			if placement.VolumeUsed > sliceIndustrialPalettePlacement.VolumeUsed && placement.ItemsPlaced > sliceIndustrialPalettePlacement.ItemsPlaced {
				sliceIndustrialPalettePlacement = placement
			}
		}
		if placement.PaletteUsed == palettePlansMap["HalfPalette"] {
			sliceHalfPalette = append(sliceHalfPalette, placement)
			if placement.VolumeUsed > bestSilceHalfPalettePlacement.VolumeUsed && placement.ItemsPlaced > bestSilceHalfPalettePlacement.ItemsPlaced {
				bestSilceHalfPalettePlacement = placement
			}
		}
	}

	fmt.Println()
	fmt.Println()
	fmt.Println("Euro Pall")
	fmt.Printf("Palette used - %v \n", bestEuroPalettePlacement.PaletteUsed)
	for i, layer := range bestEuroPalettePlacement.Layers {
		if len(layer.Items) > 0 {
			fmt.Println("layer ", i)
			fmt.Println(layer.Items)
		}
	}
	fmt.Printf("Items placed %d ", bestEuroPalettePlacement.ItemsPlaced)
	fmt.Printf("Volume used %f ", bestEuroPalettePlacement.VolumeUsed)

	fmt.Println()
	fmt.Println()
	fmt.Println("Industrial Pall")
	fmt.Printf("Palette used - %v \n", sliceIndustrialPalettePlacement.PaletteUsed)
	for i, layer := range sliceIndustrialPalettePlacement.Layers {
		if len(layer.Items) > 0 {
			fmt.Println("layer ", i)
			fmt.Println(layer.Items)
		}
	}
	fmt.Printf("Items placed %d ", sliceIndustrialPalettePlacement.ItemsPlaced)
	fmt.Printf("Volume used %f ", sliceIndustrialPalettePlacement.VolumeUsed)

	fmt.Println()
	fmt.Println()
	fmt.Println("Half Pall")
	fmt.Printf("Palette used - %v \n", bestSilceHalfPalettePlacement.PaletteUsed)
	for i, layer := range bestSilceHalfPalettePlacement.Layers {
		if len(layer.Items) > 0 {
			fmt.Println("layer ", i)
			fmt.Println(layer.Items)
		}
	}
	fmt.Printf("Items placed %d ", bestSilceHalfPalettePlacement.ItemsPlaced)
	fmt.Printf("Volume used %f ", bestSilceHalfPalettePlacement.VolumeUsed)
}

func rearrangeLayer(coordinates []Coordinate, item *Item, paletteWidth int, paletteHeight int) {
	//rotate most right and try to fit --- find most right item
	initialWidth := paletteWidth
	// aggregate all width before last element
	allWidthWithLast := 0
	allHeightWithLast := 0
	for i, coordinate := range coordinates {
		allWidthWithLast += coordinate.ItemWidth
		allHeightWithLast += coordinate.ItemHeight
		initialWidth = initialWidth - coordinate.ItemWidth
		if initialWidth <= 0 {
			// last item in the row
			lastInLayerRow := coordinates[i-1]
			nextInLayerRow := coordinate
			// remove last width from all and add height if rotated
			swappedXYSummWidth := allWidthWithLast - lastInLayerRow.ItemWidth + lastInLayerRow.ItemHeight
			swappedXYSummHeight := allHeightWithLast - lastInLayerRow.ItemHeight + lastInLayerRow.ItemWidth
			if swappedXYSummWidth <= paletteWidth && swappedXYSummHeight <= paletteHeight {
				// check that there is no overlap with item below
				fmt.Println(nextInLayerRow)
			}

		}
	}
}
