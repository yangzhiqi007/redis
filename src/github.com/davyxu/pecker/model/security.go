package model

import (
	"net/http"
)

const (
	publicKey      = "HlSGeorT9LlR1yMMUbsEFVgU2AbFIT1wDWlK35G17dSz62tKR4rWEjWW08t5vkZY"
	secrityKeyName = "PeckerKey"
)

func VerifyRequest(request *http.Request) bool {

	if request.Header.Get(secrityKeyName) == publicKey+*FlagPrivateKey {
		return true
	}

	log.Errorf("verify request failed, remote addr: %s, method: %s", request.RemoteAddr, request.Method)

	return false
}

func EncodeRequest(request *http.Request) {

	request.Header.Set(secrityKeyName, publicKey+*FlagPrivateKey)
}
