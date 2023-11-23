package argus

import (
	"github.com/boardware-cloud/common/notifications"
	"github.com/boardware-cloud/model"

	argusModel "github.com/boardware-cloud/model/argus"
)

var db = model.GetDB()

var node *Node

var argusRepository = argusModel.GetArgusRepository()
var emailSender = notifications.GetEmailSender()

func init() {
	Register()
}

func Spawn(a Argus) {
	go a.Spawn(*node)
}
