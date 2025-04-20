package models

type PageViews struct {
	Id     string `gorm:"column:id;primaryKey"`
	UserId string `gorm:"column:userId"`
}
