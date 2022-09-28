package main

import (
	"github.com/GuanghuiLiu/behavior/example/test_client/server/server"
	"github.com/GuanghuiLiu/behavior/gateway"
	"github.com/GuanghuiLiu/behavior/model"
)

func main() {

	model.InitCluster("common", "bj_0001_00091", true, 9979)
	f := server.NewRouter
	gateway.Start(true, 9989, f)
}
