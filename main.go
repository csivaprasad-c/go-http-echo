package main

import (
	"context"
	"flag"
	"fmt"
	"go-http-echo/version"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	listenFlag  = flag.String("listen", ":5678", "address and port to listen")
	textFlag    = flag.String("text", "", "text to put on the webpage")
	versionFlag = flag.Bool("version", false, "display version information")

	stdoutW = os.Stdout
	stderrW = os.Stderr
)

func main() {
	flag.Parse()

	if *versionFlag {
		_, _ = fmt.Fprintln(stderrW, version.HumanVersion)
		os.Exit(0)
	}

	if *textFlag == "" {
		_, _ = fmt.Fprintln(stderrW, "Missing -text option!")
		os.Exit(127)
	}

	args := flag.Args()
	if len(args) > 0 {
		_, _ = fmt.Fprintln(stderrW, "Too many arguments!")
		os.Exit(127)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", httpLog(stdoutW, withAppHeaders(httpEcho(*textFlag))))
	mux.HandleFunc("/health", withAppHeaders(httpHealth()))

	server := &http.Server{
		Addr:    *listenFlag,
		Handler: mux,
	}
	serverCh := make(chan struct{})
	go func() {
		log.Printf("[INFO] server is listening on %s\n", *listenFlag)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("[ERR] server exited with: %s", err)
		}
		close(serverCh)
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	<-signalCh

	log.Printf("[INFO] received interrupt, shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("[ERR] failed to shutdown server: %s", err)
	}

	// If we got this far, it was an interrupt, so don't exit cleanly
	os.Exit(2)
}

func httpEcho(v string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, v)
	}
}

func httpHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, `{"status":"ok"}`)
	}
}
