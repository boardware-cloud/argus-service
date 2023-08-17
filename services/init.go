package services

import (
	"github.com/boardware-cloud/common/config"
	"github.com/boardware-cloud/common/notifications"
	"github.com/boardware-cloud/model"
	uptime "github.com/boardware-cloud/model/argus"

	"gorm.io/gorm"
)

var DB *gorm.DB

var node UptimeNode

var emailSender notifications.Sender

func init() {
	user := config.GetString("database.user")
	password := config.GetString("database.password")
	host := config.GetString("database.host")
	port := config.GetString("database.port")
	database := config.GetString("database.database")
	var err error
	DB, err = model.NewConnection(user, password, host, port, database)
	if err != nil {
		panic(err)
	}
	// Init email sender
	emailSender.SmtpHost = config.GetString("smtp.host")
	emailSender.Port = config.GetString("smtp.port")
	emailSender.Email = config.GetString("smtp.email")
	emailSender.Password = config.GetString("smtp.password")
	DB.AutoMigrate(&uptime.Monitor{})
	DB.AutoMigrate(&uptime.UptimeNode{})
	DB.AutoMigrate(&uptime.MonitoringRecord{})
	DB.AutoMigrate(&uptime.UptimeMonitorAlert{})
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
