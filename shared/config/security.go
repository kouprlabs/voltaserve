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
	"strings"
)

type SecurityConfig struct {
	JWTSigningKey string
	CORSOrigins   []string
	APIKey        string
}

func ReadSecurity(config *SecurityConfig) {
	config.JWTSigningKey = os.Getenv("SECURITY_JWT_SIGNING_KEY")
	config.CORSOrigins = strings.Split(os.Getenv("SECURITY_CORS_ORIGINS"), ",")
	config.APIKey = os.Getenv("SECURITY_API_KEY")
}
