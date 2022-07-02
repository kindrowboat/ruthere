package ruthere

import (
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
)

const host = "127.0.0.1"
const port = "25"
const username = ""
const password = ""

var auth = smtp.PlainAuth("", username, password, host)

const addr = host + ":" + port

type localPostfixMailer struct {
	from mail.Address
}

func NewLocalPosixMailer(from mail.Address) localPostfixMailer {
	return localPostfixMailer{
		from: from,
	}
}

func (m localPostfixMailer) SendMail(toTouple []mail.Address, subject string, body string) (err error) {
	fromName := m.from.Name
	fromEmail := m.from.Address

	toAddresses := []string{}
	for _, address := range toTouple {
		toAddresses = append(toAddresses, address.String())
	}

	toHeader := strings.Join(toAddresses, ", ")
	from := mail.Address{
		Name:    fromName,
		Address: fromEmail,
	}
	fromHeader := from.String()
	subjectHeader := subject
	header := make(map[string]string)
	header["To"] = toHeader
	header["From"] = fromHeader
	header["Subject"] = subjectHeader
	header["Content-Type"] = `text/plain; charset="UTF-8"`
	msg := ""

	for k, v := range header {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	msg += "\r\n" + body
	bMsg := []byte(msg)
	// Send using local postfix service
	c, err := smtp.Dial(addr)

	if err != nil {
		return (err)
	}

	defer c.Close()
	if err = c.Mail(fromHeader); err != nil {
		return (err)
	}

	for _, to := range toTouple {
		if err = c.Rcpt(to.Address); err != nil {
			return (err)
		}
	}

	w, err := c.Data()
	if err != nil {
		return (err)
	}
	_, err = w.Write(bMsg)
	if err != nil {
		return (err)
	}

	err = w.Close()
	if err != nil {
		return (err)
	}

	err = c.Quit()
	// Or alternatively, send with remote service like Amazon SES
	// err = smtp.SendMail(addr, auth, fromEmail, toEmails, bMsg)
	// Handle response from local postfix or remote service
	if err != nil {
		return (err)
	}
	return
}
