package main

import (
	"fmt"
	"go-http-echo/version"
	"io"
	"net/http"
	"time"
)

const (
	httpHeaderAppName    string = "X-App-Name"
	httpHeaderAppVersion string = "X-App-Version"

	httpLogDateFormat string = "2006/01/02 15:04:05"
	httpLogFormat     string = "%v %s %s \"%s %s %s\" %d %d \"%s\" %v\n"
)

func withAppHeaders(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(httpHeaderAppName, version.Name)
		w.Header().Set(httpHeaderAppVersion, version.Version)
		h(w, r)
	}
}

type metaResponseWriter struct {
	writer http.ResponseWriter
	status int
	length int
}

func (w *metaResponseWriter) Header() http.Header {
	return w.writer.Header()
}

func (w *metaResponseWriter) WriteHeader(s int) {
	w.status = s
	w.writer.WriteHeader(s)
}

func (w *metaResponseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	w.length = len(b)
	return w.writer.Write(b)
}

func httpLog(out io.Writer, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var mrw metaResponseWriter
		mrw.writer = w

		defer func(start time.Time) {
			status := mrw.status
			length := mrw.length
			end := time.Now()
			dur := end.Sub(start)
			_, err := fmt.Fprintf(out, httpLogFormat, end.Format(httpLogDateFormat),
				r.Host, r.RemoteAddr, r.Method, r.URL.Path, r.Proto,
				status, length, r.UserAgent(), dur)
			if err != nil {
				return
			}
		}(time.Now())
		handlerFunc(&mrw, r)
	}
}
