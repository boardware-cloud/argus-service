package services

import (
	"github.com/boardware-cloud/argus-service/argus"
	"github.com/boardware-cloud/common/code"
	argusModel "github.com/boardware-cloud/model/argus"
	"github.com/boardware-cloud/model/common"
	"github.com/boardware-cloud/model/core"
)

func CreateMonitor(account core.Account, config argus.ArgusConfig) argus.Argus {
	entity := config.ToEntity()
	entity.AccountId = account.ID()
	db.Save(&entity)
	a := new(argus.Argus)
	a.SetEntity(entity)
	argus.Spawn(*a)
	return *a
}

func UpdateMonitor(id uint, config argus.ArgusConfig) (argus.Argus, error) {
	entity := argusRepository.GetById(id)
	if entity == nil {
		return argus.Argus{}, nil
	}
	entity.Update(config.ToEntity())
	a := new(argus.Argus)
	a.SetEntity(*entity)
	argus.Spawn(*a)
	return *a, nil
}

func GetMonitor(id uint) (argus.Argus, error) {
	model := argusRepository.GetById(id)
	a := argus.Argus{}
	if model == nil {
		return a, code.ErrNotFound
	}
	a.SetEntity(*model)
	return a, nil
}

func ListMonitors(accountId uint, index, limit int64) ([]argus.Argus, common.Pagination) {
	var list []argusModel.Argus
	pagination := common.ListEntity(&list, index, limit, "", db.Where("account_id = ?", accountId))
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

func ListRecords(argusId uint, index, limit int64) ([]argus.Record, common.Pagination) {
	var list []argusModel.ArgusRecord
	pagination := common.ListEntity(&list, index, limit, "created_at DESC", db.Where("argus_id = ?", argusId))
	var records []argus.Record
	for _, i := range list {
		records = append(records, argus.Record{
			Result:       argus.ResultStatus(i.Result),
			ResponesTime: i.ResponesTime,
			CheckedAt:    i.CreatedAt,
		})
	}
	return records, pagination
}
