package pkg

import (
	"fmt"
	"strconv"
	gomail "gopkg.in/mail.v2"
)

func Mail(Bytes int, Email string) {
	B := strconv.Itoa(Bytes)
	mail := gomail.NewMessage()
	mail.SetHeader("From", "faridunjalolov1@gmail.com")
	mail.SetHeader("To", Email)
	mail.SetHeader("Subject", "Пароль для подтверждения регистрации")
	mail.SetBody("text/html", "<body>Введите " + B + " этот пароль для подтверждение регистрации</body>")
	a := gomail.NewDialer("smtp.gmail.com",587,"faridunjalolov1@gmail.com","farn2212")
	if err := a.DialAndSend(mail); err != nil {
		fmt.Println()
		panic("Не удалось отправить пароль на почту")
	}
	fmt.Println(Bytes)
}