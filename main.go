package main

import (
	"github.com/boardware-cloud/argus-service/controllers"
	"github.com/boardware-cloud/common/config"
	"github.com/boardware-cloud/model"
)

func main() {
	port := ":" + config.GetString("server.port")
	user := config.GetString("database.user")
	password := config.GetString("database.password")
	host := config.GetString("database.host")
	dbport := config.GetString("database.port")
	database := config.GetString("database.database")
	DB, err := model.NewConnection(user, password, host, dbport, database)
	if err != nil {
		panic(err)
	}
	controllers.Init(DB)
	controllers.Run(port)
}
