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
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

const (
	MaxRangeSize     = 250 * 1024 * 1024 // 250 MB
	DefaultChunkSize = 250 * 1024 * 1024 // Default size to serve for open-ended ranges
)

type RangeInterval struct {
	Start    int64
	End      int64
	FileSize int64
}

func NewRangeInterval(header string, fileSize int64) *RangeInterval {
	ri := &RangeInterval{FileSize: fileSize}
	if header != "" {
		parts := strings.Split(header, "=")
		if len(parts) == 2 {
			ranges := strings.Split(parts[1], "-")
			ri.Start, _ = strconv.ParseInt(ranges[0], 10, 64)
			if len(ranges) > 1 && ranges[1] != "" {
				ri.End, _ = strconv.ParseInt(ranges[1], 10, 64)
			} else {
				// Indicates an open-ended range
				ri.End = 0
			}
		}
	}
	ri.adjustRange()
	return ri
}

func (ri *RangeInterval) adjustRange() {
	/* For open-ended ranges, serve 250 MB */
	if ri.End == 0 {
		ri.End = ri.Start + DefaultChunkSize - 1
	}
	/* Ensure the end is within file size limits */
	if ri.End >= ri.FileSize {
		ri.End = ri.FileSize - 1
	}
	/* Ensure the range size is within limits for closed ranges */
	if ri.End-ri.Start+1 > MaxRangeSize && ri.End != 0 {
		ri.End = ri.Start + MaxRangeSize - 1
	}
}

func (ri RangeInterval) ApplyToMinIOGetObjectOptions(opts *minio.GetObjectOptions) {
	opts.SetRange(ri.Start, ri.End)
}

func (ri RangeInterval) ApplyToFiberContext(ctx *fiber.Ctx) {
	ctx.Set("Accept-Ranges", "bytes")
	ctx.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", ri.Start, ri.End, ri.FileSize))
	ctx.Set("Content-Length", fmt.Sprintf("%d", ri.End-ri.Start+1))
}
