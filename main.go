package main

import (
	"fmt"
	"log"
	"os"
	"scrap-ifix-go/b3"
	"scrap-ifix-go/cleanData"
	"scrap-ifix-go/fundsCsv"
	"scrap-ifix-go/scrapper"
	"sync"
	"time"
)

var ipcaRatePercentage = 6.1
var riskRatePercentage = 3.0

var discountTax = (ipcaRatePercentage + riskRatePercentage) / 100

func main() {
	start := time.Now()
	
	// Check for sequential flag
	sequential := len(os.Args) > 1 && os.Args[1] == "--sequential"
	
	if sequential {
		fmt.Println("Running in SEQUENTIAL mode...")
		processFundsSequential()
	} else {
		fmt.Println("Running in CONCURRENT mode...")
		processFunds()
	}
	
	end := time.Now()
	fmt.Println("Time taken:", end.Sub(start))
}

func processFunds() {
	ifixList, err := b3.GetB3IfixData()
	if err != nil {
		log.Fatalf("Failed to fetch IFIX data from B3: %v", err)
	}
	var wg sync.WaitGroup
	fundChan := make(chan cleanData.Fund)
	
	var funds []cleanData.Fund
	
	// Start collection goroutine
	collectionDone := make(chan struct{})
	go func() {
		defer close(collectionDone)
		for fund := range fundChan {
			funds = append(funds, fund)
		}
	}()
	
	wg.Add(len(ifixList))
	
	// Start scraping goroutines
	for _, asset := range ifixList {
		go scrapeFundData(asset, &wg, fundChan)
	}
	
	// Wait for all goroutines to finish, then close channel
	go func() {
		wg.Wait()
		close(fundChan)
	}()
	
	// Wait for collection to complete
	<-collectionDone
	fmt.Printf("Successfully processed %d funds\n", len(funds))

	if err := fundsCsv.CreateCSVFromFunds(funds, "funds.csv"); err != nil {
		log.Fatalf("Failed to create CSV file: %v", err)
	}
	fmt.Println("CSV file 'funds.csv' created successfully")
}

func processFundsSequential() {
	ifixList, err := b3.GetB3IfixData()
	if err != nil {
		log.Fatalf("Failed to fetch IFIX data from B3: %v", err)
	}
	
	var funds []cleanData.Fund
	
	// Process each fund sequentially
	for i, asset := range ifixList {
		fmt.Printf("Processing fund %d/%d: %s\n", i+1, len(ifixList), asset.Code)
		
		dividend, price, err := scrapper.GetDividendByCode(asset.Code)
		if err != nil {
			fmt.Printf("Skipping fund %s: scraping failed: %v\n", asset.Code, err)
			continue
		}

		generatedFund := cleanData.IfixFundData{
			Asset:    asset,
			Price:    price,
			Dividend: dividend,
		}
		
		fund, err := generatedFund.CleanData(discountTax)
		if err != nil {
			fmt.Printf("Skipping fund %s: %v\n", asset.Code, err)
			continue
		}
		
		funds = append(funds, fund)
	}
	
	fmt.Printf("Successfully processed %d funds\n", len(funds))
	
	if len(funds) == 0 {
		log.Printf("No funds were successfully processed")
		return
	}
	
	if err := fundsCsv.CreateCSVFromFunds(funds, "funds_sequential.csv"); err != nil {
		log.Fatalf("Failed to create CSV file: %v", err)
	}
	fmt.Println("CSV file 'funds_sequential.csv' created successfully")
}

func scrapeFundData(asset b3.Asset, wg *sync.WaitGroup, ch chan cleanData.Fund) {
	defer wg.Done()

	dividend, price, err := scrapper.GetDividendByCode(asset.Code)
	if err != nil {
		fmt.Printf("Skipping fund %s: scraping failed: %v\n", asset.Code, err)
		return
	}

	generatedFund := cleanData.IfixFundData{
		Asset:    asset,
		Price:    price,
		Dividend: dividend,
	}
	fund, err := generatedFund.CleanData(discountTax)
	if err != nil {
		fmt.Printf("Skipping fund %s: %v\n", asset.Code, err)
		return
	}
	ch <- fund
}
