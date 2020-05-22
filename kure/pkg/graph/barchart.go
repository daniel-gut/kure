package graph

import "fmt"

type BarChart struct {
	Key   []string
	Count []float64
	Ratio []float64
}

func PrintBarChart(rawData []string) error {

	bc, err := fillData(rawData)
	if err != nil {
		return fmt.Errorf("error fill data: %w", err)
	}

	bc.print()

	return nil
}

func fillData(data []string) (*BarChart, error) {
	var (
		bc BarChart
	)

	uniqueKeys := make(map[string]float64)
	for _, k := range data {
		uniqueKeys[k]++
	}

	for k, v := range uniqueKeys {
		bc.Key = append(bc.Key, k)
		bc.Count = append(bc.Count, v)
	}

	bc.ratio()

	return &bc, nil

}

func (self *BarChart) ratio() {
	var sum float64

	// calculate sum
	for i := range self.Key {
		sum = sum + self.Count[i]
	}

	// calc and assign ratio
	for i := range self.Key {
		self.Ratio = append(self.Ratio, ((self.Count[i] * 100) / sum))
	}
}
