package api

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"io/ioutil"
)

const boltURL = "https://bolt.kubemq.io/v1"

type license struct {
	Error       bool   `json:"error"`
	ErrorString string `json:"error_string"`
	Data        struct {
		License string `json:"license"`
	} `json:"data"`
}
type accountStatsResponse struct {
	Error       bool   `json:"error"`
	ErrorString string `json:"error_string"`
}

type accountStatsReport struct {
	Key                    string `json:"key"`
	Host                   string `json:"host"`
	Source                 int    `json:"source"`
	IsClustered            bool   `json:"is_clustered"`
	Version                string `json:"version"`
	ActivationReports      int    `json:"activation_reports"`
	ValidationByKeyLicense int    `json:"validation_by_key_license"`
}

func SaveToFile(filename string, data string) error {
	return ioutil.WriteFile(filename, []byte(data), 0644)
}

func GetLicenseData(key string, version string) (string, error) {
	req := resty.New().R()
	lic := &license{}
	r, err := req.SetResult(lic).SetError(lic).SetQueryParam("key", key).Get(boltURL + "/license")
	if err != nil {
		return "", err
	}

	if !r.IsSuccess() || r.IsError() {
		return "", err
	}
	if lic.Error {
		return "", errors.New(lic.ErrorString)
	}
	updateAccountStats(key, version)
	return lic.Data.License, nil
}

func updateAccountStats(key, version string) {
	req := resty.New().R()
	as := &accountStatsReport{
		Key:                    key,
		Host:                   "kubemqctl",
		Source:                 3,
		Version:                version,
		ActivationReports:      1,
		IsClustered:            true,
		ValidationByKeyLicense: 1,
	}
	asr := &accountStatsResponse{}
	req.SetResult(asr).SetError(asr).SetBody(as).Post(boltURL + "/update_stats")
}
