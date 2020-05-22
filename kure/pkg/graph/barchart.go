package graph

import "fmt"

type BarChart struct {
	Key   []string
	Count []float64
	Ratio []float64
}

func PrintBarChart(rawData []map[string]string) error {

	bc, err := fillData(rawData)
	if err != nil {
		return fmt.Errorf("error fill data: %w", err)
	}

	bc.print()

	return nil
}

func fillData(data []map[string]string) (*BarChart, error) {
	var bc BarChart

	// need to find unique keys in map
	// keyMap = uniqueKeys
	// The n count for number of unique keys in map

	for _, m := range data {
		for i := range m {

			for j, key := range bc.Key {

				if key == i {
					bc.Count[j]++
					fmt.Println(key + "=" + i)
				} else {
					bc.Key = append(bc.Key, i)
					bc.Count = append(bc.Count, 1)
				}
			}
			if len(bc.Key) == 0 {
				bc.Key = append(bc.Key, i)
				bc.Count = append(bc.Count, 1)
			}
		}
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
