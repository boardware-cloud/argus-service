package services

import (
	"github.com/boardware-cloud/argus-service/argus"
	argusModel "github.com/boardware-cloud/model/argus"
	"github.com/boardware-cloud/model/common"
	"github.com/boardware-cloud/model/core"
	"gorm.io/gorm"
)

func NewArgusService(db *gorm.DB) ArgusService {
	return ArgusService{
		argusRepository: argusModel.NewArgusRepository(db),
	}
}

type ArgusService struct {
	argusRepository argusModel.ArgusRepository
}

func (as ArgusService) CreateMonitor(account core.Account, config argus.ArgusConfig) (argus.Argus, error) {
	entity := config.ToEntity()
	entity.AccountId = account.ID()
	a := new(argus.Argus)
	as.argusRepository.Save(&entity)
	argus.Spawn(*a)
	a.SetEntity(entity)
	return *a, nil
}

func (as ArgusService) ListMonitors(accountId uint, index, limit int64) ([]argus.Argus, common.Pagination) {
	var list []argusModel.Argus
	pagination := common.ListEntity(&list, index, limit, "", db.Where("account_id = ?", accountId))
	var argusList []argus.Argus
	for _, item := range list {
		argusList = append(argusList, argus.NewArgus(item))
	}
	return argusList, pagination
}
