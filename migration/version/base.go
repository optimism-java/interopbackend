package migration

import "time"

type Base struct {
	ID        int64     `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	CreatedAt time.Time `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}
