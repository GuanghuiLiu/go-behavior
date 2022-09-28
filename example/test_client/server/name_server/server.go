package main

import (
	nameServer "github.com/GuanghuiLiu/behavior/name_server"
)

func main() {
	nameServer.Start("common", 9987)
}
