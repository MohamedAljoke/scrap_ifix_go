package scrapper

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Scrapper_GetDividendSuccess(t *testing.T) {
	html := `<html>
				<head><title>Test Page</title></head>
				<body>
					<div class="headerTicker__content__price">
						<p>R$ 220,26</p>
					</div>
					<div class="indicators wrapper">
						<div class="indicators__box">
						</div>
						<div class="indicators__box">
						</div>
						<div class="indicators__box">
							<p>Test Dividend</p>
							<p>16,32 %</p>
						</div>
					</div>
				</body>
			</html>`

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, html)
	}))

	defer mockServer.Close()

	oldURL := fundsUrl

	fundsUrl = mockServer.URL

	defer func() { fundsUrl = oldURL }()
	dividend, price := GetDividendByCode("ABC11")
	expectedDividend := "16,32 %"
	expectedPrice := "R$ 220,26"

	if price != expectedPrice {
		t.Errorf("Unexpected price result, expected %s, got: %s", expectedPrice, price)
	}
	if dividend != expectedDividend {
		t.Errorf("Unexpected dividend result, expected %s, got: %s", expectedDividend, dividend)
	}
}
