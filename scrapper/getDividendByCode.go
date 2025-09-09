package scrapper

import (
	"fmt"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

var fundsUrl = "https://www.fundsexplorer.com.br/funds"

// Rate limiter to prevent "Too Many Requests" errors
var rateLimiter = time.NewTicker(500 * time.Millisecond) // 2 requests per second
var rateLimiterMutex sync.Mutex

func GetDividendByCode(code string) (string, string, error) {
	// Rate limiting - wait for the ticker
	rateLimiterMutex.Lock()
	<-rateLimiter.C
	rateLimiterMutex.Unlock()
	c := colly.NewCollector()
	
	// Set timeout for requests
	c.SetRequestTimeout(30 * time.Second)
	
	// Add error handling
	var scrapeError error
	c.OnError(func(r *colly.Response, err error) {
		scrapeError = fmt.Errorf("failed to scrape fund %s: %w", code, err)
	})

	url := fmt.Sprintf("%s/%s", fundsUrl, code)

	var price string
	var dividendText string
	
	c.OnHTML("div.headerTicker__content__price p", func(h *colly.HTMLElement) {
		price = h.Text
	})

	c.OnHTML("div.indicators.wrapper div.indicators__box:nth-of-type(3) p:nth-of-type(2)", func(h *colly.HTMLElement) {
		dividendText = h.Text
	})
	
	if err := c.Visit(url); err != nil {
		return "", "", fmt.Errorf("failed to visit %s: %w", url, err)
	}

	if scrapeError != nil {
		return "", "", scrapeError
	}

	return dividendText, price, nil
}
