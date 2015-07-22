package main

import (
	"time"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"

	"github.com/robskillington/gokit-gorilla-mux-starter/rpc"
)

func logging(logger log.Logger) func(rpc.CreateEntity) rpc.CreateEntity {
	return func(next rpc.CreateEntity) rpc.CreateEntity {
		return func(ctx context.Context, req *rpc.CreateEntityRequest) (rep *rpc.CreateEntityResponse, err error) {
			defer func(begin time.Time) {
				logger.Log("error", err, "took", time.Since(begin))
			}(time.Now())
			return next(ctx, req)
		}
	}
}

func instrument(requests metrics.Counter, duration metrics.TimeHistogram) func(rpc.CreateEntity) rpc.CreateEntity {
	return func(next rpc.CreateEntity) rpc.CreateEntity {
		return func(ctx context.Context, req *rpc.CreateEntityRequest) (rep *rpc.CreateEntityResponse, err error) {
			defer func(begin time.Time) {
				requests.Add(1)
				duration.Observe(time.Since(begin))
			}(time.Now())
			return next(ctx, req)
		}
	}
}
