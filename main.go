package main

import (
	"github.com/micro/go-micro/client"
	"github.com/micro/go-web"
	proto "github.com/micro/auth-srv/proto/account"
	"github.com/micro/auth-web/handler"
)

func main() {
	service := web.NewService(
		web.Name("go.micro.web.auth"),
		web.Handler(handler.Router()),
	)

	service.Init()

	handler.Init(
		"templates",
		proto.NewAccountClient("go.micro.srv.auth", client.DefaultClient),
	)

	service.Run()
}
