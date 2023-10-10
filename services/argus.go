package services

import (
	"github.com/boardware-cloud/argus-service/argus"
	"github.com/boardware-cloud/model/core"
)

func CreateMonitor(account core.Account, config argus.ArgusConfig) {
	a := config.ToEntity(account)
	db.Save(&a)
}

func UpdateMonitor(monitorId any, config argus.ArgusConfig) {
	// a := config.ToEntity(account)
	// db.Save(&a)
}

func GetMonitor(id uint) {

}
