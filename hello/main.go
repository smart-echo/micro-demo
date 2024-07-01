package main

import (
	"context"
	"sync"

	"github.com/smart-echo/micro-demo/hello/handler"
	pb "github.com/smart-echo/micro-demo/hello/proto"

	ot "github.com/smart-echo/micro-plugins/wrapper/trace/opentelemetry"
	"github.com/smart-echo/micro"
	"github.com/smart-echo/micro/logger"
	"github.com/smart-echo/micro/server"

)

var (
	service = "hello"
	version = "latest"
)

func main() {
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	// Create tracer
	tp, err := initTracerProvider(service, version)
	if err != nil {
		logger.Fatal("failed to initial tracer provider: %v", err)
	}
	traceOpts := ot.WithTraceProvider(tp)
	defer func(ctx context.Context) {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Infof("error shutting down the tracer provider: %v", err)
		}
	}(ctx)

	// Create service
	srv := micro.NewService(
		micro.BeforeStart(func() error {
			logger.Infof("Starting service %s", service)
			return nil
		}),
		micro.BeforeStop(func() error {
			logger.Infof("Shutting down service %s", service)
			cancel()
			return nil
		}),
		micro.AfterStop(func() error {
			wg.Wait()
			return nil
		}),
		micro.WrapCall(ot.NewCallWrapper(traceOpts)),
		micro.WrapClient(ot.NewClientWrapper(traceOpts)),
		micro.WrapHandler(ot.NewHandlerWrapper(traceOpts)),
		micro.WrapSubscriber(ot.NewSubscriberWrapper(traceOpts)),
	)
	srv.Init(
		micro.Name(service),
		micro.Version(version),
	)

	srv.Server().Init(
		server.Wait(&wg),
	)

	ctx = server.NewContext(ctx, srv.Server())

	// Register handler
	if err := pb.RegisterHelloHandler(srv.Server(), new(handler.Hello)); err != nil {
		logger.Fatal(err)
	}
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
