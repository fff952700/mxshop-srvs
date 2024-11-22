package model

import "time"

type BaseModel struct {
	Id        uint32    `gorm:"primary_key;AUTO_INCREMENT"`
	CreatedAt time.Time `gorm:"type:datetime;DEFAULT:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:datetime;DEFAULT:CURRENT_TIMESTAMP"`
	IsDel     bool      `gorm:"type:tinyint(1);DEFAULT:1;comment:'false 1 (not deleted), true 2 (deleted)'"`
}

type User struct {
	BaseModel
	Mobile   string     `gorm:"type:varchar(11);unique;NOT NULL;"`
	Nickname string     `gorm:"type:varchar(255);"`
	Password string     `gorm:"type:varchar(255);NOT NULL;"`
	Birthday *time.Time `gorm:"type:datetime;DEFAULT:CURRENT_TIMESTAMP"`
	Gender   uint8      `gorm:"type:tinyint(1);DEFAULT:1;comment:'male 1, female 2'"`
	Role     uint8      `gorm:"type:tinyint(1);NOT NULL;DEFAULT:2;comment:'admin 1, user 2'"`
}
