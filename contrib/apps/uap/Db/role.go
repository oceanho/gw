package Db

import "github.com/oceanho/gw/backend/gwdb"

type Role struct {
	gwdb.Model
	gwdb.HasTenantState
	Name       string `gorm:"type:varchar(32);not null"`
	Descriptor string `gorm:"type:varchar(128);not null"`
	gwdb.HasCreationState
	gwdb.HasModificationState
}

func (Role) TableName() string {
	return getTableName("role")
}
