// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package config

import (
	"os"
	"strconv"
)

type SMTPConfig struct {
	Host          string
	Port          int
	Secure        bool
	Username      string
	Password      string
	SenderAddress string
	SenderName    string
}

func ReadSMTP(config *SMTPConfig) {
	config.Host = os.Getenv("SMTP_HOST")
	if len(os.Getenv("SMTP_PORT")) > 0 {
		v, err := strconv.ParseInt(os.Getenv("SMTP_PORT"), 10, 32)
		if err != nil {
			panic(err)
		}
		config.Port = int(v)
	}
	if len(os.Getenv("SMTP_SECURE")) > 0 {
		v, err := strconv.ParseBool(os.Getenv("SMTP_SECURE"))
		if err != nil {
			panic(err)
		}
		config.Secure = v
	}
	config.Username = os.Getenv("SMTP_USERNAME")
	config.Password = os.Getenv("SMTP_PASSWORD")
	config.SenderAddress = os.Getenv("SMTP_SENDER_ADDRESS")
	config.SenderName = os.Getenv("SMTP_SENDER_NAME")
}
