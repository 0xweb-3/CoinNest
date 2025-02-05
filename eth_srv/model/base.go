package model

import (
	"time"

	"github.com/sony/sonyflake"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

// Sonyflake generator
var sf *sonyflake.Sonyflake

func init() {
	st := sonyflake.Settings{
		StartTime: time.Now(),
	}
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("Sonyflake not created")
	}
}

// BeforeCreate will set a unique ID using Sonyflake
func (base *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := sf.NextID()
	if err != nil {
		return err
	}
	base.ID = id
	return nil
}
