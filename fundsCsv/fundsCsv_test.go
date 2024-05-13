package fundsCsv

import (
	"encoding/csv"
	"os"
	"reflect"
	cleandata "scrap-ifix-go/cleanData"
	"testing"
)

func Test_FundsCsv_Success(t *testing.T) {

	mockFunds := []cleandata.Fund{
		{
			Code:     "ABC",
			Yield:    "3%",
			Price:    100.50,
			MaxPrice: 120.75,
		},
	}
	file, err := os.CreateTemp("", "test_fund.csv")
	if err != nil {
		t.Fatalf("error creating temporary file %v", err.Error())
	}
	defer os.Remove(file.Name())

	err = CreateCSVFromFunds(mockFunds, file.Name())
	if err != nil {
		t.Fatalf("error writing to csv %v", err.Error())
	}
	csvFile, err := os.Open(file.Name())
	if err != nil {
		t.Fatalf("error opening csv file %v", err.Error())
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("error reading CSV file: %v", err)
	}
	expectedHeader := []string{"Code", "Yield", "Price", "MaxPrice"}
	if !reflect.DeepEqual(records[0], expectedHeader) {
		t.Errorf("header mismatch, expected: %v, got: %v", expectedHeader, records[0])
	}
	expectedRecord := []string{"ABC", "3%", "100.5", "120.75"}
	if !reflect.DeepEqual(records[1], expectedRecord) {
		t.Errorf("record mismatch, expected: %v, got: %v", expectedRecord, records[1])
	}
}

func Test_FundsCsv_Throws(t *testing.T) {
	mockFunds := []cleandata.Fund{
		{
			Code:     "ABC",
			Yield:    "3%",
			Price:    100.50,
			MaxPrice: 120.75,
		},
	}

	err := CreateCSVFromFunds(mockFunds, "/path/to/a/directory")
	if err == nil {
		t.Errorf("error writing to CSV: %v", err)
	}
}
