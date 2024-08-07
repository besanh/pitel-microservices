package mail

import (
	"bytes"
	"html/template"
	textTemplate "text/template"

	"gopkg.in/gomail.v2"
)

type MailConfig struct {
	Host     string
	User     string
	Password string
	Port     int
	From     string
	IsUseTLS bool
	IsUseSSL bool
}

type MailUtil struct {
	*MailConfig
}

var Client *MailUtil

func NewMailClient(config MailConfig) *MailUtil {
	return &MailUtil{
		MailConfig: &config,
	}
}

func (m *MailUtil) GetDialer() *gomail.Dialer {
	return gomail.NewPlainDialer(m.Host, m.Port, m.User, m.Password)
}

func (m *MailUtil) Send(subject, receiver string, message *gomail.Message) error {
	client := m.GetDialer()
	m.IsUseTLS = true
	m.IsUseSSL = false
	message.SetHeader("From", m.From)
	message.SetHeader("To", receiver)
	message.SetHeader("Subject", subject)
	return client.DialAndSend(message)
}

func ParseTemplate(templateFilepath string, data interface{}) (str string, err error) {
	t, err := template.ParseFiles(templateFilepath)
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return
	}
	str = buf.String()
	return
}

func ParseFromString(htmlString string, data interface{}) string {
	t, err := textTemplate.New("Email").Parse(htmlString)
	if err != nil {
		return ""
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return ""
	}
	return buf.String()
}
