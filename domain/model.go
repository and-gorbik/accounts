package domain

import (
	"time"

	"github.com/google/uuid"
)

type AccountModel struct {
	ID     int32     `db:"id"`
	Status string    `db:"status"`
	Email  string    `db:"email"`
	Sex    string    `db:"sex"`
	Birth  time.Time `db:"birth"`
	Joined time.Time `db:"joined"`

	Name         *string    `db:"name"`
	Surname      *string    `db:"surname"`
	Phone        *string    `db:"phone"`
	CountryID    *uuid.UUID `db:"country_id"`
	CityID       *uuid.UUID `db:"city_id"`
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
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

type CountryModel struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}
