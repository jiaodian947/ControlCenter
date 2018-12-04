package transaction

import (
	"time"

	quicklz "github.com/dgryski/go-quicklz"
)

func Decompress(src []byte) (data []byte, err error) {
	data, err = quicklz.Decompress(src)
	return
}

type BaseTransaction struct {
	url        string
	identity   string
	connId     int64
	serverId   int
	gameId     int
	platform   string
	err        error
	errcode    int
	verifytime time.Time
}

func (v *BaseTransaction) Platform() string {
	return v.platform
}

func (v *BaseTransaction) SetUrl(url string) {
	v.url = url
}

func (v *BaseTransaction) Url() string {
	return v.url
}

func (v *BaseTransaction) SetIdentity(identity string) {
	v.identity = identity
}

func (v *BaseTransaction) Identity() string {
	return v.identity
}

func (v *BaseTransaction) SetConnId(connid int64) {
	v.connId = connid
}

func (v *BaseTransaction) ConnId() int64 {
	return v.connId
}

func (v *BaseTransaction) SetServerId(id int) {
	v.serverId = id
}

func (v *BaseTransaction) ServerId() int {
	return v.serverId
}

func (v *BaseTransaction) SetGameId(id int) {
	v.gameId = id
}

func (v *BaseTransaction) GameId() int {
	return v.gameId
}

func (v *BaseTransaction) VerifyTime() time.Time {
	return v.verifytime
}

func (v *BaseTransaction) Err() error {
	return v.err
}

func (v *BaseTransaction) ErrCode() int {
	return v.errcode
}
