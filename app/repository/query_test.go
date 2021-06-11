package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"accounts/domain"
	"accounts/util"
)

func Test_buildAccountSearchQuery_Success(t *testing.T) {
	f := NewFilter()
	f.Eq(AccountSex, "m")
	f.Domain(AccountEmail, "test.ru")
	f.Any(AccountFirstname, []interface{}{"Андрей", "Иван"})
	f.Null(CountryName, false)
	f.Limit = 10

	sql, values, err := buildAccountSearchQuery(f)
	if err != nil {
		t.Fatal(err)
	}

	expected := "SELECT account.id, account.email, account.sex, account.name, country.name FROM account "
	expected += "JOIN country ON country.id = account.country_id "
	expected += "WHERE account.sex = $1 AND account.email LIKE $2 "
	expected += "AND account.name IN ($3,$4) AND country.name IS NOT NULL "
	expected += "ORDER BY account.id DESC "
	expected += "LIMIT 10"

	assert.Equal(t, expected, sql)
	assert.Equal(t, 4, len(values))
}

func Test_buildAccountUpdateQuery_Success(t *testing.T) {
	email := domain.FieldEmail("test@test.ru")
	acc := domain.AccountUpdate{
		ID:    1,
		Email: &email,
	}

	sql, values, err := buildAccountUpdateQuery(acc, 2, 3)
	if err != nil {
		t.Fatal(err)
	}

	expected := "UPDATE account SET account.city_id = $1, "
	expected += "account.country_id = $2, account.email = $3 "
	expected += "WHERE account.id = $4"

	assert.Equal(t, expected, sql)
	assert.Equal(t, 4, len(values))
}

func Test_buildAccountInsertQuery_Success(t *testing.T) {
	now := time.Now()
	account := domain.AccountModel{
		ID:           1,
		Status:       "заняты",
		Email:        "test@test.ru",
		Sex:          "m",
		Birth:        now,
		Joined:       now,
		Name:         util.PtrString("Иван"),
		Surname:      util.PtrString("Иванов"),
		Phone:        util.PtrString("8(999)7654321"),
		CountryID:    util.PtrInt32(1),
		CityID:       util.PtrInt32(2),
		PremiumStart: &now,
		PremiumEnd:   &now,
	}

	sql, values, err := buildAccountInsertQuery(account)
	if err != nil {
		t.Fatal(err)
	}

	expected := "INSERT INTO account (account.id,account.status,account.email,"
	expected += "account.sex,account.birth,account.name,account.surname,account.phone,"
	expected += "account.country_id,account.city_id,account.joined,account.prem_start,account.prem_end) "
	expected += "VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING account.id"

	assert.Equal(t, expected, sql)
	assert.Equal(t, 13, len(values))
}

func Test_buildCityInsertQuery_Success(t *testing.T) {
	city := domain.CityModel{
		Name: "Москва",
	}

	sql, values, err := buildCityInsertQuery(city)
	if err != nil {
		t.Fatal(err)
	}

	expected := "INSERT INTO city (city.name) VALUES ($1) "
	expected += "ON CONFLICT(city.name) DO NOTHING RETURNING city.id"

	assert.Equal(t, expected, sql)
	assert.Equal(t, 1, len(values))
}

func Test_buildCountryInsertQuery_Success(t *testing.T) {
	country := domain.CountryModel{
		Name: "Россия",
	}

	sql, values, err := buildCountryInsertQuery(country)
	if err != nil {
		t.Fatal(err)
	}

	expected := "INSERT INTO country (country.name) VALUES ($1) "
	expected += "ON CONFLICT(country.name) DO NOTHING RETURNING country.id"

	assert.Equal(t, expected, sql)
	assert.Equal(t, 1, len(values))
}
