package main

import (
	"fmt"

	"example.com/price-calculator/filemanager"
	"example.com/price-calculator/prices"
)

func main() {
	taxRates := []float64{0, 0.07, 0.1, 0.15}
	doneChans := make([]chan bool, len(taxRates))
	errorChans := make([]chan error, len(taxRates))

	for index, taxRate := range taxRates {
		doneChans[index] = make(chan bool)
		errorChans[index] = make(chan error)

		fm := filemanager.New("Prices.txt", fmt.Sprintf("result_%.0f.json", taxRate*100))
		//cmd := cmdmanager.New()

		priceJob := prices.NewTaxIncludedPriceJob(fm, taxRate)
		go priceJob.Process(doneChans[index], errorChans[index])
	}

	for indexRate := range taxRates {
		select {
		case err := <-errorChans[indexRate]:
			if err != nil {
				fmt.Println(err)
			}
		case <-doneChans[indexRate]:
			fmt.Println("Done")
		}
	}
}
