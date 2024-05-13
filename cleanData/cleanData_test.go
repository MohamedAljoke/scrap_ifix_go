package cleandata

import (
	"io"
	"os"
	"scrap-ifix-go/b3"
	"strings"
	"testing"
)

func Test_CleanDataSuccess(t *testing.T) {
	discountTax := 0.9
	ifixFund := &IfixFundData{
		Asset:    b3.Asset{},
		Price:    "R$ 100",
		Dividend: "% 10",
	}
	fundData, err := ifixFund.CleanData(discountTax)

	if err != nil {
		t.Errorf("expected error to be nil got %s", err.Error())
	}
	if fundData.Price != 100 {
		t.Errorf("expected price to be 100 got %f", fundData.Price)
	}
}

func Test_CleanData_ErrorParsing(t *testing.T) {
	ifixFund := &IfixFundData{
		Asset:    b3.Asset{},
		Price:    "invalid",
		Dividend: "invalid",
	}
	_, err := ifixFund.CleanData(0.1)

	if err == nil {
		t.Errorf("Expected error but got nil")
	}
}
func Test_CleanData_PrintError(t *testing.T) {
	oldOut := os.Stdout
	r, w, _ := os.Pipe()

	os.Stdout = w

	defer func() {
		os.Stdout = oldOut
	}()

	ifixFund := &IfixFundData{
		Price:    "invalid",
		Dividend: "invalid",
	}
	_, _ = ifixFund.CleanData(0.1)

	_ = w.Close()
	out, _ := io.ReadAll(r)
	expectedErrorMsg := "Error occured while parsing dividend"
	if !strings.Contains(string(out), expectedErrorMsg) {
		t.Errorf("intro text not correct got %s", string(out))
	}
}
func Test_CleanData_PrintPriceError(t *testing.T) {
	oldOut := os.Stdout
	r, w, _ := os.Pipe()

	os.Stdout = w

	defer func() {
		os.Stdout = oldOut
	}()

	ifixFund := &IfixFundData{
		Price:    "invalid",
		Dividend: "% 10",
	}
	_, _ = ifixFund.CleanData(0.1)

	_ = w.Close()
	out, _ := io.ReadAll(r)
	expectedErrorMsg := "Error occured while parsing price"
	if !strings.Contains(string(out), expectedErrorMsg) {
		t.Errorf("intro text not correct got %s", string(out))
	}
}
