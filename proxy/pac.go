package proxy

import (
	"fmt"
	"net/http"
)

// Simple PAC configuration which will route all the traffic through the proxy.
// Reference: https://en.wikipedia.org/wiki/Proxy_auto-config
const pacTpl = `
function FindProxyForURL(url, host) {
	return "PROXY 127.0.0.1:%s; DIRECT";
}
`

// NewPacHandler Creates a new PAC handler which redirect all the traffic to the
// given port.
func NewPacHandler(port string) func(http.ResponseWriter, http.Request) {
	pac := fmt.Sprintf(pacTpl, port)
	return func(res http.ResponseWriter, req http.Request) {
		fmt.Fprintf(res, pac)
	}
}
