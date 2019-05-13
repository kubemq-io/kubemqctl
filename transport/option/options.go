package option

import "fmt"

type ConnectionType int

const (
	ConnectionTypeUndefeind ConnectionType = iota
	ConnectionTypeGrpc
	ConnectionTypeRest
)

var ConnectionTypeNames = map[ConnectionType]string{
	ConnectionTypeGrpc: "gRPC",
	ConnectionTypeRest: "Rest",
}

type Options struct {
	Kind      ConnectionType
	Host      string
	Port      int
	ApiPort   int
	IsSecured bool
	CertFile  string
}

func NewOptions(connType ConnectionType, host string, port int, isSecured bool, certFile string) *Options {
	return &Options{
		Kind:      connType,
		Host:      host,
		Port:      port,
		IsSecured: isSecured,
		CertFile:  certFile,
	}
}

func (o *Options) Uri() string {
	scheme := "http"
	if o.IsSecured {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s:%d", scheme, o.Host, o.Port)
}

func (o *Options) WebsocketUri() string {
	scheme := "ws"
	if o.IsSecured {
		scheme = "wss"
	}
	return fmt.Sprintf("%s://%s:%d", scheme, o.Host, o.Port)
}
