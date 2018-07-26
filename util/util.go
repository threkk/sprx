// Package util Contains a few useful functions to be shared across the
// different packages.
package util

import (
	"io"
	"net/http"
)

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers#hbh
var hbhHeaders = [...]string{"Connection", "Keep-Alive", "Proxy-Authenticate",
	"Proxy-Authorization", "TE", "Trailer", "Transfer-Encoding", "Upgrade"}

// isHbHHeader - Checks if a given header is a Hop-by-Hop header and needs to be
// skipped.
func isHbHHeader(header string) bool {
	for _, h := range hbhHeaders {
		if h == header {
			return true
		}
	}
	return false
}

// CopyHeader - Copies Headers from one request to another one.
func CopyHeader(dst, src http.Header) {
	for k, vs := range src {
		// If the header is a Hop-By-Hop header, don't copy it.
		if !isHbHHeader(k) {
			for _, v := range vs {
				dst.Add(k, v)
			}
		}
	}
}

// Transfer - Transfers the content of a reader to a writer.
func Transfer(dst io.WriteCloser, src io.ReadCloser) {
	defer dst.Close()
	defer src.Close()
	io.Copy(dst, src)
}
