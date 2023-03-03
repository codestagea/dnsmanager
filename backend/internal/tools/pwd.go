package tools

import (
	"encoding/base64"
	"errors"
	"github.com/codestagea/bindmgr/global"
	"golang.org/x/crypto/bcrypt"
	"unicode"
)

func HashPwdCompare(hashPwd string, pwd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(pwd))
	if err != nil {
		return false, err
	}
	return true, nil
}

func VerifyPwdStrength(s string) error {
	strength := make(map[string]bool)
	if len(s) < 8 {
		return errors.New("password length should more than 8 characters")
	}
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			strength["number"] = true
		case unicode.IsUpper(c):
			strength["upper"] = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			strength["special"] = true
		case unicode.IsLower(c):
			strength["lower"] = true
		default:
		}
	}
	if len(strength) < 3 {
		return errors.New("password should contains at least 3 of uppercase, lowercase, number and special character")
	}
	return nil
}

func PwdHash(pwd string) (string, error) {
	if pwd == "" {
		return "", errors.New("password cannot be empty")
	}

	if hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost); err != nil {
		return "", err
	} else {
		return string(hash), nil
	}
}

func PwdRsaDecode(pwd string) (string, error) {
	text, err := base64.StdEncoding.DecodeString(pwd)
	if err != nil {
		return "", err
	}
	if plainPwd, err := RsaDecrypt(global.GVA_CONF.Rsa.PrivateKey, text); err != nil {
		return "", err
	} else {
		return string(plainPwd), nil
	}

}
