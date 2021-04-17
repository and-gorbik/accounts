package models

type Person struct {
	ID        int32   `db:"id"`
	Email     string  `db:"email"`
	Sex       string  `db:"sex"`
	Birth     int64   `db:"birth"`
	Name      *string `db:"name"`
	Surname   *string `db:"surname"`
	Phone     *string `db:"phone"`
	CountryID *int32  `db:"country_id"`
	CityID    *int32  `db:"city_id"`
}

type Account struct {
	ID           int32  `db:"id"`
	Joined       int64  `db:"joined"`
	Status       string `db:"status"`
	PremiumStart *int64 `db:"prem_start"`
	PremiumEnd   *int64 `db:"prem_end"`
}

type Like struct {
	LikerID   int32 `db:"liker_id"`
	LikeeID   int32 `db:"likee_id"`
	Timestamp int64 `db:"ts"`
}

type Interest struct {
	AccountID int32  `db:"account_id"`
	Name      string `db:"name"`
}

type City struct {
	ID   int32  `db:"id"`
	Name string `db:"name"`
}

type Country struct {
	ID   int32  `db:"id"`
	Name string `db:"name"`
}
