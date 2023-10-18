package argus

import (
	"github.com/boardware-cloud/common/notifications"
	argusModel "github.com/boardware-cloud/model/argus"
	"gorm.io/gorm"
)

var db *gorm.DB

var node *Node

var argusRepository argusModel.ArgusRepository
var emailSender notifications.Sender

func Init(inject *gorm.DB, _emailSender notifications.Sender) {
	db = inject
	emailSender = _emailSender
	argusRepository = argusModel.NewArgusRepository(db)
	Register()
}

func Spawn(a Argus) {
	go a.Spawn(*node)
}
