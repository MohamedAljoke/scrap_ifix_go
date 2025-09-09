package fundsCsv

import (
	"encoding/csv"
	"os"
	"scrap-ifix-go/cleanData"
	"strconv"
)

func CreateCSVFromFunds(funds []cleanData.Fund, filename string) error {

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"Code", "Yield", "Price", "MaxPrice"}

	if err := writer.Write(header); err != nil {
		return err
	}

	for _, fund := range funds {
		row := []string{
			fund.Code,
			fund.Yield,
			strconv.FormatFloat(fund.Price, 'f', -1, 64),
			strconv.FormatFloat(fund.MaxPrice, 'f', -1, 64),
		}

		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}
