package main

import (
	"fmt"
)

type Item struct {
	Width, Depth, Height int
}

type Palette struct {
	Width, Depth, Height int
}

type PlacedItem struct {
	X, Y, Z              int
	Width, Depth, Height int
}

type Layer struct {
	Items  []PlacedItem
	Height int
}

func main() {
	palette := Palette{Width: 1000, Depth: 1000, Height: 1000}
	item := Item{Width: 200, Depth: 391, Height: 200}

	layers, itemsPlaced, volumeUtilized := placeItems(palette, item)
	fmt.Printf("Number of layers: %d\n", len(layers))
	fmt.Printf("Number of items placed: %d\n", itemsPlaced)
	fmt.Printf("Volume utilized: %.2f%%\n", volumeUtilized)

	for i, layer := range layers {
		fmt.Printf("Layer %d (Height %d):\n", i+1, layer.Height)
		for _, placedItem := range layer.Items {
			fmt.Printf("  Item placed at (%d, %d, %d) of size (%d, %d, %d)\n", placedItem.X, placedItem.Y, placedItem.Z, placedItem.Width, placedItem.Depth, placedItem.Height)
		}
	}
}

func placeItems(palette Palette, item Item) ([]Layer, int, float64) {
	var layers []Layer
	remainingHeight := palette.Height
	itemsPlaced := 0
	totalVolume := 0
	paletteVolume := palette.Width * palette.Depth * palette.Height

	for remainingHeight > 0 {
		layer, layerHeight := placeLayer(palette, item, remainingHeight)
		if layerHeight == 0 {
			break
		}
		layers = append(layers, layer)
		remainingHeight -= layerHeight
		itemsPlaced += len(layer.Items)
		totalVolume += layerHeight * palette.Width * palette.Depth
	}

	volumeUtilized := (float64(totalVolume) / float64(paletteVolume)) * 100
	return layers, itemsPlaced, volumeUtilized
}

func placeLayer(palette Palette, item Item, remainingHeight int) (Layer, int) {
	layer := Layer{}
	layerHeight := 0
	space := make([][]bool, palette.Width)
	for i := range space {
		space[i] = make([]bool, palette.Depth)
	}

	for x := 0; x < palette.Width; x++ {
		for y := 0; y < palette.Depth; y++ {
			if !space[x][y] {
				rotatedItem := item.rotateAndPlace(palette, x, y, space, remainingHeight)
				if rotatedItem.Height > 0 {
					layerHeight = max(layerHeight, rotatedItem.Height)
					layer.Items = append(layer.Items, rotatedItem)
				}
			}
		}
	}

	layer.Height = layerHeight
	return layer, layerHeight
}

func (item Item) rotateAndPlace(palette Palette, x, y int, space [][]bool, remainingHeight int) PlacedItem {
	for rotation := 0; rotation < 6; rotation++ {
		rotatedItem := item.rotate(rotation)
		if rotatedItem.Height <= remainingHeight && canPlace(space, x, y, rotatedItem) {
			place(space, x, y, rotatedItem)
			return PlacedItem{X: x, Y: y, Z: palette.Height - remainingHeight, Width: rotatedItem.Width, Depth: rotatedItem.Depth, Height: rotatedItem.Height}
		}
	}
	return PlacedItem{}
}

func (item Item) rotate(rotation int) Item {
	switch rotation {
	case 1:
		return Item{Width: item.Depth, Depth: item.Height, Height: item.Width}
	case 2:
		return Item{Width: item.Height, Depth: item.Width, Height: item.Depth}
	case 3:
		return Item{Width: item.Width, Depth: item.Height, Height: item.Depth}
	case 4:
		return Item{Width: item.Height, Depth: item.Depth, Height: item.Width}
	case 5:
		return Item{Width: item.Depth, Depth: item.Width, Height: item.Height}
	default:
		return item
	}
}

func canPlace(space [][]bool, x, y int, item Item) bool {
	if x+item.Width > len(space) || y+item.Depth > len(space[0]) {
		return false
	}
	for i := x; i < x+item.Width; i++ {
		for j := y; j < y+item.Depth; j++ {
			if space[i][j] {
				return false
			}
		}
	}
	return true
}

func place(space [][]bool, x, y int, item Item) {
	for i := x; i < x+item.Width; i++ {
		for j := y; j < y+item.Depth; j++ {
			space[i][j] = true
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
