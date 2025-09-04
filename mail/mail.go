package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/smtp"

	_ "embed"

	"example.com/url-shortner/config"
	"example.com/url-shortner/model"
	"example.com/url-shortner/utils"
	"github.com/jordan-wright/email"
)

//go:embed template/sign-up-success.html
var signUpTemplate string

//go:embed assets/logo.png
var logoImg []byte

func logoBase64() string {
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(logoImg)
}

func SendSignedUpUserMail(u model.User) error {

	type SignUpUserMailOpt struct {
		APP_NAME      string
		USER_EMAIL    string
		LOGIN_URL     string
		SUPPORT_EMAIL string
		IMG_BASE_URL  template.URL
	}

	t, parseErr := template.New("welcome").Parse(string(signUpTemplate))
	if parseErr != nil {
		return parseErr
	}

	data := SignUpUserMailOpt{
		APP_NAME:      config.Config.App.Name,
		USER_EMAIL:    u.Email,
		LOGIN_URL:     "http://localhost/logiin",
		SUPPORT_EMAIL: "support@support.com",
		IMG_BASE_URL:  template.URL(logoBase64()),
	}

	var buf bytes.Buffer
	executeErr := t.Execute(&buf, data)
	if executeErr != nil {
		return executeErr
	}

	parsedHTML := buf.String()
	utils.Log.Info("User signed up mail to send: ", u)

	subject := fmt.Sprintf("Welcome aboard - %s !", config.Config.App.Name)

	e := email.NewEmail()
	e.From = "Your Name <your@email.com>"
	e.To = []string{data.USER_EMAIL}
	e.Subject = subject
	e.HTML = []byte(parsedHTML)
	e.Text = []byte(parsedHTML)

	smtpUrl := fmt.Sprintf("%s:%s", config.Config.SMTP.Host, config.Config.SMTP.Port)

	err := e.Send(smtpUrl,
		smtp.PlainAuth("", config.Config.SMTP.Username, config.Config.SMTP.Password, config.Config.SMTP.Host))

	if err != nil {
		panic(err)
	}

	return nil
}
