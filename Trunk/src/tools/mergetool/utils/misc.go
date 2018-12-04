package utils

import (
	"encoding/hex"
	"strconv"
	"sync/atomic"

	"github.com/google/uuid"
)

var id = int32(5000)

func UUID() string {
	uid := uuid.New()
	var buf [32]byte
	hex.Encode(buf[:], uid[:])
	return string(buf[:])
}

func NewRoleId(s string) string {
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	old := int(i64)
	serverid := old - old/10000
	val := atomic.AddInt32(&id, 1)
	return strconv.FormatInt(int64(int(val)*10000+serverid), 10)
}
