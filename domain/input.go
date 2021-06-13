package domain

import (
	"errors"

	"accounts/util"

	"github.com/google/uuid"
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

func (a *AccountInput) AccountModel(cityID, countryID uuid.UUID) *AccountModel {
	if !a.validated {
		return nil
	}

	table := &AccountModel{
		ID:        int32(*a.ID),
		Status:    string(*a.Status),
		Email:     string(*a.Email),
		Sex:       string(*a.Sex),
		Birth:     *util.TimestampToDatetime((*int64)(a.Birth)),
		Name:      (*string)(a.Name),
		Surname:   (*string)(a.Surname),
		Phone:     (*string)(a.Phone),
		CountryID: util.PtrUUID(countryID),
		CityID:    util.PtrUUID(cityID),
		Joined:    *util.TimestampToDatetime((*int64)(a.Joined)),
	}

	if a.Premium != nil {
		table.PremiumStart = util.TimestampToDatetime((*int64)(a.Premium.Start))
		table.PremiumEnd = util.TimestampToDatetime((*int64)(a.Premium.End))
	}

	return table
}

func (a *AccountInput) LikeModels() []LikeModel {
	if a.Likes == nil || len(a.Likes) == 0 || !a.validated {
		return nil
	}

	models := make([]LikeModel, 0, len(a.Likes))
	for _, like := range a.Likes {
		models = append(models, LikeModel{
			LikerID:   int32(*a.ID),
			LikeeID:   int32(*like.UserID),
			Timestamp: *util.TimestampToDatetime((*int64)(like.Timestamp)),
		})
	}

	return models
}

func (a *AccountInput) InterestModels() []InterestModel {
	if a.Interests == nil || len(a.Interests) == 0 || !a.validated {
		return nil
	}

	models := make([]InterestModel, 0, len(a.Interests))
	for _, interest := range a.Interests {
		models = append(models, InterestModel{
			AccountID: int32(*a.ID),
			Name:      string(*interest),
		})
	}

	return models
}

func (a *AccountInput) CityModel() *CityModel {
	if a.City == nil || !a.validated {
		return nil
	}

	return &CityModel{
		ID:   uuid.New(),
		Name: string(*a.City),
	}
}

func (a *AccountInput) CountryModel() *CountryModel {
	if a.Country == nil || !a.validated {
		return nil
	}

	return &CountryModel{
		ID:   uuid.New(),
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
	ID FieldID `json:"-"`

	Email   *FieldEmail   `json:"email,omitempty"`
	Birth   *FieldBirth   `json:"birth,omitempty"`
	City    *FieldCity    `json:"city,omitempty"`
	Country *FieldCountry `json:"country,omitempty"`
	Status  *FieldStatus  `json:"status,omitempty"`

	validated bool
}

func (a *AccountUpdate) Validate() error {
	defer func() {
		a.validated = true
	}()

	return checkValidators(&a.ID, a.Email, a.Birth, a.City, a.Country, a.Status)
}

func (a *AccountUpdate) AccountModel(cityID, countryID *uuid.UUID) *AccountModel {
	if !a.validated {
		return nil
	}

	return &AccountModel{
		ID:        int32(a.ID),
		Status:    string(*a.Status),
		Email:     string(*a.Email),
		Birth:     *util.TimestampToDatetime((*int64)(a.Birth)),
		CityID:    cityID,
		CountryID: countryID,
	}
}

func (a *AccountUpdate) CityModel() *CityModel {
	if a.City == nil || !a.validated {
		return nil
	}

	return &CityModel{
		Name: string(*a.City),
	}
}

func (a *AccountUpdate) CountryModel() *CountryModel {
	if a.Country == nil || !a.validated {
		return nil
	}

	return &CountryModel{
		Name: string(*a.Country),
	}
}

type LikesInput struct {
	Likes     []LikeInput `json:"likes"`
	validated bool
}

func (li *LikesInput) Validate() error {
	defer func() {
		li.validated = true
	}()

	for _, like := range li.Likes {
		if err := like.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (li *LikesInput) LikeModels() []LikeModel {
	if !li.validated {
		return nil
	}

	likeModels := make([]LikeModel, 0, len(li.Likes))
	for _, like := range li.Likes {
		likeModels = append(likeModels, LikeModel{
			LikerID:   int32(*like.Liker),
			LikeeID:   int32(*like.Likee),
			Timestamp: *util.TimestampToDatetime((*int64)(like.Timestamp)),
		})
	}

	return likeModels
}

type LikeInput struct {
	Likee     *FieldID        `json:"likee"`
	Liker     *FieldID        `json:"liker"`
	Timestamp *FieldTimestamp `json:"ts"`
}

func (li *LikeInput) Validate() error {
	if util.AnyIsNil(li.Likee, li.Liker, li.Timestamp) {
		return errEmptyField
	}

	return checkValidators(li.Likee, li.Liker, li.Timestamp)
}
