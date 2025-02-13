// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameDocTransaction = "doc_transactions"

// DocTransaction mapped from table <doc_transactions>
type DocTransaction struct {
	ID         *string    `gorm:"column:id;type:text;primaryKey" json:"id"`
	SenderID   string     `gorm:"column:sender_id;type:text;not null;index:idx_doc_transactions_sender_receiver,priority:1" json:"sender_id"`
	ReceiverID *string    `gorm:"column:receiver_id;type:text;index:idx_doc_transactions_sender_receiver,priority:2" json:"receiver_id"`
	Amount     int32      `gorm:"column:amount;type:integer;not null" json:"amount"`
	CreatedAt  *time.Time `gorm:"column:created_at;type:timestamp without time zone;index:idx_doc_transactions_created_at,priority:1;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName DocTransaction's table name
func (*DocTransaction) TableName() string {
	return TableNameDocTransaction
}
