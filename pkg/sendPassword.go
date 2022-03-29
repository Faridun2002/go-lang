package pkg

import (
	"fmt"
	gomail "gopkg.in/mail.v2"
	Type "regist/types"
)

func SendPassword(user Type.Rec_user) {
	mail := gomail.NewMessage()
	mail.SetHeader("From", "faridunjalolov1@gmail.com")
	mail.SetHeader("To", user.Email)
	mail.SetHeader("Subject", "Пароль для входа")
	mail.SetBody("text/html", "<html><body><b>Ваш логин:</b> " + user.Login + "<br><b>Ваш пароль:</b> " + user.Pass + "</body></html>")
	a := gomail.NewDialer("smtp.gmail.com",587,"faridunjalolov1@gmail.com","farn2212")
	if err := a.DialAndSend(mail); err != nil {
		fmt.Println()
		panic("Не удалось отправить пароль для восстановления на почту")
	}
}