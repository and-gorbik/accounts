package domain

import (
	"time"
)

type PersonModel struct {
	ID      int32     `db:"account_id"`
	Status  string    `db:"status"`
	Email   string    `db:"email"`
	Sex     string    `db:"sex"`
	Birth   time.Time `db:"birth"`
	Name    *string   `db:"name"`
	Surname *string   `db:"surname"`
	Phone   *string   `db:"phone"`

	CountryID *int32 `db:"country_id"`
	CityID    *int32 `db:"city_id"`
}

type AccountModel struct {
	ID           int32      `db:"id"`
	Joined       time.Time  `db:"joined"`
	PremiumStart *time.Time `db:"prem_start"`
	PremiumEnd   *time.Time `db:"prem_end"`
}

type LikeModel struct {
	LikerID   int32     `db:"liker_id"`
	LikeeID   int32     `db:"likee_id"`
	Timestamp time.Time `db:"ts"`
}

type InterestModel struct {
	AccountID int32  `db:"account_id"`
	Name      string `db:"name"`
}

type CityModel struct {
	ID   int32  `db:"id"`
	Name string `db:"name"`
}

type CountryModel struct {
	ID   int32  `db:"id"`
	Name string `db:"name"`
}
