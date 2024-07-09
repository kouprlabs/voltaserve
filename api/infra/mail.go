// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package infra

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"voltaserve/config"
	"voltaserve/log"

	"gopkg.in/gomail.v2"
	"sigs.k8s.io/yaml"

	"text/template"
)

type MessageParams struct {
	Subject string
}

type MailTemplate struct {
	dialer *gomail.Dialer
	config config.SMTPConfig
}

func NewMailTemplate() *MailTemplate {
	mt := new(MailTemplate)
	mt.config = config.GetConfig().SMTP
	mt.dialer = gomail.NewDialer(mt.config.Host, mt.config.Port, mt.config.Username, mt.config.Password)
	return mt
}

func (mt *MailTemplate) Send(templateName string, address string, variables map[string]string) error {
	html, err := mt.getText(filepath.FromSlash("templates/"+templateName+"/template.html"), variables)
	if err != nil {
		return err
	}
	text, err := mt.getText(filepath.FromSlash("templates/"+templateName+"/template.txt"), variables)
	if err != nil {
		return err
	}
	params, err := mt.getMessageParams(templateName)
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

func (mt *MailTemplate) getText(path string, variables map[string]string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			log.GetLogger().Error(err)
		}
	}(f)
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

func (mt *MailTemplate) getMessageParams(templateName string) (*MessageParams, error) {
	f, err := os.Open(filepath.FromSlash("templates/" + templateName + "/params.yml"))
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			log.GetLogger().Error(err)
		}
	}(f)
	b, _ := io.ReadAll(f)
	res := &MessageParams{}
	if err := yaml.Unmarshal(b, res); err != nil {
		return nil, err
	}
	return res, nil
}
