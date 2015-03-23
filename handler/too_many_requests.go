package handler

import (
	"fmt"
	"net"
	"net/http"
)

func TooManyRequests(w http.ResponseWriter, r *http.Request) {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	http.Error(w,
		fmt.Sprintf(
			"Too many requests from your IP %s; try again in %s seconds",
			host,
			w.Header().Get("Retry-After"),
		),
		429,
	)
}
