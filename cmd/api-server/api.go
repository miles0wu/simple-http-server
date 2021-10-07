package main

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

type AccessLog struct {
	URI  string
	IP   string
	Code int
}

func action(c *cli.Context) (err error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)

	server := &http.Server{
		Addr:    "0.0.0.0:" + c.String(_port),
		Handler: mux,
	}
	return server.ListenAndServe()
}

func index(w http.ResponseWriter, r *http.Request) {
	for key, values := range r.Header {
		w.Header().Add(key, strings.Join(values, ";"))
	}

	w.Header().Add("VERSION", os.Getenv("VERSION"))

	accessLog := AccessLog{
		URI:  r.RequestURI,
		IP:   "",
		Code: http.StatusOK,
	}

	reqIp, err := GetIP(r)
	if err != nil {
		accessLog.Code = http.StatusBadRequest
	}
	accessLog.IP = reqIp

	if accessLog.URI == "/healthz" {
		accessLog.Code = http.StatusOK
	}

	w.WriteHeader(accessLog.Code)

	j, err := json.Marshal(accessLog)
	if err != nil {
		log.Println("fail to encode log: ", accessLog)
		return
	}
	log.Println(string(j))
}

func GetIP(r *http.Request) (string, error) {
	ip := r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	ip = r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i, nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	return "", errors.New("no valid ip found")
}
