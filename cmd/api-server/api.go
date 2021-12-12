package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"simple-http-server/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
)

type AccessLog struct {
	URI  string
	IP   string
	Code int
}

func action(c *cli.Context) (err error) {
	metrics.Register()

	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/healthz", healthz)
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    "0.0.0.0:" + c.String(_port),
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// add graceful shutdown
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return server.Shutdown(ctx)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok\n")
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func index(w http.ResponseWriter, r *http.Request) {
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()
	delay := randInt(0, 2000)
	time.Sleep(time.Millisecond * time.Duration(delay))

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

	w.WriteHeader(accessLog.Code)

	j, err := json.Marshal(accessLog)
	if err != nil {
		log.Println("fail to encode log: ", accessLog)
		return
	}
	log.Println(string(j))
	log.Printf("Respond in %d ms", delay)
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
