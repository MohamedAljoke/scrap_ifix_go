package scrapper

import (
	"fmt"

	"github.com/gocolly/colly"
)

var fundsUrl = "https://www.fundsexplorer.com.br/funds"

func GetDividendByCode(code string) (string, string) {
	c := colly.NewCollector()
	url := fmt.Sprintf("%s/%s", fundsUrl, code)

	var price string
	var dividendText string
	// headerTicker__content__price
	c.OnHTML("div.headerTicker__content__price p", func(h *colly.HTMLElement) {
		price = h.Text
	})

	c.OnHTML("div.indicators.wrapper div.indicators__box:nth-of-type(3) p:nth-of-type(2)", func(h *colly.HTMLElement) {
		dividendText = h.Text
	})
	c.Visit(url)

	return dividendText, price
}
