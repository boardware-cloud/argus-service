package services

import (
	"context"

	"github.com/boardware-cloud/argus-service/argus"
	"github.com/boardware-cloud/common/notifications"
	argusModel "github.com/boardware-cloud/model/argus"
	"github.com/boardware-cloud/model/core"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var db *gorm.DB

var emailSender notifications.Sender

var accountRepository core.AccountRepository
var argusRepository argusModel.ArgusRepository

func Init(inject context.Context) {
	db = inject.Value("db").(*gorm.DB)
	core.Init(db)
	argusModel.Init(db)
	argus.Init(db)
	emailSender.SmtpHost = viper.GetString("smtp.host")
	emailSender.Port = viper.GetString("smtp.port")
	emailSender.Email = viper.GetString("smtp.email")
	emailSender.Password = viper.GetString("smtp.password")
	accountRepository = core.NewAccountRepository(db)
	argusRepository = argusModel.NewArgusRepository(db)
}
