package bot

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Sirupsen/logrus"
)

// Currex represents one currency exchange action
type Currex struct {
	From   string
	To     string
	Amount float64
	log    *logrus.Logger
}

// Convert the currency and returns the result
func (c *Currex) Convert() (success bool, from float64, to float64, err error) {
	url := fmt.Sprintf("https://www.google.com/finance/converter?a=%f&from=%s&to=%s", c.Amount, c.From, c.To)
	c.log.Debugf("Calling URL: %s", url)

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return
	}

	r := strings.TrimSpace(doc.Find("#currency_converter_result").Text())
	split := strings.Split(r, " = ")

	if len(split) == 2 {
		var q []string

		rx := regexp.MustCompile("(\\d*(\\.?\\d*))")

		q = rx.FindStringSubmatch(split[0])
		from, err = strconv.ParseFloat(q[0], 64)

		q = rx.FindStringSubmatch(split[1])
		to, err = strconv.ParseFloat(q[0], 64)

		success = true
	} else if len(split) > 2 {
		err = fmt.Errorf("The API returned a wrong value: %s", r)
	}

	return
}

// Validate the currency
func (c *Currex) Validate(currency string) error {
	if len(currency) < 3 {
		return errors.New("Please check your currencies")
	}

	return nil
}
