package model

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func GenToken() string {
	return strconv.Itoa(int(rand.Int31()))
}

// 会话
type Session struct {
	Token   string
	Account string
	Name    string
	Actor   string
	Time    time.Time
}

type SessionManager struct {
	SessionList     []*Session
	sessionListGuad sync.RWMutex `json:"-"`
}

func (self *SessionManager) GenSession(account string) (ses *Session) {
	ses = &Session{
		Token:   GenToken(),
		Account: account,
		Name:    account,
		Time:    time.Now(),
	}

	self.sessionListGuad.Lock()

	defer self.sessionListGuad.Unlock()

	for index, libSes := range self.SessionList {
		if libSes.Account == account {
			self.SessionList[index] = ses
			return
		}
	}

	self.SessionList = append(self.SessionList, ses)

	return
}

func (self *SessionManager) RemoveTimeoutSession() {

	self.sessionListGuad.Lock()

	defer self.sessionListGuad.Unlock()

	for len(self.SessionList) > 0 {

		var found bool
		for index, ses := range self.SessionList {
			if !ses.InTime() {
				// 每次删一个
				self.SessionList = append(self.SessionList[:index], self.SessionList[index+1:]...)
				found = true
				break
			}
		}

		// 没有发现超时的任务
		if !found {
			break
		}
	}

}

// 5天有效时间
func (self *Session) InTime() bool {
	return time.Now().Sub(self.Time) < time.Hour*24*5
}

func HandleAuthSession(ctx *gin.Context) {

	token := ctx.Request.Header.Get("Access-Token")

	DB.sessionListGuad.RLock()
	defer DB.sessionListGuad.RUnlock()

	for _, ses := range DB.SessionList {
		if ses.Token == token {
			ctx.Set("Session", ses)
			return
		}
	}

	ctx.String(http.StatusUnauthorized, "")
	ctx.Abort()
}
