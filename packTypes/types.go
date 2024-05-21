package PackerTypes

type ResultPlacement struct {
	Layers []*Layer
	Volume int
}

type Palette struct {
	Width, Height, Depth int
}

type Item struct {
	Width, Height, Depth int
	X, Y, Z              int
}

type Layer struct {
	Items  []*Item
	Height int
}

type PaletteMap = map[string]Palette

/*
*
Useful function for setting all dimensions of the item - using pointer reciever
*/
func (item *Item) setItemWDH(width, height, depth int) {
	item.Width = width
	item.Height = height
	item.Depth = depth
}

/*
Function for setting width of the item - using pointer reciever
*/
func (item *Item) setWidth(width int) {
	item.Width = width
}

/*
Function for setting height of the item - using pointer reciever
*/
func (item *Item) setHeight(height int) {
	item.Height = height
}

/*
Function for setting depth of the item - using pointer reciever
*/
func (item *Item) setDepth(depth int) {
	item.Depth = depth
}

/**
Function for rotating item dimensions based on rotation provided
*/

func (item *Item) rotate(rotation int) *Item {
	switch rotation {
	case 1:
		return &Item{Width: item.Depth, Depth: item.Height, Height: item.Width}
	case 2:
		return &Item{Width: item.Height, Depth: item.Width, Height: item.Depth}
	case 3:
		return &Item{Width: item.Width, Depth: item.Height, Height: item.Depth}
	case 4:
		return &Item{Width: item.Height, Depth: item.Depth, Height: item.Width}
	case 5:
		return &Item{Width: item.Depth, Depth: item.Width, Height: item.Height}
	default:
		return item
	}
}

func GetPalettePlans() PaletteMap {
	initMap := make(map[string]Palette)

	initMap["EuroPalette"] = Palette{
		Width:  1200,
		Height: 800,
		Depth:  800,
	}

	initMap["IndustriePalette"] = Palette{
		Width:  1200,
		Height: 1000,
		Depth:  800,
	}

	initMap["HalfPalette"] = Palette{
		Width:  800,
		Height: 600,
		Depth:  800,
	}

	return initMap
}

/**
Palette plans
*/
