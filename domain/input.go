package domain

import (
	"errors"

	"accounts/util"
)

var (
	errEmptyField = errors.New("empty field")
)

type AccountInput struct {
	// required fields
	ID     *FieldID     `json:"id"`
	Email  *FieldEmail  `json:"email"`
	Sex    *FieldSex    `json:"sex"`
	Birth  *FieldBirth  `json:"birth"`
	Joined *FieldJoined `json:"joined"`
	Status *FieldStatus `json:"status"`

	// not required fields
	Name      *FieldFirstname     `json:"fname"`
	Surname   *FieldSurname       `json:"sname"`
	Phone     *FieldPhone         `json:"phone"`
	Country   *FieldCountry       `json:"country"`
	City      *FieldCity          `json:"city"`
	Interests []*FieldInterest    `json:"interests"`
	Premium   *PremiumInput       `json:"premium"`
	Likes     []*AccountLikeInput `json:"likes"`

	validated bool
}

func (a *AccountInput) Validate() error {
	defer func() {
		a.validated = true
	}()

	if util.AnyIsNil(a.ID, a.Email, a.Sex, a.Birth, a.Joined, a.Status) {
		return errEmptyField
	}

	for _, interest := range a.Interests {
		if interest == nil {
			return errEmptyField
		}
		if interest.Validate() != nil {
			return errInvalidValue
		}
	}

	for _, like := range a.Likes {
		if like == nil {
			return errEmptyField
		}
		if like.Validate() != nil {
			return errInvalidValue
		}
	}

	return checkValidators(
		a.ID, a.Email, a.Sex, a.Birth, a.Name,
		a.Surname, a.Phone, a.Country, a.City,
		a.Joined, a.Status, a.Premium,
	)
}

func (a *AccountInput) GetPerson() PersonTable {
	if !a.validated {
		return PersonTable{}
	}

	return PersonTable{
		ID:      int32(*a.ID),
		Email:   string(*a.Email),
		Sex:     string(*a.Sex),
		Birth:   int64(*a.Birth),
		Name:    (*string)(a.Name),
		Surname: (*string)(a.Surname),
		Phone:   (*string)(a.Phone),
	}
}

func (a *AccountInput) GetAccount() AccountTable {
	if !a.validated {
		return AccountTable{}
	}

	table := AccountTable{
		ID:     int32(*a.ID),
		Joined: int64(*a.Joined),
		Status: string(*a.Status),
	}

	if a.Premium != nil {
		table.PremiumStart = (*int64)(a.Premium.Start)
		table.PremiumEnd = (*int64)(a.Premium.End)
	}

	return table
}

func (a *AccountInput) GetLikes() []LikeTable {
	if !a.validated || a.Likes == nil || len(a.Likes) == 0 {
		return nil
	}

	tables := make([]LikeTable, 0, len(a.Likes))
	for _, like := range a.Likes {
		tables = append(tables, LikeTable{
			LikerID:   int32(*a.ID),
			LikeeID:   int32(*like.UserID),
			Timestamp: int64(*like.Timestamp),
		})
	}

	return tables
}

func (a *AccountInput) GetInterests() []InterestTable {
	if !a.validated || a.Interests == nil || len(a.Interests) == 0 {
		return nil
	}

	tables := make([]InterestTable, 0, len(a.Interests))
	for _, interest := range a.Interests {
		tables = append(tables, InterestTable{
			AccountID: int32(*a.ID),
			Name:      string(*interest),
		})
	}

	return tables
}

func (a *AccountInput) GetCity() *CityTable {
	if !a.validated || a.City == nil {
		return nil
	}

	return &CityTable{
		Name: string(*a.City),
	}
}

func (a *AccountInput) GetCountry() *CountryTable {
	if !a.validated || a.Country == nil {
		return nil
	}

	return &CountryTable{
		Name: string(*a.Country),
	}
}

type PremiumInput struct {
	Start *FieldPremium `json:"start"`
	End   *FieldPremium `json:"finish"`
}

func (p *PremiumInput) Validate() error {
	if p == nil {
		return nil
	}

	if util.AnyIsNil(p.Start, p.End) {
		return errEmptyField
	}

	return checkValidators(p.Start, p.End)
}

type AccountLikeInput struct {
	UserID    *FieldID        `json:"id"`
	Timestamp *FieldTimestamp `json:"ts"`
}

func (a *AccountLikeInput) Validate() error {
	if a == nil {
		return nil
	}

	if util.AnyIsNil(a.UserID, a.Timestamp) {
		return errEmptyField
	}

	return checkValidators(a.UserID, a.Timestamp)
}

type AccountUpdate struct {
	ID      FieldID     `json:"-"`
	Email   FieldEmail  `json:"email,omitempty"`
	Birth   FieldBirth  `json:"birth,omitempty"`
	City    *FieldID    `json:"city,omitempty"`
	Country *FieldID    `json:"country,omitempty"`
	Status  FieldStatus `json:"status,omitempty"`
}

type LikeInput struct {
	Likee     FieldID        `json:"likee"`
	Liker     FieldID        `json:"liker"`
	Timestamp FieldTimestamp `json:"ts"`
}
