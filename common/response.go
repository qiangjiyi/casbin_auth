package common

import "github.com/qiangjiyi/casbin_auth/common/errors"

// Response common request response struct
type Response struct {
	ErrCode    errors.ErrorCodeType `json:"err_code"`
	ErrMessage string               `json:"err_message"`
	Payload    interface{}          `json:"payload,omitempty"`
}
