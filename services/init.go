package services

import (
	"github.com/boardware-cloud/common/constants"
	"github.com/boardware-cloud/common/notifications"
	"github.com/boardware-cloud/model"
	argus "github.com/boardware-cloud/model/argus"
	"github.com/boardware-cloud/model/billing"
	"github.com/spf13/viper"

	"gorm.io/gorm"
)

var DB *gorm.DB

var node UptimeNode

var emailSender notifications.Sender

var serviceTitle = "Uptime monitor"
var serviceDescription = "Uptime monitoring system"
var serviceUrl = "/uptime"

func init() {
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	user := viper.GetString("database.user")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")
	database := viper.GetString("database.database")
	DB, err = model.NewConnection(user, password, host, port, database)
	emailSender.SmtpHost = viper.GetString("smtp.host")
	emailSender.Port = viper.GetString("smtp.port")
	emailSender.Email = viper.GetString("smtp.email")
	emailSender.Password = viper.GetString("smtp.password")
	if err != nil {
		panic(err)
	}
	argus.Init(DB)
	// Billing configuartion
	billing.Init(DB)
	billing.AutoMigrate(constants.ARGUS,
		serviceTitle,
		serviceDescription,
		serviceUrl,
		constants.SAAS,
	)
	node = NewUptimeNode()
	node.Register()
	go KeepAlive()
	go KeepCheckNodes()
	go KeepCheckMontiors()
}

type List[T any] struct {
	Data       []T
	Pagination Pagination
}

type Pagination struct {
	Index int64
	Limit int64
	Total int64
}
