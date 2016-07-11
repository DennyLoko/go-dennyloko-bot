package currency

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/Sirupsen/logrus"
)

var cache map[string]*Exchange

type Exchange struct {
	Rate float64
	Time time.Time
}

// Currex represents one currency exchange action
type Currex struct {
	From   string
	To     string
	Amount float64
	Log    *logrus.Logger
}

func init() {
	cache = make(map[string]*Exchange, 0)
}

// Convert the currency and returns the result
func (c *Currex) Convert() (success bool, from float64, to float64, err error) {
	key := fmt.Sprintf("%s%s", c.From, c.To)

	if exc, ok := cache[key]; ok && (time.Now()).Sub(exc.Time) < (5*time.Minute) {
		c.Log.Debug("Returning from cache...")
		from = c.Amount
		to = exc.Rate * c.Amount
		success = true

		return
	}

	url := fmt.Sprintf("https://www.google.com/finance/converter?a=%d&from=%s&to=%s", 1, c.From, c.To)

	c.Log.Debugf("Calling URL: %s", url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return
	}

	c.Log.Debug("Parsing document...")
	r := strings.TrimSpace(doc.Find("#currency_converter_result").Text())
	split := strings.Split(r, " = ")
	c.Log.Debugf("We got: %v", split)

	if len(split) == 2 {
		var q []string

		rx := regexp.MustCompile("(\\d*(\\.?\\d*))")

		q = rx.FindStringSubmatch(split[1])
		to, err = strconv.ParseFloat(q[0], 64)

		c.Log.Debug("Creating exchange object...")
		exc := &Exchange{
			Rate: to,
			Time: time.Now(),
		}

		c.Log.Debug("Creating cache entry...")
		cache[key] = exc

		from = c.Amount
		to = to * c.Amount
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
