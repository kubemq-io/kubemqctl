package kubemq

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

//const kubemqAPI = "http://localhost:3000/v1"
const kubemqAPI = "https://bolt.kubemq.io/v1"

type RegistrationRequest struct {
	Name            string `json:"name"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	GeneratedKey    string `json:"generated_key"`
	RegistrationKey string `json:"registration_key"`
}

type RegistrationResponse struct {
	Error       bool   `json:"error"`
	ErrorString string `json:"error_string"`
	Data        struct {
		Name     string `json:"name"`
		Username string `json:"username"`
	} `json:"data"`
}

func RegisterKubeMQ(request *RegistrationRequest) error {
	req := resty.New().R()
	res := &RegistrationResponse{}
	r, err := req.SetResult(res).SetError(res).SetBody(request).Post(kubemqAPI + "/register-ext-source")
	if err != nil {
		return err
	}
	if r != nil && res.Error {
		return fmt.Errorf(res.ErrorString)
	}
	return nil
}

func ValidateRegisterKubeMQ(key, reg string) error {
	req := resty.New().R()
	res := &RegistrationResponse{}
	r, err := req.SetResult(res).SetError(res).SetQueryParam("key", key).SetQueryParam("reg", reg).Get(kubemqAPI + "/validate-ext-source")
	if err != nil {
		return err
	}
	if r != nil && res.Error {
		return fmt.Errorf(res.ErrorString)
	}
	return nil
}
