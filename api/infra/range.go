package infra

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

type RangeInterval struct {
	Start int64
	End   int64
}

func NewRangeInterval(header string) *RangeInterval {
	ri := new(RangeInterval)
	if header != "" {
		parts := strings.Split(header, "=")
		if len(parts) == 2 {
			ranges := strings.Split(parts[1], "-")
			ri.Start, _ = strconv.ParseInt(ranges[0], 10, 64)
			if len(ranges) > 1 && ranges[1] != "" {
				ri.End, _ = strconv.ParseInt(ranges[1], 10, 64)
			}
		}
	}
	return ri
}

func (ri RangeInterval) ApplyToMinIOGetObjectOptions(opts *minio.GetObjectOptions) {
	if ri.End != 0 {
		opts.SetRange(ri.Start, ri.End)
	} else {
		opts.SetRange(ri.Start, 0)
	}
}

func (ri RangeInterval) ApplyToFiberContext(partialSize int64, totalSize int64, ctx *fiber.Ctx) {
	ctx.Set("Accept-Ranges", "bytes")
	if ri.Start != 0 || ri.End != 0 {
		if ri.End != 0 {
			ctx.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", ri.Start, ri.End, totalSize))
			ctx.Set("Content-Length", fmt.Sprintf("%d", partialSize))
		} else {
			ctx.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", ri.Start, totalSize-1, totalSize))
			ctx.Set("Content-Length", fmt.Sprintf("%d", partialSize))
		}
	} else {
		ctx.Set("Content-Range", fmt.Sprintf("bytes 0-%d/%d", totalSize-1, totalSize))
		ctx.Set("Content-Length", fmt.Sprintf("%d", totalSize))
	}
}
