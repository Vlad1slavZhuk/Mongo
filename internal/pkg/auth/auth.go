package auth

import (
	"Mongo/internal/pkg/constErr"
	"Mongo/internal/pkg/data"
	"log"
	"net/http"
	"strings"
)

func IsAccountExists(acc *data.Account, baseAcc []*data.Account) (*data.Account, error) {
	for _, account := range baseAcc {
		if account.Username == acc.Username && account.Password == acc.Password {
			return account, nil
		}
	}
	return nil, constErr.NotFoundAcc
}

func GetTokenFromHeader(r *http.Request) (string, error) {
	text := r.Header.Get("Authorization")
	arr := strings.Fields(text)
	if len(arr) != 2 || arr[0] != "Bearer" {
		return "", constErr.NoValidToken
	}
	return arr[1], nil
}

func AccountIdentification(r *http.Request, baseAcc []*data.Account) (*data.Account, error) {
	token, err := GetTokenFromHeader(r)
	if err != nil {
		log.Println(1)
		return nil, err // emptyToken or incorrect token
	}

	if err = VerifyToken(token); err != nil {
		log.Println(2)
		return nil, err // no valid token
	}

	acc, err := ContainsToken(token, baseAcc)
	if err != nil {
		log.Println(3)
		return nil, err // no contains token
	}
	return acc, nil
}
