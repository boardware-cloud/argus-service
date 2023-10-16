package argus

import (
	argusModel "github.com/boardware-cloud/model/argus"
	"gorm.io/gorm"
)

var db *gorm.DB

var node *Node

var argusRepository argusModel.ArgusRepository

func Init(inject *gorm.DB) {
	db = inject
	Register()
	argusRepository = argusModel.NewArgusRepository(db)
}

func Spawn(a Argus) {
	go a.Spawn(*node)
}
