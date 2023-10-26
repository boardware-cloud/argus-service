package services

import (
	"github.com/boardware-cloud/argus-service/argus"
	argusModel "github.com/boardware-cloud/model/argus"
	"github.com/boardware-cloud/model/common"
	"github.com/boardware-cloud/model/core"
)

func CreateMonitor(account core.Account, config argus.ArgusConfig) (argus.Argus, error) {
	entity := config.ToEntity()
	entity.AccountId = account.ID()
	a := new(argus.Argus)
	db.Save(&entity)
	argus.Spawn(*a)
	a.SetEntity(entity)
	return *a, nil
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
