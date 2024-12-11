package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type BaseModel struct {
	Id        int32     `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt time.Time `gorm:"type:datetime;DEFAULT:CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt time.Time `gorm:"type:datetime;DEFAULT:CURRENT_TIMESTAMP" json:"-"`
	IsDel     bool      `gorm:"type:tinyint(1);DEFAULT:0;comment:'false 0 (not deleted), true 1 (deleted)'" json:"-"`
}

type GormList []string

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), g)
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}
