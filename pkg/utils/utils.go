package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/bschlaman/b-utils/pkg/logger"
)

// Adapter is a middleware adapter
type Adapter func(h http.Handler) http.Handler

// ReqData contains useful components of a request
type ReqData struct {
	Method      string `json:"method"`
	UrlPath     string `json:"urlPath"`
	RFC3339Time string `json:"time"`
	UnixTime    int64  `json:"unix"`
	Addr        string `json:"addr"`
	UserAgent   string `json:"user_agent"`
}

// ReqDataSimple contains a subset of request components
// useful for debugging.
type ReqDataSimple struct {
	Method  string `json:"method"`
	UrlPath string `json:"urlPath"`
}

// ParseRequest parses the http request and marshals it into json.
func ParseRequest(r *http.Request) ([]byte, error) {
	currTime := time.Now()
	jd, err := json.Marshal(&ReqData{
		r.Method,
		r.URL.Path,
		currTime.Format(time.RFC3339),
		currTime.Unix(),
		r.RemoteAddr,
		r.UserAgent(),
	})
	if err != nil {
		return nil, err
	}
	return jd, nil
}

// ParseRequestSimple parses the http request and marshals it into json.
// A smaller subset of the request attributes are selected.
func ParseRequestSimple(r *http.Request) ([]byte, error) {
	jd, err := json.Marshal(&ReqDataSimple{
		r.Method,
		r.URL.Path,
	})
	if err != nil {
		return nil, err
	}
	return jd, nil
}

// LogParseRequest parses the request and logs it
func LogParseRequest(l *logger.BLogger, r *http.Request) error {
	parsedReqBytes, err := ParseRequest(r)
	if err != nil {
		l.Error(err)
		return err
	}
	l.Info("received request:", string(parsedReqBytes))
	return nil
}

// LogParseRequestSimple parses the request and logs it
func LogParseRequestSimple(l *logger.BLogger, r *http.Request) error {
	parsedReqBytes, err := ParseRequestSimple(r)
	if err != nil {
		l.Error(err)
		return err
	}
	l.Info("received request:", string(parsedReqBytes))
	return nil
}

// LogReq returns an adapter that attempts to log and parse the request
// If an error is encountered, the error is logged by LogParseRequest
func LogReq(l *logger.BLogger) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			LogParseRequest(l, r)
			h.ServeHTTP(w, r)
		})
	}
}

// LogReqSimple returns an adapter that attempts to log and parse the request
// If an error is encountered, the error is logged by LogParseRequestSimple
func LogReqSimple(l *logger.BLogger) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			LogParseRequestSimple(l, r)
			h.ServeHTTP(w, r)
		})
	}
}

// EchoHandle returns an http.Handler that returns the
// output of ParseRequest in the http response. This is
// useful for debugging purposes
func EchoHandle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parsedReqBytes, _ := ParseRequest(r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(parsedReqBytes))
	})
}

// EchoDelayHandle does the same thing as EchoHandle
// but adds a delay in seconds that can be passed in as
// a query string parameter. This is useful for mocking network latency
func EchoDelayHandle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parsedReqBytes, _ := ParseRequest(r)

		maxDelaySeconds := 5
		delaySeconds, err := strconv.Atoi(r.URL.Query().Get("t"))
		if err != nil || delaySeconds > maxDelaySeconds {
			delaySeconds = maxDelaySeconds
		}

		// TODO: unmarshalling / remarshalling parsedReqBytes
		// is wasteful and error-prone; the unmarshal error case
		// also needs to be handled
		echoRequest := make(map[string]interface{})
		json.Unmarshal(parsedReqBytes, &echoRequest)

		echoRequest["delay"] = delaySeconds
		parsedReqBytes, _ = json.Marshal(&echoRequest)

		time.Sleep(time.Duration(delaySeconds) * time.Second)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(parsedReqBytes))
	})
}
