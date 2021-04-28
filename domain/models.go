package domain

type Person struct {
	ID        FieldID         `db:"account_id"`
	Email     FieldEmail      `db:"email"`
	Sex       FieldSex        `db:"sex"`
	Birth     FieldBirth      `db:"birth"`
	Name      *FieldFirstname `db:"name"`
	Surname   *FieldSurname   `db:"surname"`
	Phone     *FieldPhone     `db:"phone"`
	CountryID *FieldID        `db:"country_id"`
	CityID    *FieldID        `db:"city_id"`
}

type Account struct {
	ID           FieldID      `db:"id"`
	Joined       FieldJoined  `db:"joined"`
	Status       FieldStatus  `db:"status"`
	PremiumStart FieldPremium `db:"prem_start"`
	PremiumEnd   FieldPremium `db:"prem_end"`
}

type Like struct {
	LikerID   FieldID `db:"liker_id"`
	LikeeID   FieldID `db:"likee_id"`
	Timestamp int64   `db:"ts"`
}

type Interest struct {
	AccountID FieldID       `db:"account_id"`
	Name      FieldInterest `db:"name"`
}

type City struct {
	ID   FieldID   `db:"id"`
	Name FieldCity `db:"name"`
}

type Country struct {
	ID   FieldID      `db:"id"`
	Name FieldCountry `db:"name"`
}
