package util

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/qiangjiyi/casbin_auth/common/errors"
	"github.com/satori/go.uuid"
	"strconv"
	"time"
)

// Md5Encrypt encrypt plaintext using "md5" hash algorithm
func Md5Encrypt(plaintext string) string {
	hash := md5.New()
	hash.Write([]byte(plaintext))
	return hex.EncodeToString(hash.Sum(nil))
}

// Sha256Encrypt
func Sha256Encrypt(plaintext string) string {
	hash := sha256.New()
	hash.Write([]byte(plaintext))
	return hex.EncodeToString(hash.Sum(nil))
}

// ParseRequestParams handle and parse request params
func ParseRequestParams(requestBody []byte, params interface{}) errors.AuthError {
	if err := json.Unmarshal(requestBody, params); err != nil {
		logs.GetLogger().Printf("Request param unmarshal error ---> %+v", err)
		return errors.ErrDeserializeDataFail
	}
	return nil
}

// GenerateUserLoginToken generate user login token by userId, current timestamp and uuid
func GenerateUserLoginToken(userId int) (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		logs.Error("generate UUID failed ---> %+v", )
		return "", err
	}
	return Sha256Encrypt(strconv.Itoa(userId) + strconv.Itoa(time.Now().Nanosecond()) + u.String()), nil
}
