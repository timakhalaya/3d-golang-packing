package svg

import (
	"fmt"
	. "packages/packTypes"
	"strings"
)

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
