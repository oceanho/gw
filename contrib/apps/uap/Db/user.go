package Db

import (
	"fmt"
	"github.com/oceanho/gw/backend/gwdb"
	"time"
)

var tablePrefix = "gw_uap"

func getTableName(name string) string {
	return fmt.Sprintf("%s_%s", tablePrefix, name)
}

type User struct {
	gwdb.Model
	TenantId  uint64 `gorm:"default:0;not null;UNIQUEINDEX:idx_tenant_id_passport"`
	Passport  string `gorm:"type:varchar(32);UNIQUEINDEX:idx_tenant_id_passport;not null"`
	Secret    string `gorm:"type:varchar(128);not null"`
	IsUser    bool   `gorm:"default:0;not null"`
	IsAdmin   bool   `gorm:"default:0;not null"`
	IsTenancy bool   `gorm:"default:0;not null"`
	gwdb.HasLockState
	gwdb.HasCreationState
	gwdb.HasModificationState
	gwdb.HasSoftDeletionState
	gwdb.HasActivationState
}

func (User) TableName() string {
	return getTableName("user")
}

type UserProfile struct {
	gwdb.Model
	gwdb.HasTenantState
	Gender   uint8  `gorm:"default:4"` // 1.man, 2.woman, 3.custom, 4.unknown
	UserID   uint64 `gorm:"index"`
	Name     string `gorm:"type:varchar(64);index"`
	Email    string `gorm:"type:varchar(128);index"`
	Phone    string `gorm:"type:varchar(16);index"`
	Avatar   string `gorm:"type:varchar(256)"`
	Address  string `gorm:"type:varchar(256)"`
	PostCode string `gorm:"type:varchar(16)"`
	BirthDay *time.Time
	gwdb.HasCreationState
	gwdb.HasModificationState
}

func (UserProfile) TableName() string {
	return getTableName("user_profile")
}

type UserAccessKeySecret struct {
	gwdb.Model
	gwdb.HasTenantState
	Key    string
	Secret string
	gwdb.HasCreationState
	gwdb.HasActivationState
	gwdb.HasModificationState
}

func (UserAccessKeySecret) TableName() string {
	return getTableName("user_aks")
}
