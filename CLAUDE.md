# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Build and Run
- `go run main.go` - Run the main scraping application
- `go build` - Build the executable
- `go mod tidy` - Update dependencies

### Testing
- `go test ./...` - Run all tests across all packages
- `go test -v ./...` - Run tests with verbose output
- `go test ./b3` - Test specific package (b3, scrapper, cleanData, fundsCsv)
- `go test -cover ./...` - Run tests with coverage report

## Architecture Overview

This is a Go web scraping application that analyzes Brazilian real estate investment funds (REITs) from the IFIX index. The application follows a concurrent pipeline pattern:

### Core Flow
1. **Data Source**: Fetches IFIX fund list from B3 (Brazilian stock exchange) API
2. **Concurrent Scraping**: Uses goroutines and channels to scrape fund data in parallel from FundsExplorer
3. **Data Processing**: Cleans and calculates investment metrics (yield, max price based on discount rate)
4. **CSV Export**: Outputs results to CSV file for analysis

### Package Structure
- `b3/` - Fetches IFIX fund list from B3 API endpoint
- `scrapper/` - Web scraping using Colly to get dividend and price data from FundsExplorer
- `cleanData/` - Data processing and financial calculations (yield analysis, discount rate application)
- `fundsCsv/` - CSV file generation for output

### Key Technical Patterns
- **Concurrency**: Uses `sync.WaitGroup` and channels for parallel web scraping
- **Error Handling**: Consistent error propagation through the pipeline
- **Data Flow**: Channel-based communication between goroutines for fund data collection

### Configuration
- `taxaIpcaPorcentage = 6.1` - IPCA inflation rate
- `riskRatePercentage = 3.0` - Risk rate for calculations
- Combined discount rate used for maximum price calculations