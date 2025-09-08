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

type MailType string

const SUPPORT_EMAIL = "support@support.com"
const (
	MailTypeSignUp        MailType = "sign_up"
	MailTypeSendOTP       MailType = "send_otp"
	MailTypeURLRegistered MailType = "url_registered"
)

type MailOptions interface{}

type AppConfigOptions struct {
	APP_NAME string
}

type SignUpMailOptions struct {
	AppConfigOptions
	USER_EMAIL    string
	LOGIN_URL     string
	SUPPORT_EMAIL string
	IMG_BASE_URL  template.URL
}

type SendOTPMailOptions struct {
	AppConfigOptions
	ACTION_TYPE   string
	USER_EMAIL    string
	OTP_CODE      string
	SUPPORT_EMAIL string
	IMG_BASE_URL  template.URL
}

type URLRegisteredMailOptions struct {
	AppConfigOptions
	USER_EMAIL    string
	SHORT_URL     string
	ORIGINAL_URL  string
	SUPPORT_EMAIL string
	IMG_BASE_URL  template.URL
}

//go:embed template/sign-up-success.html
var signUpTemplate string

//go:embed template/send-otp.html
var sendOtpTemplate string

//go:embed template/url-registered.html
var urlRegisteredTemplate string

//go:embed assets/logo.png
var logoImg []byte

var mailTemplates = map[MailType]string{
	MailTypeSignUp:        signUpTemplate,
	MailTypeSendOTP:       sendOtpTemplate,
	MailTypeURLRegistered: urlRegisteredTemplate,
}

func logoBase64() string {
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(logoImg)
}

func sendMail(mailType MailType, opts MailOptions, toEmail string, subject string) error {
	tmplStr := mailTemplates[mailType]

	switch mailType {
	case MailTypeSignUp:
		signUpOpts, ok := opts.(SignUpMailOptions)
		if !ok {
			return fmt.Errorf("opts must be SignUpMailOptions for MailTypeSignUp")
		}

		signUpOpts.APP_NAME = config.Config.APP.Name
		opts = signUpOpts
	case MailTypeSendOTP:
		loginOtpOpts, ok := opts.(SendOTPMailOptions)
		if !ok {
			return fmt.Errorf("opts must be SendOTPMailOptions for MailTypeSendOTP")
		}
		loginOtpOpts.APP_NAME = config.Config.APP.Name
		opts = loginOtpOpts
	case MailTypeURLRegistered:
		urlRegisteredOpts, ok := opts.(URLRegisteredMailOptions)
		if !ok {
			return fmt.Errorf("opts must be URLRegisteredMailOptions for MailTypeURLRegistered")
		}
		urlRegisteredOpts.APP_NAME = config.Config.APP.Name
		opts = urlRegisteredOpts
	default:
		return fmt.Errorf("unknown mail type: %s", mailType)
	}

	t, parseErr := template.New(string(mailType)).Parse(tmplStr)
	if parseErr != nil {
		return parseErr
	}

	var buf bytes.Buffer
	executeErr := t.Execute(&buf, opts)
	if executeErr != nil {
		return executeErr
	}

	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", config.Config.APP.Name, SUPPORT_EMAIL)
	e.To = []string{toEmail}
	e.Subject = subject
	e.HTML = buf.Bytes()
	e.Text = buf.Bytes()

	smtpUrl := fmt.Sprintf("%s:%s", config.Config.SMTP.Host, config.Config.SMTP.Port)
	err := e.Send(smtpUrl,
		smtp.PlainAuth("", config.Config.SMTP.Username, config.Config.SMTP.Password, config.Config.SMTP.Host))
	if err != nil {
		return err
	}

	return nil
}

func SendSignedUpUserMail(u model.User) error {
	data := SignUpMailOptions{
		USER_EMAIL:    u.Email,
		LOGIN_URL:     fmt.Sprintf("http://%s:%s/user/login", config.Config.APP.Host, config.Config.APP.Port),
		SUPPORT_EMAIL: SUPPORT_EMAIL,
		IMG_BASE_URL:  template.URL(logoBase64()),
		AppConfigOptions: AppConfigOptions{
			APP_NAME: config.Config.APP.Name,
		},
	}

	return sendMail(MailTypeSignUp, data, u.Email, "Welcome to "+config.Config.APP.Name)
}

func SendOtpUserMail(o model.Otp) error {
	data := SendOTPMailOptions{
		USER_EMAIL:    o.Key,
		SUPPORT_EMAIL: SUPPORT_EMAIL,
		IMG_BASE_URL:  template.URL(logoBase64()),
		ACTION_TYPE:   string(o.Action),
		OTP_CODE:      o.OtpCode,
		AppConfigOptions: AppConfigOptions{
			APP_NAME: config.Config.APP.Name,
		},
	}

	subject := fmt.Sprintf("Your OTP to %s for %s", o.Action, config.Config.APP.Name)
	return sendMail(MailTypeSendOTP, data, o.Key, subject)
}

func SendShortUrlUserMail(u model.Url) error {
	user, userErr := model.GetUserById(u.UserID)
	if userErr != nil {
		return userErr
	}

	data := URLRegisteredMailOptions{
		USER_EMAIL:    user.Email,
		SHORT_URL:     utils.GetShortUrl(u.Code),
		ORIGINAL_URL:  u.Url,
		SUPPORT_EMAIL: SUPPORT_EMAIL,
		IMG_BASE_URL:  template.URL(logoBase64()),
		AppConfigOptions: AppConfigOptions{
			APP_NAME: config.Config.APP.Name,
		},
	}

	return sendMail(MailTypeURLRegistered, data, user.Email, "Your shortened URL is ready")
}
