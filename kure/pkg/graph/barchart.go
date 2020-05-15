package graph

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

type BarChart struct {
	Key   []string
	Value []float64
	Ratio []float64
}

func NewBarChart() *BarChart {
	return &BarChart{
		Key:   []string{},
		Value: []float64{},
		Ratio: []float64{},
	}
}

func MapToBarChart(rawData map[string]float64) (*BarChart, error) {

	var (
		keys   []string
		values []float64
		ratio  []float64
		bc     BarChart
	)

	for k, v := range rawData {
		keys = append(keys, k)
		values = append(values, v)
		// otherwise lenght is 0
		ratio = append(ratio, float64(0))
	}

	bc.Key = keys
	bc.Value = values
	bc.Ratio = ratio

	return &bc, nil

}

func (self *BarChart) Print() {

	self.RatioCalc()

	writer := tabwriter.NewWriter(os.Stdout, 2, 2, 2, ' ', 0)

	title := fmt.Sprintf("%s\t%s\t%s\t%s", "Bucket Name", "Ratio", "Graph", "Count")
	fmt.Fprintln(writer, title)

	boarder := fmt.Sprintf("%s\t%s\t%s\t%s", "-----------", "-----", "-----", "-----")
	fmt.Fprintln(writer, boarder)

	for i := range self.Key {

		barLength := int(self.Ratio[i]) / 2 // 100% == 50 Blocks
		bar := strings.Repeat(string('█'), barLength)

		output := fmt.Sprintf("%s\t%5.2f%%\t%s\t%5.0f", self.Key[i], self.Ratio[i], bar, self.Value[i])

		fmt.Fprintln(writer, output)
	}

	writer.Flush()

}

func (self *BarChart) RatioCalc() {
	var sum float64

	// calculate sum
	for i := range self.Key {
		sum = sum + self.Value[i]
	}

	// calc and assign ratio
	for i := range self.Key {
		self.Ratio[i] = ((self.Value[i] * 100) / sum)
	}

}
