package graph

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func (bc BarChart) print() {

	writer := tabwriter.NewWriter(os.Stdout, 2, 2, 2, ' ', 0)

	title := fmt.Sprintf("%s\t%s\t%s\t%s", "Bucket Name", "Ratio", "Graph", "Count")
	fmt.Fprintln(writer, title)

	boarder := fmt.Sprintf("%s\t%s\t%s\t%s", "-----------", "-----", "-----", "-----")
	fmt.Fprintln(writer, boarder)

	for i := range bc.Key {

		barLength := int(bc.Ratio[i]) / 2 // 100% == 50 Blocks
		bar := strings.Repeat("█", barLength)

		output := fmt.Sprintf("%s\t%5.2f%%\t%s\t%5.0f", bc.Key[i], bc.Ratio[i], bar, bc.Count[i])

		fmt.Fprintln(writer, output)
	}

	writer.Flush()

}
