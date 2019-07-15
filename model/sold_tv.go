package model

import "time"

type SoldTv struct {
	ID       int64 `sql:"id"`
	TvID     int64 `sql:"tv_id" db:"tv_id"`
	Returned bool
	SaleDate time.Time `sql:"sale_date" db:"sale_date"`
	Quantity int64
}
