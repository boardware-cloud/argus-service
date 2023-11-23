package services

import (
	"github.com/boardware-cloud/model"
	argusModel "github.com/boardware-cloud/model/argus"
)

var db = model.GetDB()

var argusRepository = argusModel.GetArgusRepository()
