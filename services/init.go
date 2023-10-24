package services

import (
	"github.com/boardware-cloud/argus-service/argus"
	"github.com/boardware-cloud/common/notifications"
	argusModel "github.com/boardware-cloud/model/argus"
	"github.com/boardware-cloud/model/core"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var db *gorm.DB

var argusRepository argusModel.ArgusRepository

func Init(inject *gorm.DB) {
	db = inject
	core.Init(db)
	argusModel.Init(db)
	var emailSender notifications.Sender
	emailSender.SmtpHost = viper.GetString("smtp.host")
	emailSender.Port = viper.GetString("smtp.port")
	emailSender.Email = viper.GetString("smtp.email")
	emailSender.Password = viper.GetString("smtp.password")
	argusRepository = argusModel.NewArgusRepository(db)
	argus.Init(db, emailSender)
}
