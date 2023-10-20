package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewServer(host, port string) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: promhttp.Handler(),
	}
}
