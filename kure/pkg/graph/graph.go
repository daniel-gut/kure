package graph

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

type graphData struct {
	bucketName string
	count      float64
	ratio      float64
}

// Print displays graph output
func Print(rawData map[string]float64) {

	const (
		width  = 50
		height = 20

		cellEmpty = ' '
		cellFull  = '█'
	)

	var gd = make([]graphData, len(rawData))

	// rewind the slice (allow appending from the beginning)
	gd = gd[:0]

	for k, v := range rawData {

		maxChars := 20
		if len(k) < 20 {
			maxChars = len(k)
		}

		data := graphData{bucketName: k[:maxChars], count: v}
		gd = append(gd, data)

	}

	// calculate ratio and assign value to struct
	gd = ratio(gd)

	// Print graph
	printGraph(gd)
}

// Clear clears the screen
func Clear() {
	fmt.Print("\033[2J")
}

// MoveTopLeft moves the cursor to the top left position of the screen
func MoveTopLeft() {
	fmt.Print("\033[H")
}

func ratio(data []graphData) []graphData {
	var sum float64

	// calculate sum
	for _, d := range data {
		sum = sum + d.count
	}

	// calc and assign ratio
	for i, d := range data {
		data[i].ratio = ((d.count * 100) / sum)
	}

	return data
}

func printGraph(gd []graphData) {

	writer := tabwriter.NewWriter(os.Stdout, 2, 2, 2, ' ', 0)

	title := fmt.Sprintf("%s\t%s\t%s\t%s", "Bucket Name", "Ratio", "Graph", "Count")
	fmt.Fprintln(writer, title)

	boarder := fmt.Sprintf("%s\t%s\t%s\t%s", "-----------", "-----", "-----", "-----")
	fmt.Fprintln(writer, boarder)

	for _, d := range gd {

		barLength := int(d.ratio) / 2
		bar := strings.Repeat(string('█'), barLength)

		output := fmt.Sprintf("%s\t%5.2f%%\t%s\t%5.0f", d.bucketName, d.ratio, bar, d.count)

		fmt.Fprintln(writer, output)
	}

	writer.Flush()
}
