package services

import (
	"github.com/boardware-cloud/argus-service/argus"
	argusModel "github.com/boardware-cloud/model/argus"
	"github.com/boardware-cloud/model/common"
	"github.com/boardware-cloud/model/core"
)

func CreateMonitor(account core.Account, config argus.ArgusConfig) argus.Argus {
	entity := config.ToEntity(account)
	entity.AccountId = account.ID
	db.Save(&entity)
	a := argus.Argus{}
	a.SetEntity(entity)
	return a
}

func UpdateMonitor(monitorId any, config argus.ArgusConfig) {
	// a := config.ToEntity(account)
	// db.Save(&a)
}

func GetMonitor(id uint) argus.Argus {
	var model argusModel.Argus
	db.Find(&model, id)
	a := argus.Argus{}
	a.SetEntity(model)
	return a
}

func ListMonitors(accountId uint, index, limit int64) ([]argus.Argus, common.Pagination) {
	var list []argusModel.Argus
	pagination := common.ListEntity(&list, index, limit, db.Where("account_id = ?", accountId))
	var argusList []argus.Argus
	for _, item := range list {
		argusList = append(argusList, argus.NewArgus(item))
	}
	return argusList, pagination
}

func DeleteMonitor(a argus.Argus) {
	entity := a.Entity()
	db.Delete(&entity)
}
