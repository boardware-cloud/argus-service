package services

import (
	"github.com/boardware-cloud/argus-service/argus"
	errorCode "github.com/boardware-cloud/common/code"
	argusModel "github.com/boardware-cloud/model/argus"
	"github.com/boardware-cloud/model/common"
	"github.com/boardware-cloud/model/core"
)

var argusService *ArgusService

func GetArgusService() *ArgusService {
	if argusService == nil {
		argusService = NewArgusService()
	}
	return argusService
}

func NewArgusService() *ArgusService {
	return &ArgusService{
		argusRepository: argusModel.GetArgusRepository(),
	}
}

type ArgusService struct {
	argusRepository *argusModel.ArgusRepository
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

func (as ArgusService) UpdateMonitor(id uint, config argus.ArgusConfig) (argus.Argus, error) {
	a, err := as.GetMonitor(id)
	if err != nil {
		return a, err
	}
	entity := a.Entity()
	entity.Update(config.ToEntity())
	a.SetEntity(entity)
	argus.Spawn(a)
	return a, nil
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

func (as ArgusService) GetMonitor(id uint) (argus.Argus, error) {
	model := argusRepository.GetById(id)
	a := argus.Argus{}
	if model == nil {
		return a, errorCode.ErrNotFound
	}
	a.SetEntity(*model)
	return a, nil
}
