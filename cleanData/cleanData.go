package cleandata

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

func (f *IfixFundData) CleanData(descountTax float64) (Fund, error) {
	f.Dividend = strings.ReplaceAll(f.Dividend, "%", "")
	f.Dividend = strings.ReplaceAll(f.Dividend, ",", ".")
	f.Price = strings.ReplaceAll(f.Price, "R$", "")
	f.Price = strings.ReplaceAll(f.Price, ",", ".")
	f.Price = strings.TrimSpace(f.Price)
	f.Dividend = strings.TrimSpace(f.Dividend)

	dividendNumber, err := strconv.ParseFloat(f.Dividend, 64)
	if err != nil {
		fmt.Println("Error occured while parsing dividend", err)
		return Fund{}, errors.New("error occured while parsing dividend")
	}
	priceNumber, err := strconv.ParseFloat(f.Price, 64)
	if err != nil {
		fmt.Println("Error occured while parsing price", err)
		return Fund{}, errors.New("error occured while parsing price")
	}

	maxPrice := dividendNumber / descountTax

	fund := Fund{f.Code, f.Dividend, priceNumber, maxPrice}
	return fund, nil
}
