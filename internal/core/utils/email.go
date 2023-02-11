package utils

import (
	"crypto/tls"
	"fmt"
	"git.friends.com/PetLand/UserService/v2/internal/core"
	"git.friends.com/PetLand/UserService/v2/internal/genErr"
	"gopkg.in/gomail.v2"
)

func SendEmail(email string, code string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", "petland23@mail.ru")
	msg.SetHeader("To", email)
	//fmt.Println("To:", email)
	msg.SetHeader("Subject", "Verification code")
	msg.SetBody("text/html", fmt.Sprintf("<b>%s</b>", code))
	//fmt.Printf("<b>%s</b>\n", code)

	n := gomail.NewDialer("smtp.mail.ru", 465, "petland23@mail.ru", "uCjsve57KRhjdik4FUBt")
	n.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		return genErr.NewError(err, core.ErrWrongType)
	}

	return nil
}
