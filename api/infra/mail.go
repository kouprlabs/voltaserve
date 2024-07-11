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
	"io/fs"
	"path/filepath"
	"text/template"

	"gopkg.in/gomail.v2"
	"sigs.k8s.io/yaml"

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/log"
	"github.com/kouprlabs/voltaserve/api/templates"
)

type dialer interface {
	DialAndSend(m ...*gomail.Message) error
}

type MessageParams struct {
	Subject string
}

type MailTemplate struct {
	dialer dialer
	config config.SMTPConfig
}

func NewMailTemplate(config config.SMTPConfig) *MailTemplate {
	return NewMailTemplateWithDialer(config, gomail.NewDialer(config.Host, config.Port, config.Username, config.Password))
}

func NewMailTemplateWithDialer(config config.SMTPConfig, dialer dialer) *MailTemplate {
	return &MailTemplate{
		config: config,
		dialer: dialer,
	}
}

func (mt *MailTemplate) Send(templateName string, address string, variables map[string]string) error {
	html, err := mt.getText(filepath.FromSlash(templateName+"/template.html"), variables)
	if err != nil {
		return err
	}
	text, err := mt.getText(filepath.FromSlash(templateName+"/template.txt"), variables)
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
		return fmt.Errorf("dial and sending mail: %w", err)
	}
	return nil
}

func (mt *MailTemplate) getText(path string, variables map[string]string) (string, error) {
	f, err := templates.FS.Open(path)
	if err != nil {
		return "", err
	}
	defer func(f fs.File) {
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
	f, err := templates.FS.Open(filepath.FromSlash(templateName + "/params.yml"))
	if err != nil {
		return nil, err
	}
	defer func(f fs.File) {
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
