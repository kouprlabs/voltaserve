// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package helper

import (
	"time"
)

func NewTimeString() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func TimeToString(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

func StringToTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

func StringToTimestamp(s string) int64 {
	return StringToTime(s).UnixMilli()
}

func TimeToTimestamp(t time.Time) int64 {
	return t.UnixMilli()
}

func ToUTCString(s *string) string {
	if s == nil {
		return ""
	}
	parsedTime, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		return ""
	}
	return parsedTime.Format(time.RFC1123)
}
