package api

import (
	"errors"
	"github.com/go-resty/resty/v2"
)

const boltURL = "https://bolt.kubemq.io"

type response struct {
	Error       bool   `json:"error"`
	ErrorString string `json:"error_string"`
	Data        string `json:"data"`
}

func GetLicenseDataByToken(key string) (string, error) {
	req := resty.New().R()
	lic := &response{}
	r, err := req.SetResult(lic).SetError(lic).SetQueryParam("key", key).Get(boltURL + "/getlicensebykey")
	if err != nil {
		return "", err
	}

	if !r.IsSuccess() || r.IsError() {
		return "", err
	}
	if lic.Error {
		return "", errors.New(lic.ErrorString)
	}

	return lic.Data, nil
}
