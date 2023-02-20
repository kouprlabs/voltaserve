package infra

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"voltaserve/config"

	"gopkg.in/gomail.v2"
	"sigs.k8s.io/yaml"

	"text/template"
)

type MessageParams struct {
	Subject string
}

type MailTemplate struct {
	dialer    *gomail.Dialer
	imageProc *ImageProcessor
	config    config.SmtpConfig
}

func NewMailTemplate() *MailTemplate {
	mt := new(MailTemplate)
	mt.config = config.GetConfig().Smtp
	mt.dialer = gomail.NewDialer(mt.config.Host, mt.config.Port, mt.config.Username, mt.config.Password)
	if mt.config.Secure {
		mt.dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	mt.imageProc = NewImageProcessor()
	return mt
}

func (mt *MailTemplate) Send(templateName string, address string, variables map[string]string) error {
	html, err := mt.GetText(filepath.FromSlash("templates/"+templateName+"/template.html"), variables)
	if err != nil {
		return err
	}
	text, err := mt.GetText(filepath.FromSlash("templates/"+templateName+"/template.txt"), variables)
	if err != nil {
		return err
	}
	params, err := mt.GetMessageParams(templateName)
	if err != nil {
		return err
	}
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf(`"%s" <%s>`, mt.config.SenderName, mt.config.SenderAddress))
	m.SetHeader("To", address)
	m.SetHeader("Subject", params.Subject)
	m.SetBody("text/plain ", text)
	m.AddAlternative("text/html", html)
	if err := mt.dialer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func (mt *MailTemplate) GetText(path string, variables map[string]string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	b, _ := io.ReadAll(f)
	html := string(b)
	tmpl, err := template.New("").Parse(html)
	if err != nil {
		return "", nil
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, variables)
	if err != nil {
		return "", nil
	}
	return buf.String(), nil
}

func (mt *MailTemplate) GetMessageParams(templateName string) (*MessageParams, error) {
	f, err := os.Open(filepath.FromSlash("templates/" + templateName + "/params.yml"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, _ := io.ReadAll(f)
	res := &MessageParams{}
	if err := yaml.Unmarshal(b, res); err != nil {
		return nil, err
	}
	return res, nil
}
