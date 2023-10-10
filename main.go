package main

import (
	"context"

	"github.com/boardware-cloud/argus-service/controllers"
	"github.com/boardware-cloud/model"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.ReadInConfig()
	port := ":" + viper.GetString("server.port")
	user := viper.GetString("database.user")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	dbport := viper.GetString("database.port")
	database := viper.GetString("database.database")
	DB, err := model.NewConnection(user, password, host, dbport, database)
	if err != nil {
		panic(err)
	}
	context := context.WithValue(context.Background(), "db", DB)
	controllers.Init(context)
	controllers.Run(port)
}
