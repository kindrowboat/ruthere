package ruthere

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/mail"
	"time"
)

type httpChecker struct {
	mailer      localPostfixMailer
	toAddresses []mail.Address
	isUp        map[string]bool
	interval    time.Duration
}

func NewHttpChecker(fromAddress mail.Address, toAddresses []mail.Address, sitesToCheck []string, interval time.Duration) httpChecker {
	checker := httpChecker{
		mailer:      NewLocalPosixMailer(fromAddress),
		toAddresses: toAddresses,
		isUp:        map[string]bool{},
		interval:    interval,
	}
	for _, site := range sitesToCheck {
		checker.isUp[site] = true
	}
	return checker
}

func (c httpChecker) Run() {
	for true {
		for key := range c.isUp {
			log.Printf("Checking %s", key)
			response, err := http.Get(key)
			if err != nil {
				c.recordAndReportDown(key, fmt.Sprintf("Could not connect to server: %s", err.Error()))
				continue
			}
			if response.StatusCode < 200 || response.StatusCode > 299 {
				body, err := io.ReadAll(response.Body)
				response.Body.Close()
				if err != nil {
					c.recordAndReportDown(key, fmt.Sprintf("Could not read response: %s", err.Error()))
					continue
				}
				message := fmt.Sprintf("Returned status %d with body:\n%s", response.StatusCode, body)
				c.recordAndReportDown(key, message)
				continue
			} else {
				c.isUp[key] = true
			}
		}
		log.Print("Checks finished, sleeping... zzz...")
		time.Sleep(c.interval)
	}
	return
}

func (c *httpChecker) recordAndReportDown(site string, message string) (err error) {
	if c.isUp[site] {
		c.isUp[site] = false
		log.Printf("%s is down, sending notification", site)
		subject := fmt.Sprintf("ALERT: %s is DOWN", site)
		err = c.mailer.SendMail(c.toAddresses, subject, message)
		if err != nil {
			log.Printf("Could not send alert email: %s", err.Error())
		}
	} else {
		log.Printf("%s is still down", site)
	}
	return
}
