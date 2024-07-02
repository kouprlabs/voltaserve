package helper

import (
	"time"

	"github.com/google/uuid"
	"github.com/speps/go-hashids/v2"
)

func NewID() string {
	hd := hashids.NewData()
	hd.Salt = uuid.NewString()
	h, err := hashids.NewWithData(hd)
	if err != nil {
		panic(err)
	}
	id, err := h.EncodeInt64([]int64{time.Now().UTC().UnixNano()})
	if err != nil {
		panic(err)
	}
	return id
}
