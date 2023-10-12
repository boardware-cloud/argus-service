package argus

import (
	"gorm.io/gorm"
)

var db *gorm.DB

var node *Node

func Init(inject *gorm.DB) {
	db = inject
	Register()
}

func Spawn(a Argus) {
	go a.Spawn(*node)
}
