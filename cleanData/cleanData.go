package cleanData

import (
	"errors"
	"fmt"
	"scrap-ifix-go/b3"
	"strconv"
	"strings"
)

type IfixFundData struct {
	b3.Asset
	Price    string
	Dividend string
}
type Fund struct {
	Code     string
	Yield    string
	Price    float64
	MaxPrice float64
}

func (f *IfixFundData) CleanData(discountTax float64) (Fund, error) {
	f.Dividend = strings.ReplaceAll(f.Dividend, "%", "")
	f.Dividend = strings.ReplaceAll(f.Dividend, ",", ".")
	f.Price = strings.ReplaceAll(f.Price, "R$", "")
	f.Price = strings.ReplaceAll(f.Price, ",", ".")
	f.Price = strings.TrimSpace(f.Price)
	f.Dividend = strings.TrimSpace(f.Dividend)

	// Check for empty strings before parsing
	if f.Dividend == "" {
		fmt.Printf("Empty dividend for fund %s, skipping\n", f.Code)
		return Fund{}, errors.New("empty dividend value")
	}
	if f.Price == "" {
		fmt.Printf("Empty price for fund %s, skipping\n", f.Code)
		return Fund{}, errors.New("empty price value")
	}

	dividendNumber, err := strconv.ParseFloat(f.Dividend, 64)
	if err != nil {
		fmt.Println("Error occurred while parsing dividend", err)
		return Fund{}, errors.New("error occurred while parsing dividend")
	}
	priceNumber, err := strconv.ParseFloat(f.Price, 64)
	if err != nil {
		fmt.Println("Error occurred while parsing price", err)
		return Fund{}, errors.New("error occurred while parsing price")
	}

	maxPrice := dividendNumber / discountTax

	fund := Fund{f.Code, f.Dividend, priceNumber, maxPrice}
	return fund, nil
}
