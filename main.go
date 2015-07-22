package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	stdlog "log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/statsd"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"github.com/robskillington/gokit-gorilla-mux-starter/deps"
	"github.com/robskillington/gokit-gorilla-mux-starter/rpc"
	"github.com/robskillington/gokit-gorilla-mux-starter/services"
)

type callback func() (interface{}, error)

func main() {
	// Flag domain
	fs := flag.NewFlagSet("", flag.ExitOnError)

	httpAddr := fs.String("http.addr", ":8000", "Address for HTTP (JSON) server")
	debugAddr := fs.String("debug.addr", ":8001", "Address for HTTP debug/instrumentation server")

	flag.Usage = fs.Usage // only show our flags
	fs.Parse(os.Args[1:])

	// `package log` domain
	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(os.Stderr)
	logger = kitlog.NewContext(logger).With("ts", kitlog.DefaultTimestampUTC)
	stdlog.SetOutput(kitlog.NewStdlibAdapter(logger)) // redirect stdlib logging to us
	stdlog.SetFlags(0)                                // flags are handled in our logger

	// `package metrics` domain
	requests := metrics.NewMultiCounter(
		statsd.NewCounter(ioutil.Discard, "requests_total", time.Second),
	)
	duration := metrics.NewTimeHistogram(time.Nanosecond, metrics.NewMultiHistogram(
		statsd.NewHistogram(ioutil.Discard, "duration_nanoseconds_total", time.Second),
	))

	// Dependencies and services
	dependencies := &deps.All{
		EntityService: &services.EntityService{},
	}

	// RPCs
	var createEntity rpc.CreateEntity = rpc.NewCreateEntity(dependencies)
	createEntity = logging(logger)(createEntity)
	createEntity = instrument(requests, duration)(createEntity)

	// Mechanical stuff
	rand.Seed(time.Now().UnixNano())
	root := context.Background()
	errc := make(chan error)

	go func() {
		errc <- interrupt()
	}()

	// Transport: HTTP (JSON)
	go func() {
		ctx, cancel := context.WithCancel(root)
		defer cancel()

		router := mux.NewRouter()

		router.HandleFunc("/entity", func(w http.ResponseWriter, r *http.Request) {
			// vars := mux.Vars(r)

			var request rpc.CreateEntityRequest
			httpJsonBodyEndpoint(w, r, request, func() (interface{}, error) {
				return createEntity(ctx, &request)
			})
		}).Methods("POST")

		logger.Log("addr", *httpAddr, "transport", "HTTP/JSON")
		errc <- http.ListenAndServe(*httpAddr, router)
	}()

	// Transport: HTTP (debug/instrumentation)
	go func() {
		logger.Log("addr", *debugAddr, "transport", "debug")
		errc <- http.ListenAndServe(*debugAddr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/health" {
				io.WriteString(w, "OK ;-)\n")
			} else {
				w.WriteHeader(http.StatusNotFound)
				io.WriteString(w, "404 not found\n")
			}
		}))
	}()

	logger.Log("fatal", <-errc)
}

func interrupt() error {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	return fmt.Errorf("%s", <-c)
}
