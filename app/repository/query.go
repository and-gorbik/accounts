package repository

import (
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"

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

func buildAccountUpdateQuery(a domain.AccountUpdate, cityID, countryID uuid.UUID) (string, []interface{}, error) {
	setMap := make(map[string]interface{})
	if cityID != uuid.Nil {
		setMap[shortName(AccountCityID)] = cityID
	}
	if countryID != uuid.Nil {
		setMap[shortName(AccountCountryID)] = countryID
	}

	if a.Email != nil {
		setMap[shortName(AccountEmail)] = a.Email
	}
	if a.Birth != nil {
		setMap[shortName(AccountBirth)] = util.TimestampToDatetime((*int64)(a.Birth))
	}
	if a.Status != nil {
		setMap[shortName(AccountStatus)] = util.TimestampToDatetime((*int64)(a.Birth))
	}

	return squirrel.Update(TableAccount).
		SetMap(setMap).
		Where(squirrel.Eq{AccountID: a.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
}

func buildAccountInsertQuery(a domain.AccountModel) (string, []interface{}, error) {
	return squirrel.Insert(TableAccount).
		Columns(shortName(AccountID), shortName(AccountStatus), shortName(AccountEmail),
			shortName(AccountSex), shortName(AccountBirth), shortName(AccountFirstname),
			shortName(AccountSurname), shortName(AccountPhone), shortName(AccountCountryID),
			shortName(AccountCityID), shortName(AccountJoined), shortName(AccountPremStart), shortName(AccountPremEnd)).
		Values(a.ID, a.Status, a.Email, a.Sex, a.Birth, a.Name, a.Surname,
			a.Phone, a.CountryID, a.CityID, a.Joined, a.PremiumStart, a.PremiumEnd).
		Suffix(returning(AccountID)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
}

func buildCityInsertQuery(c domain.CityModel) (sql string, values []interface{}, err error) {
	insertCity := squirrel.Insert(TableCity).
		Columns(shortName(CityID), shortName(CityName)).
		Values(c.ID, c.Name).
		Suffix(onConflictDoNothing(CityName)).
		Suffix(returning(shortName(CityID)))

	sql, values, err = cte("inserted", insertCity)
	if err != nil {
		return
	}

	return union(
		squirrel.Select("inserted.id").From("inserted").Prefix(sql, values...),
		squirrel.Select(CityID).From(TableCity).Where(squirrel.Eq{CityName: c.Name}),
	)
}

func buildCountryInsertQuery(c domain.CountryModel) (sql string, values []interface{}, err error) {
	insertCountry := squirrel.Insert(TableCountry).
		Columns(shortName(CountryID), shortName(CountryName)).
		Values(c.ID, c.Name).
		Suffix(onConflictDoNothing(CountryName)).
		Suffix(returning(shortName(CountryID)))

	sql, values, err = cte("inserted", insertCountry)
	if err != nil {
		return
	}

	return union(
		squirrel.Select("inserted.id").From("inserted").Prefix(sql, values...),
		squirrel.Select(CountryID).From(TableCountry).Where(squirrel.Eq{CountryName: c.Name}),
	)
}

func join(table, left, right string) string {
	return fmt.Sprintf("%s ON %s = %s", table, left, right)
}

func returning(columns ...string) string {
	return "RETURNING " + strings.Join(columns, ",")
}

func onConflictDoNothing(column string) string {
	return fmt.Sprintf("ON CONFLICT(%s) DO NOTHING", shortName(column))
}

func shortName(column string) string {
	lst := strings.Split(column, ".")
	if len(lst) != 2 {
		return column
	}

	return lst[1]
}

func union(l, r squirrel.SelectBuilder) (string, []interface{}, error) {
	lSQL, values, err := l.ToSql()
	if err != nil {
		return "", nil, err
	}

	rSQL, rValues, err := r.ToSql()
	if err != nil {
		return "", nil, err
	}

	values = append(values, rValues...)
	sql, err := squirrel.Dollar.ReplacePlaceholders(fmt.Sprintf(`%s UNION %s`, lSQL, rSQL))
	if err != nil {
		return "", nil, err
	}

	return sql, values, nil
}

func cte(name string, q squirrel.Sqlizer) (sql string, values []interface{}, err error) {
	sql, values, err = q.ToSql()
	if err != nil {
		return
	}

	sql = fmt.Sprintf("WITH %s AS (%s)", name, sql)
	return
}
