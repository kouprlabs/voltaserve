// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package infra_test

import (
	_ "embed"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/gomail.v2"

	"github.com/kouprlabs/voltaserve/shared/config"

	"github.com/kouprlabs/voltaserve/api/infra"
)

type DialMock struct {
	Err  error
	Body string
}

func (d *DialMock) DialAndSend(m ...*gomail.Message) error {
	var body strings.Builder

	for _, message := range m {
		_, err := message.WriteTo(&body)
		if err != nil {
			return fmt.Errorf("write to: %w", err)
		}
	}

	d.Body = body.String()

	return d.Err
}

//go:embed fixtures/templates/join-organization.eml
var joinOrganization string

//go:embed fixtures/templates/sign-up-and-join-organization.eml
var signupAndJoinOrganization string

type MailTemplateSuite struct {
	suite.Suite
}

func TestMailTemplateSuite(t *testing.T) {
	// Avoid the mock being instantiated, because here we are testing the real implementation
	t.Setenv("TEST", "")
	suite.Run(t, new(MailTemplateSuite))
}

func (s *MailTemplateSuite) TestSend() {
	tests := map[string]struct {
		TemplateName string
		Address      string
		Variables    map[string]string
		ExpectedBody string
	}{
		"join-organization": {
			TemplateName: "join-organization",
			Address:      `"Someone" <someone@example.com>`,
			Variables: map[string]string{
				"USER_FULL_NAME":    "Someone",
				"ORGANIZATION_NAME": "ACME",
				"UI_URL":            "example.com",
			},
			ExpectedBody: joinOrganization,
		},
		"signup-and-join-organization": {
			TemplateName: "signup-and-join-organization",
			Address:      `"Someone" <someone@example.com>`,
			Variables: map[string]string{
				"USER_FULL_NAME":    "Someone",
				"ORGANIZATION_NAME": "ACME",
				"UI_URL":            "example.com",
			},
			ExpectedBody: signupAndJoinOrganization,
		},
	}

	for name, tc := range tests {
		s.Run(name, func() {
			dialMock := &DialMock{}
			mt := infra.NewMailTemplateWithDialer(config.SMTPConfig{
				Host:          "localhost",
				SenderName:    "Voltaserve",
				SenderAddress: "voltaserve@example.com",
			}, dialMock, false)

			// gomail is non-deterministic in its headers, so we'll brute force our expected body.
			s.EventuallyWithT(func(t *assert.CollectT) {
				err := mt.Send(tc.TemplateName, tc.Address, tc.Variables)
				require.NoError(t, err)

				simplifiedBody := regexp.MustCompile("boundary=.+").ReplaceAllString(dialMock.Body, "boundary=XXX")
				simplifiedBody = regexp.MustCompile("--.+(|--)").ReplaceAllString(simplifiedBody, "--XXX$1")
				simplifiedBody = regexp.MustCompile("Date: .+").ReplaceAllString(simplifiedBody, "Date: Now")
				simplifiedBody = strings.ReplaceAll(simplifiedBody, "\r\n", "\n")
				assert.Equal(t, tc.ExpectedBody, simplifiedBody)
			}, 1*time.Second, 1)
		})
	}
}
