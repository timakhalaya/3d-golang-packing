package PackerTypes

type ResultPlacement struct {
	Layers []*Layer
	Volume int
}

type Palette struct {
	Width, Height, Depth int
}

/*
Result chan Layer struct needed to push items in chan
*/
type ResultChanLayers struct {
	Layers      []Layer
	VolumeUsed  float64
	ItemsPlaced int
	PaletteUsed Palette
}

type Item struct {
	RotationCalls        int
	Width, Height, Depth int
	X, Y, Z              int
}

type Layer struct {
	Items  []Coordinate
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

func (item *Item) Rotate(rotation int) *Item {
	item.RotationCalls = item.RotationCalls + 1
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

	initMap["IndustrialPalette"] = Palette{
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

type Coordinate struct {
	X, Y, Z                          int
	ItemWidth, ItemHeight, ItemDepth int
}

// RotateItem rotates the item dimensions to maximize space utilization within the layer
func RotateItem(item *Item) {
	// Swap width and height
	item.Width, item.Height = item.Height, item.Width
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
}

/**
Palette plans
*/
