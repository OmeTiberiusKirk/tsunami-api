package tdss

import (
	"gorm.io/gorm"
)

type TdssIntf interface {
	Controllers() ControllersIntf
	Services() ServicesIntf
}

type Tdss struct {
	DB *gorm.DB
}

func NewTdss(db *gorm.DB) (tdss TdssIntf) {
	tdss = &Tdss{
		DB: db,
	}
	return
}
