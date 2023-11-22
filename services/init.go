package services

import (
	"github.com/boardware-cloud/argus-service/argus"
	"github.com/boardware-cloud/common/config"
	"github.com/boardware-cloud/common/notifications"
	argusModel "github.com/boardware-cloud/model/argus"
	"github.com/boardware-cloud/model/core"
	"gorm.io/gorm"
)

var db *gorm.DB

var argusRepository argusModel.ArgusRepository

func Init(inject *gorm.DB) {
	db = inject
	core.Init(db)
	argusModel.Init(db)
	var emailSender notifications.Sender
	emailSender.SmtpHost = config.GetString("smtp.host")
	emailSender.Port = config.GetString("smtp.port")
	emailSender.Email = config.GetString("smtp.email")
	emailSender.Password = config.GetString("smtp.password")
	argusRepository = argusModel.NewArgusRepository(db)
	argus.Init(db, emailSender)
}
