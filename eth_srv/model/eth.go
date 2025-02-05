package model

type User struct {
	BaseModel
	Nickname string `gorm:"type:varchar(20);not null;default:''"`
	Phone    string `gorm:"type:varchar(20);unique;not null;default:'';index:idx_phone_phone;"`
}
