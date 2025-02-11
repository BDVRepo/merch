// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameDocMerch = "doc_merchs"

// DocMerch mapped from table <doc_merchs>
type DocMerch struct {
	ID        *string    `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name      string     `gorm:"column:name;type:text;not null" json:"name"`
	Price     int32      `gorm:"column:price;type:integer;not null" json:"price"`
	CreatedAt *time.Time `gorm:"column:created_at;type:timestamp without time zone;index:idx_doc_merchs_created_at,priority:1;default:now()" json:"created_at"`
}

// TableName DocMerch's table name
func (*DocMerch) TableName() string {
	return TableNameDocMerch
}
