package cmd

import (
	"io/ioutil"
	"log"
	"os"
)

var defaultConfig = `statsAddress: "http://localhost:8080/v1/stats" #the address of Stats endpoint, you can replace the localhost:8080 with your address
healthAddress: "http://localhost:8080/health" # the address of Health endpoint , you can replace the localhost:8080 with your address
metricsAddress: "http://localhost:8080/metrics" #the address of Health endpoint, you can replace the localhost:8080 with your address
monitorAddress: "ws://localhost:8080/v1/stats" #the address of Monitor endpoint, you can replace the localhost:8080 with your address
connections:
  - kind: 1 # 1 - grpc 2- rest
    host: "localhost" # host destination
    port: 50000 # port destination
    isSecured: false # set using https
    certFile: "" # set location of cert file
  - kind: 2 # 1 - grpc 2- rest
    host: "localhost" # host destination
    port: 9090  # port destination
    isSecured: false  # set using https
    certFile: "" # set location of cert file - not in use for Rest
`

func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
func checkConfigFile() {
	if !exists(".config.yaml") {
		err := ioutil.WriteFile(".config.yaml", []byte(defaultConfig), 0644)
		if err != nil {
			log.Fatalf("error creating config file, %s", err.Error())
		}
	}
}
