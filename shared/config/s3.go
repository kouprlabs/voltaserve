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

type S3Config struct {
	URL       string
	AccessKey string
	SecretKey string
	Region    string
	Bucket    string
	Secure    bool
}

func ReadS3(config *S3Config) {
	config.URL = os.Getenv("S3_URL")
	config.AccessKey = os.Getenv("S3_ACCESS_KEY")
	config.SecretKey = os.Getenv("S3_SECRET_KEY")
	config.Region = os.Getenv("S3_REGION")
	config.Bucket = os.Getenv("S3_BUCKET")
	if len(os.Getenv("S3_SECURE")) > 0 {
		v, err := strconv.ParseBool(os.Getenv("S3_SECURE"))
		if err != nil {
			panic(err)
		}
		config.Secure = v
	}
}
