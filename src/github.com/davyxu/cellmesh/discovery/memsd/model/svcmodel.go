package model

import (
	"github.com/davyxu/cellmesh/util"
	"github.com/davyxu/cellnet"
	"strings"
)

const (
	ServiceKeyPrefix = "_svcdesc_"
)

var (
	Queue cellnet.EventQueue
	IDGen = meshutil.NewUUID64Generator()

	Listener cellnet.Peer
	Debug    bool

	Version = "1.3.1"

	sesByToken = map[string]cellnet.Session{}
)

func IsServiceKey(rawkey string) bool {

	return strings.HasPrefix(rawkey, ServiceKeyPrefix)
}

func GetSvcIDByServiceKey(rawkey string) string {

	if IsServiceKey(rawkey) {
		return rawkey[len(ServiceKeyPrefix):]
	}

	return ""
}

func init() {
	IDGen.AddTimeComponent(8)
	IDGen.AddSeqComponent(8, 0)
}

func GetSessionToken(ses cellnet.Session) (token string) {
	ses.(cellnet.ContextSet).FetchContext("token", &token)

	return
}

func SessionCount() int {
	return len(sesByToken)
}

func AddSession(ses cellnet.Session, token string) {
	sesByToken[token] = ses
}

func RemoveSession(ses cellnet.Session) {
	token := GetSessionToken(ses)
	delete(sesByToken, token)
}

func Broadcast(msg interface{}) {
	for _, ses := range sesByToken {
		ses.Send(msg)
	}
}

func TokenExists(token string) (ret bool) {

	if _, ok := sesByToken[token]; ok {
		return true
	}

	return false
}
