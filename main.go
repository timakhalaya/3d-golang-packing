package main

import (
	"fmt"
	"strings"
)

type Item struct {
	Width  int
	Height int
	Depth  int
}

type Coordinate struct {
	X, Y, Z                          int
	itemWidth, itemHeight, itemDepth int
}

type Layer struct {
	Items []*Coordinate
}

func generateSVG(layers []Layer, itemWidth, itemHeight, paletteWidth, paletteHeight int) string {
	var svg strings.Builder

	// Start SVG element
	svg.WriteString("<svg xmlns=\"http://www.w3.org/2000/svg\" ")
	svg.WriteString(fmt.Sprintf("width=\"%d\" height=\"%d\">\n", paletteWidth*itemWidth, paletteHeight*itemHeight))

	// Draw palette rectangle
	svg.WriteString(fmt.Sprintf("<rect x=\"0\" y=\"0\" width=\"%d\" height=\"%d\" fill=\"none\" stroke=\"black\"/>\n", paletteWidth*itemWidth, paletteHeight*itemHeight))

	// Draw items
	for _, layer := range layers {
		for _, item := range layer.Items {
			// Draw item border
			svg.WriteString(fmt.Sprintf("<rect x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" fill=\"none\" stroke=\"black\"/>\n", item.X*itemWidth, item.Y*itemHeight, itemWidth, itemHeight))
			// Draw filled item
			svg.WriteString(fmt.Sprintf("<rect x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" fill=\"blue\"/>\n", item.X*itemWidth, item.Y*itemHeight, itemWidth, itemHeight))
		}
	}

	// End SVG element
	svg.WriteString("</svg>")

	return svg.String()
}

func placeItemsOnPalette(paletteWidth, paletteHeight, paletteDepth int, item Item) ([]Layer, int, float64) {
	var layers []Layer
	var totalItems int
	totalVolume := paletteWidth * paletteHeight * paletteDepth
	hasBeenRotated := false

	/**
	Check if one of the item sides is close to pallete size in proximity of 15%
	*/
	RotationIn3d(&item, paletteWidth, paletteHeight, paletteDepth)

	//the depth should be added based on item depth
	for z := 0; z < paletteDepth; z += item.Depth {
		remainingHeight := paletteHeight
		var layerItems []*Coordinate

		// y is free until we iterate over all items height
		for y := 0; y < paletteHeight; y += item.Height {
			if remainingHeight < item.Height {
				RotateItem(&item)
				hasBeenRotated = true
				if remainingHeight < item.Height {
					break
				}
			}
			remainingHeight -= item.Height

			remainingWidth := paletteWidth
			for x := 0; x < paletteWidth; x += item.Width {
				// test also if height is left check volume of space left and try to fit
				if (remainingWidth < item.Width) && remainingHeight == 0 {
					break
				} else if (remainingWidth < item.Width) && remainingHeight > 0 {
					if !hasBeenRotated {
						RotateItem(&item)
						hasBeenRotated = true
					}
					if remainingHeight*item.Width < item.Width*item.Height {
						break
					}
				}
				remainingWidth -= item.Width

				if x+item.Width <= paletteWidth && y+item.Height <= paletteHeight && z+item.Depth <= paletteDepth {
					layerItems = append(layerItems, &Coordinate{X: x, Y: y, Z: z, itemDepth: item.Depth, itemWidth: item.Width, itemHeight: item.Height})
					totalItems++
				}
			}
			if hasBeenRotated {
				RotateItem(&item)
				hasBeenRotated = false
			}
		}
		if hasBeenRotated {
			RotateItem(&item)
			hasBeenRotated = false
		}
		/**
		Layer might not be efficient then we need to recalculate
		*/
		if len(layerItems) > 0 {
			itemOneVolume := item.Width * item.Height
			itemAllVolume := len(layerItems) * itemOneVolume
			paletteVolume := paletteWidth * paletteHeight
			if paletteVolume-itemAllVolume > itemOneVolume {
				rearrangeLayer(layerItems, &item, paletteWidth, paletteHeight)
			}
		}
		layers = append(layers, Layer{Items: layerItems})
		//after layer has been added we need to rotate item to it's previous state
	}

	volumeUsage := float64(totalItems*item.Width*item.Height*item.Depth) / float64(totalVolume) * 100.0
	return layers, totalItems, volumeUsage
}

func rearrangeLayer(coordinates []*Coordinate, item *Item, paletteWidth int, paletteHeight int) {
	fmt.Println("There should be still available place left!!!")
	//rotate most right and try to fit --- find most right item
	initialWidth := paletteWidth
	// aggregate all width before last element
	allWidthWithLast := 0
	allHeightWithLast := 0
	for i, coordinate := range coordinates {
		allWidthWithLast += coordinate.itemWidth
		allHeightWithLast += coordinate.itemHeight
		initialWidth = initialWidth - coordinate.itemWidth
		if initialWidth <= 0 {
			// last item in the row
			lastInLayerRow := coordinates[i-1]
			nextInLayerRow := coordinate
			// remove last width from all and add height if rotated
			swappedXYSummWidth := allWidthWithLast - lastInLayerRow.itemWidth + lastInLayerRow.itemHeight
			swappedXYSummHeight := allHeightWithLast - lastInLayerRow.itemHeight + lastInLayerRow.itemWidth
			if swappedXYSummWidth <= paletteWidth && swappedXYSummHeight <= paletteHeight {
				// check that there is no overlap with item below
				fmt.Println(nextInLayerRow)
			}

		}
	}
}

// RotateItem rotates the item dimensions to maximize space utilization within the layer
func RotateItem(item *Item) {
	// Swap width and height
	item.Width, item.Height = item.Height, item.Width
}

func main() {
	paletteWidth := 120
	paletteHeight := 80
	paletteDepth := 180
	/**
	Check rotation
	*/
	item := Item{Width: 120, Height: 80, Depth: 155}

	result, totalItems, volumeUsage := placeItemsOnPalette(paletteWidth, paletteHeight, paletteDepth, item)

	//svg := generateSVG(result, item.Width, item.Height, paletteWidth, paletteHeight)
	//fmt.Println(svg)

	for i, layer := range result {
		fmt.Printf("Layer %d:\n", i+1)
		for _, coordinate := range layer.Items {
			fmt.Printf("Item at coordinate (%d, %d, %d) placed with dimensions of (width: %d , height: %d, depth: %d) \n", coordinate.X, coordinate.Y, coordinate.Z, coordinate.itemWidth, coordinate.itemHeight, coordinate.itemDepth)
		}
	}

	fmt.Printf("Total Items Placed: %d\n", totalItems)
	fmt.Printf("Volume Usage: %.2f%%\n", volumeUsage)
}

func RotationIn3d(item *Item, palWidth int, palHeight int, palDepth int) {
	/*
		Большую сторону надо двигать к наиболее большой стороне паллеты
	*/
	// need to rotate item height and depth
	if item.Height > palHeight && item.Height <= palDepth && palHeight >= item.Depth {
		itemDepth := item.Depth
		item.Depth = item.Height
		item.Height = itemDepth
	}

	// need to rotate item height and width
	if item.Height > palHeight && item.Height <= palWidth && palHeight >= item.Width {
		itemHeight := item.Height
		item.Height = item.Width
		item.Width = itemHeight
	}
	fmt.Println(item)
}
