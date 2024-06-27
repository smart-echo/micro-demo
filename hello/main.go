package main

import (

	"github.com/smart-echo/micro-demo/hello/handler"
	pb "github.com/smart-echo/micro-demo/hello/proto"

	"github.com/smart-echo/micro"
	"github.com/smart-echo/micro/logger"

)

var (
	service = "hello"
	version = "latest"
)

func main() {
	// Create service
	srv := micro.NewService(
	)
	srv.Init(
		micro.Name(service),
		micro.Version(version),
	)

	// Register handler
	if err := pb.RegisterHelloHandler(srv.Server(), new(handler.Hello)); err != nil {
		logger.Fatal(err)
	}
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
