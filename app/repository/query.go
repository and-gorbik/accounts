package repository

import (
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"

	"accounts/domain"
	"accounts/util"
)

const (
	TableAccount  = "account"
	TableLike     = "likes"
	TableInterest = "interest"
	TableCity     = "city"
	TableCountry  = "country"

	AccountID         = "account.id"
	AccountStatus     = "account.status"
	AccountEmail      = "account.email"
	AccountSex        = "account.sex"
	AccountBirth      = "account.birth"
	AccountJoined     = "account.joined"
	AccountFirstname  = "account.name"
	AccountSurname    = "account.surname"
	AccountPhone      = "account.phone"
	AccountCountryID  = "account.country_id"
	AccountCityID     = "account.city_id"
	AccountPremStart  = "account.prem_start"
	AccountPremEnd    = "account.prem_end"
	LikesLikerID      = "likes.liker_id"
	LikesLikeeID      = "likes.likee_id"
	LikesTimestamp    = "likes.ts"
	InterestAccountID = "interest.account_id"
	InterestName      = "interest.name"
	CityID            = "city.id"
	CityName          = "city.name"
	CountryID         = "country.id"
	CountryName       = "country.name"
)

func buildAccountSearchQuery(f *Filter) (string, []interface{}, error) {
	where, params, err := f.Build()
	if err != nil {
		return "", nil, err
	}

	q := squirrel.Select(AccountID, AccountEmail).
		PlaceholderFormat(squirrel.Dollar).
		From(TableAccount).
		Where(where, params...).
		OrderBy(AccountID + " DESC")

	for column := range f.Columns() {
		switch column {
		case AccountSex, AccountStatus, AccountBirth, AccountPhone, AccountFirstname, AccountSurname:
			q = q.Column(column)
		case CityName:
			q = q.Column(column).Join(join(TableCity, CityID, AccountCityID))
		case CountryName:
			q = q.Column(column).Join(join(TableCountry, CountryID, AccountCountryID))
		case LikesLikerID:
			q = q.Join(join(TableLike, LikesLikerID, AccountID))
		case InterestName:
			q = q.Join(join(TableInterest, InterestAccountID, AccountID))
		}
	}

	if f.Limit != 0 {
		q = q.Limit(uint64(f.Limit))
	}

	return q.ToSql()
}

func buildAccountUpdateQuery(a domain.AccountUpdate, cityID, countryID int32) (string, []interface{}, error) {
	setMap := make(map[string]interface{})
	setMap[AccountCityID] = cityID
	setMap[AccountCountryID] = countryID

	if a.Email != nil {
		setMap[AccountEmail] = a.Email
	}
	if a.Birth != nil {
		setMap[AccountBirth] = util.TimestampToDatetime((*int64)(a.Birth))
	}
	if a.Status != nil {
		setMap[AccountStatus] = util.TimestampToDatetime((*int64)(a.Birth))
	}

	return squirrel.Update(TableAccount).
		SetMap(setMap).
		Where(squirrel.Eq{AccountID: a.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
}

func buildAccountInsertQuery(a domain.AccountModel) (string, []interface{}, error) {
	return squirrel.Insert(TableAccount).
		Columns(AccountID, AccountStatus, AccountEmail,
			AccountSex, AccountBirth, AccountFirstname,
			AccountSurname, AccountPhone, AccountCountryID,
			AccountCityID, AccountJoined, AccountPremStart, AccountPremEnd).
		Values(a.ID, a.Status, a.Email, a.Sex, a.Birth, a.Name, a.Surname,
			a.Phone, a.CountryID, a.CityID, a.Joined, a.PremiumStart, a.PremiumEnd).
		Suffix(returning(AccountID)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
}

func buildCityInsertQuery(c domain.CityModel) (string, []interface{}, error) {
	return squirrel.Insert(TableCity).
		Columns(CityName).
		Values(c.Name).
		Suffix(onConflictDoNothing(CityName)).
		Suffix(returning(CityID)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
}

func buildCountryInsertQuery(c domain.CountryModel) (string, []interface{}, error) {
	return squirrel.Insert(TableCountry).
		Columns(CountryName).
		Values(c.Name).
		Suffix(onConflictDoNothing(CountryName)).
		Suffix(returning(CountryID)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
}

func join(table, left, right string) string {
	return fmt.Sprintf("%s ON %s = %s", table, left, right)
}

func returning(columns ...string) string {
	return "RETURNING " + strings.Join(columns, ",")
}

func onConflictDoNothing(column string) string {
	return fmt.Sprintf("ON CONFLICT(%s) DO NOTHING", column)
}
