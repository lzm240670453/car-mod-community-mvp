package model

type TradeIntent struct {
	ID       int64  `gorm:"column:id;primaryKey" json:"id"`
	PartID   int64  `gorm:"column:part_id" json:"partId"`
	BuyerID  int64  `gorm:"column:buyer_id" json:"buyerId"`
	SellerID int64  `gorm:"column:seller_id" json:"sellerId"`
	Message  string `gorm:"column:message" json:"message"`
	Status   int8   `gorm:"column:status" json:"status"`
	Timestamp
}

func (TradeIntent) TableName() string {
	return "trade_intents"
}
