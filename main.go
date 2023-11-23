package main

import (
	"github.com/boardware-cloud/argus-service/controllers"
	"github.com/boardware-cloud/common/config"
)

func main() {
	port := ":" + config.GetString("server.port")
	controllers.Run(port)
}
