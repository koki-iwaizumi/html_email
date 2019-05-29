package main

import (
	"./model"
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"html/template"
	_ "math/rand"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
	"time"
)

const (
	EMAIL_HOST                = "*"
	EMAIL_PORT                = "*"
	EMAIL_USER                = "*"
	EMAIL_PASSWORD            = "*"
	EMAIL_FROMNAME            = "*"
	EMAIL_FROMADDRESS         = "*"
	EMAIL_SUBJECT             = "IoTデータロガーのご紹介！"
	EMAIL_PATH                = "view/email.html"
	EMAIL_TARGET_STATUS       = "未送信"
	EMAIL_TARGET_STATUS_AFTER = "メール送信済み"
)

type EmailData struct {
	Person *model.Person
	Year   string
	Month  string
	Day    string
	Title  string
}

func main() {

	err := emailMain()
	if err != nil {
		fmt.Println(err)
		fmt.Println("FINISH ERROR")
	} else {
		fmt.Println("FINISH SUCCESS")
	}

	// Y押したら終了
	for {
		fmt.Println("PRESS Y")
		p := StrStdin()
		if p == "Y" {
			break
		}
	}

}

func emailMain() (err error) {

	// データベース接続
	err = model.Connect()
	if err != nil {
		fmt.Println("model.Connect() ERROR")
		return err
	} else {
		fmt.Println("model.Connect() SUCCESS")
	}

	// データベース取得
	persons := *model.Personrepo.FindByJinto(EMAIL_TARGET_STATUS)

	var person_key = 0
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// メールテンプレート取得
	email_tmp := template.Must(template.ParseFiles(EMAIL_PATH))
	t := time.Now()
	t.Format("2006.01.02")

	data := EmailData{
		Year:  t.Format("2006"),
		Month: t.Format("01"),
		Day:   t.Format("02"),
		Title: EMAIL_SUBJECT,
	}

	for {
		select {
		case <-ticker.C:
			fmt.Println(persons[person_key])

			// テンプレート取得
			var tmp bytes.Buffer
			data.Person = &persons[person_key]
			err = email_tmp.Execute(&tmp, data)
			if err != nil {
				fmt.Println("EMAIL TMP ERROR : person_id=" + string(persons[person_key].Id))
				return err
			} else {
				fmt.Println("EMAIL TMP SUCCESS")
			}
			body := createBody(&tmp, EMAIL_SUBJECT, persons[person_key].Email, persons[person_key].Name)

			// メール送信
			err = sendEmail(&body, persons[person_key].Email)
			if err != nil {
				fmt.Println("EMAIL SEND ERROR : person_id=" + string(persons[person_key].Id))
				return err
			} else {
				fmt.Println("EMAIL SEND SUCCESS")

				persons[person_key].Jinto = EMAIL_TARGET_STATUS_AFTER

				err = model.Personrepo.Update(&persons[person_key])
				if err != nil {
					fmt.Println("EMAIL SAVE ERROR : person_id=" + string(persons[person_key].Id))
					return err
				}
			}

			fmt.Printf("EMAIL NOW -> %v\n", time.Now())

			person_key++

			if len(persons) == person_key {
				return err
			}
		}
	}
}

func createBody(tmp *bytes.Buffer, subject string, tomail string, toname string) (body bytes.Buffer) {

	to := mail.Address{toname, tomail}
	from := mail.Address{EMAIL_FROMNAME, EMAIL_FROMADDRESS}

	body.WriteString("To: " + to.String() + "\r\n")
	body.WriteString("From: " + from.String() + "\r\n")
	body.WriteString(encodeSubject(subject))
	body.WriteString("Mime-Version: 1.0\r\n")

	body.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	body.WriteString("Content-Transfer-Encoding: base64\r\n")
	body.WriteString("\r\n")
	body.WriteString(add76crlf(base64.StdEncoding.EncodeToString(tmp.Bytes())))
	body.WriteString("\r\n")
	body.WriteString("\r\n")

	return body
}

func sendEmail(body *bytes.Buffer, tomail string) (err error) {

	//SMTPサーバー
	host, _, _ := net.SplitHostPort(EMAIL_HOST + ":" + EMAIL_PORT)
	auth := smtp.PlainAuth("", EMAIL_USER, EMAIL_PASSWORD, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", EMAIL_HOST+":"+EMAIL_PORT, tlsconfig)

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer c.Close()

	// Auth
	err = c.Auth(auth)
	if err != nil {
		return err
	}

	err = c.Mail(EMAIL_FROMADDRESS)
	if err != nil {
		return err
	}

	err = c.Rcpt(tomail)
	if err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = body.WriteTo(w)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	err = c.Quit()
	if err != nil {
		return err
	}

	return err
}

// 76バイト毎にCRLFを挿入する
func add76crlf(msg string) string {
	var buffer bytes.Buffer
	for k, c := range strings.Split(msg, "") {
		buffer.WriteString(c)
		if k%76 == 75 {
			buffer.WriteString("\r\n")
		}
	}
	return buffer.String()
}

// UTF8文字列を指定文字数で分割
func utf8Split(utf8string string, length int) []string {
	resultString := []string{}
	var buffer bytes.Buffer
	for k, c := range strings.Split(utf8string, "") {
		buffer.WriteString(c)
		if k%length == length-1 {
			resultString = append(resultString, buffer.String())
			buffer.Reset()
		}
	}
	if buffer.Len() > 0 {
		resultString = append(resultString, buffer.String())
	}
	return resultString
}

// サブジェクトをMIMEエンコードする
func encodeSubject(subject string) string {
	var buffer bytes.Buffer
	buffer.WriteString("Subject:")
	for _, line := range utf8Split(subject, 13) {
		buffer.WriteString(" =?utf-8?B?")
		buffer.WriteString(base64.StdEncoding.EncodeToString([]byte(line)))
		buffer.WriteString("?=\r\n")
	}
	return buffer.String()
}

// 標準入力
func StrStdin() (stringInput string) {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	stringInput = scanner.Text()

	stringInput = strings.TrimSpace(stringInput)
	return
}
