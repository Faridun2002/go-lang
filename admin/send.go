package admin

import (
	"fmt"
	"net/smtp"
	"time"
	"io/ioutil"
	"path/filepath"
	"encoding/base64"
	"strings"
	"mime/multipart"
	"bytes"
	"net/http"
	pkg "regist/pkg"
)

type Sender struct {
	Auth smtp.Auth
}
type Message struct {
	To          []string
	Subject     string
	Body        string
	Attachments map[string][]byte
}

var (
	
	host 		= "smtp.gmail.com"
	username 	= "faridunjalolov1@gmail.com"
	password 	= "farn2212"
	portNumber 	= "587"
	Time_now time.Time
)

func Send() {
	sender := New()
	m := NewMessage("Test", "Body message.")
	m.To = []string{"faridunjalolov0@gmail.com"}
	m.AttachFile("templates/excel/" + Time_now.Format("2006.01.02") + ".xlsx")
	fmt.Println(sender.Send(m))
}

func New() *Sender {
	Auth := smtp.PlainAuth("", username, password, host)
	return &Sender{Auth}
}

func (s *Sender) Send(m *Message) error {
	return smtp.SendMail(fmt.Sprintf("%s:%s", host, portNumber), s.Auth, username, m.To, m.ToBytes())
}

func NewMessage(s, b string) *Message {
	return &Message{Subject: s, Body: b, Attachments: make(map[string][]byte)}
}

func (m *Message) AttachFile(src string) error {
	b, err := ioutil.ReadFile(src)
	pkg.ForError(err)
	_, fileName := filepath.Split(src)
	m.Attachments[fileName] = b
	return nil
}

func (m *Message) ToBytes() []byte {
	buf := bytes.NewBuffer(nil)
	withAttachments := len(m.Attachments) > 0
	buf.WriteString(fmt.Sprintf("Subject: %s\n", m.Subject))
	buf.WriteString(fmt.Sprintf("To: %s\n", strings.Join(m.To, ",")))
	buf.WriteString("MIME-Version: 1.0\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	if withAttachments {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	} else {
		buf.WriteString("Content-Type: text/plain; charset=utf-8\n")
	}
	buf.WriteString(m.Body)
	if withAttachments {
		for k, v := range m.Attachments {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v)))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", k))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}

		buf.WriteString("--")
	}
	return buf.Bytes()
}
