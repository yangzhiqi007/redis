package model

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"
)

var (
	AdminPassword string
)

func InitAdminPassword() {
	rand.Seed(time.Now().UnixNano())

	code := GenToken()

	log.Infof("admin: %s", code)

	AdminPassword = fmt.Sprintf("%x", md5.Sum([]byte(code)))
}
