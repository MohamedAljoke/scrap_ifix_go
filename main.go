package main

import (
	"fmt"
	"log"
	"scrap-ifix-go/b3"
	cleandata "scrap-ifix-go/cleanData"
	"scrap-ifix-go/fundsCsv"
	"scrap-ifix-go/scrapper"
	"sync"
	"time"
)

var taxaIpcaPorcentage = 6.1
var riskRatePercentage = 3.0

var descountTax = (taxaIpcaPorcentage + riskRatePercentage) / 100

func main() {
	start := time.Now()
	verifyWhereToInvest()
	end := time.Now()
	fmt.Println("Time taken:", end.Sub(start))
}

func verifyWhereToInvest() {
	ifxiList, err := b3.GetB3IfixData()
	if err != nil {
		log.Fatal("Error", err)
	}
	var wg sync.WaitGroup

	fundChan := make(chan cleandata.Fund)
	defer close(fundChan)
	wg.Add(len(ifxiList))

	for _, b := range ifxiList {
		go scrapperFundForIfix(b, &wg, fundChan)
	}

	var funds []cleandata.Fund
	go func() {
		for fund := range fundChan {
			funds = append(funds, fund)
		}
	}()
	wg.Wait()
	fmt.Println(funds)
	fundsCsv.CreateCSVFromFunds(funds, "funds.csv")
	fmt.Println("It is done")
}

func scrapperFundForIfix(ifix b3.Asset, wg *sync.WaitGroup, ch chan cleandata.Fund) {
	dividend, price := scrapper.GetDividendByCode(ifix.Code)
	defer wg.Done()

	generatedFund := cleandata.IfixFundData{
		Asset:    ifix,
		Price:    price,
		Dividend: dividend,
	}
	fund, err := generatedFund.CleanData(descountTax)
	if err != nil {
		log.Fatal("Error", err)
	}
	ch <- fund
}
