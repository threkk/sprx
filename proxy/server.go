// Package proxy Simple proxy server heavily based on
// https://medium.com/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c
// with several improvements.
package proxy

import (
	"crypto/tls"
	"github.com/threkk/sprx/util"
	"io"
	"net"
	"net/http"
	"time"
)

// handleHTTP - Handles HTTP requests by resolving the original request and
// copying the return to the response.
func handleHTTP(res http.ResponseWriter, req *http.Request) {
	// Resolve the request using the default transport.
	trip, err := http.DefaultTransport.RoundTrip(req)

	// If error, return 503.
	if err != nil {
		http.Error(res, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer trip.Body.Close()

	// Copy the headers of the trip to the response.
	util.CopyHeader(trip.Header, res.Header())
	res.WriteHeader(trip.StatusCode)
	// Copy the body of the trip to the response.
	io.Copy(res, trip.Body)
}

// handleHTTPS - Handles the HTTP Connect method used by the proxy to handle
// HTTPS traffic and other protocols.
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods/CONNECT
// Esentially, it hijacks the connection and sues raw TCP.
func handleHTTPS(res http.ResponseWriter, req *http.Request) {
	// Open a TCP connection to the destination.
	dst, err := net.DialTimeout("tcp", req.Host, 10*time.Second)
	if err != nil {
		http.Error(res, err.Error(), http.StatusServiceUnavailable)
		return
	}

	// Connection successful, we'll return a 200.
	res.WriteHeader(http.StatusOK)

	// Hijack the connection to use raw TCP packages.
	hijacker, ok := res.(http.Hijacker)
	if !ok {
		http.Error(res, "Connection hijacking not supported, likely beacuse using HTTP/2 server", http.StatusServiceUnavailable)
		return
	}

	// We ignore the buffer as we are going to copy everything from the
	// destination to the client without any intermetiate step.
	client, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(res, err.Error(), http.StatusServiceUnavailable)
		return
	}

	// Transfer everything between client and destination.
	go util.Transfer(dst, client)
	go util.Transfer(client, dst)
}

// NewProxy - Creates a new proxy server.
func NewProxy() *http.Server {
	server := &http.Server{
		// Handle HTTPS and HTTP traffic.
		Handler: http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if req.Method == http.MethodConnect {
				handleHTTPS(res, req)
			} else {
				handleHTTP(res, req)
			}
		}),
		// Disable HTTP/2, required for hijacking.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	return server
}
