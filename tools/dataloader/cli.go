package dataloader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli/v2"

	"accounts/models"
)

const (
	flagConn = "conn"
)

var (
	errEmptyConn = errors.New("empty connection string")
	errEmptyArg  = errors.New("command line args doesn't exist")
)

func Run() *cli.App {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     flagConn,
			Usage:    "connection string",
			Required: true,
		},
	}

	app.Action = run

	return app
}

func run(ctx *cli.Context) error {
	conn, err := sqlx.Connect("pgx", ctx.String(flagConn))
	if err != nil {
		return err
	}

	defer conn.Close()

	accounts, err := readFile(ctx.Args().First())
	if err != nil {
		return err
	}

	if err = writeToDB(conn, accounts); err != nil {
		return err
	}

	return nil
}

func readFile(path string) ([]Account, error) {
	if path == "" {
		return nil, errEmptyArg
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	account := []Account{}
	if err = json.Unmarshal(data, &account); err != nil {
		return nil, err
	}

	return account, nil
}

func writeToDB(conn *sqlx.DB, accounts []Account) error {
	countries, cities, err := writeCountriesAndCities(conn, accounts)
	if err != nil {
		return err
	}

	if err = writeAccountsAndPersons(conn, accounts, countries, cities); err != nil {
		return err
	}

	return writeLikesAndInterests(conn, accounts)
}

func writeCountriesAndCities(conn *sqlx.DB, accs []Account) (countries, cities map[string]int32, err error) {
	lenCountries := int32(0)
	lenCities := int32(0)

	for _, acc := range accs {
		if acc.Country != nil {
			if _, ok := countries[*acc.Country]; !ok {
				countries[*acc.Country] = lenCountries
				lenCountries++
			}
		}

		if acc.City != nil {
			if _, ok := cities[*acc.City]; !ok {
				cities[*acc.City] = lenCities
				lenCities++
			}
		}
	}

	queryTotal := `INSERT INTO %s(id, name) VALUES %s;`

	queriesCity := make([]string, 0, lenCities)
	for name, id := range cities {
		queriesCity = append(queriesCity, fmt.Sprintf(`(%d, '%s')`, id, name))
	}

	if _, err = conn.Exec(fmt.Sprintf(queryTotal, "city", strings.Join(queriesCity, ", "))); err != nil {
		return
	}

	queriesCountry := make([]string, 0, lenCountries)
	for name, id := range cities {
		queriesCountry = append(queriesCountry, fmt.Sprintf(`(%d, '%s')`, id, name))
	}

	_, err = conn.Exec(fmt.Sprintf(queryTotal, "country", strings.Join(queriesCountry, ", ")))
	return
}

func writeAccountsAndPersons(conn *sqlx.DB, accs []Account, countries, cities map[string]int32) error {
	accounts := make([]models.Account, 0, len(accs))
	persons := make([]models.Person, 0, len(accs))

	for _, acc := range accs {
		accounts = append(accounts, newAccount(&acc))

		var countryID, cityID *int32
		if acc.Country != nil {
			countryID = int32Ptr(countries[*acc.Country])
		}

		if acc.City != nil {
			cityID = int32Ptr(cities[*acc.City])
		}

		persons = append(persons, newPerson(&acc, countryID, cityID))
	}

	queryAccountTotal := `INSERT INTO account(id, joined, status, premium_start, premium_end) VALUES %s;`
	queriesAccount := make([]string, 0, len(accounts))
	for _, a := range accounts {
		queriesAccount = append(
			queriesAccount,
			fmt.Sprintf(`(%d, %d, '%s', %d, %d)`, a.ID, a.Joined, a.Status, a.PremiumStart, a.PremiumEnd),
		)
	}

	if _, err := conn.Exec(fmt.Sprintf(queryAccountTotal, strings.Join(queriesAccount, ", "))); err != nil {
		return err
	}

	queryPersonTotal := `INSERT INTO person(account_id, email, sex, birth, name, surname, phone, country_id, city_id) VALUES %s;`
	queriesPerson := make([]string, 0, len(persons))
	for _, p := range persons {
		queriesPerson = append(
			queriesPerson,
			fmt.Sprintf(`(%d, '%s', '%s', %d, '%s', '%s', '%s', %d, %d)`, p.ID, p.Email, p.Sex, p.Birth, *p.Name, *p.Surname, *p.Phone, *p.CountryID, *p.CityID),
		)
	}

	if _, err := conn.Exec(fmt.Sprintf(queryPersonTotal, strings.Join(queriesPerson, ", "))); err != nil {
		return err
	}

	return nil
}

func writeLikesAndInterests(conn *sqlx.DB, accs []Account) error {
	likes := make([]models.Like, 0)
	interests := make([]models.Interest, 0)

	for _, acc := range accs {
		likes = append(likes, newLikes(&acc)...)
		interests = append(interests, newInterests(&acc)...)
	}

	queryLikeTotal := `INSERT INTO like(liker_id, likee_id, ts) VALUES %s;`
	queriesLike := make([]string, 0, len(likes))
	for _, like := range likes {
		queriesLike = append(queriesLike, fmt.Sprintf(`(%d, %d, %d)`, like.LikerID, like.LikeeID, like.Timestamp))
	}

	if _, err := conn.Exec(fmt.Sprintf(queryLikeTotal, strings.Join(queriesLike, ", "))); err != nil {
		return err
	}

	queryInterestTotal := `INSERT INTO interest(account_id, name) VALUES %s;`
	queriesInterest := make([]string, 0, len(interests))
	for _, i := range interests {
		queriesInterest = append(queriesInterest, fmt.Sprintf(`(%d, '%s')`, i.AccountID, i.Name))
	}

	if _, err := conn.Exec(fmt.Sprintf(queryInterestTotal, strings.Join(queriesInterest, ", "))); err != nil {
		return err
	}

	return nil
}

func newAccount(acc *Account) models.Account {
	a := models.Account{
		ID:     acc.ID,
		Joined: acc.Joined,
		Status: acc.Status,
	}

	if acc.Premium != nil {
		a.PremiumStart = &acc.Premium.Start
		a.PremiumEnd = &acc.Premium.End
	}

	return a
}

func newPerson(acc *Account, countryID, cityID *int32) models.Person {
	return models.Person{
		ID:        acc.ID,
		Email:     acc.Email,
		Sex:       acc.Sex,
		Birth:     acc.Birth,
		Name:      acc.Name,
		Surname:   acc.Surname,
		Phone:     acc.Phone,
		CountryID: countryID,
		CityID:    cityID,
	}
}

func newLikes(acc *Account) []models.Like {
	likes := make([]models.Like, 0, len(acc.Likes))

	for _, like := range acc.Likes {
		likes = append(likes, models.Like{
			LikerID:   acc.ID,
			LikeeID:   like.UserID,
			Timestamp: like.Timestamp,
		})
	}

	return likes
}

func newInterests(acc *Account) []models.Interest {
	interests := make([]models.Interest, 0, len(acc.Interests))

	for _, interest := range acc.Interests {
		interests = append(interests, models.Interest{
			AccountID: acc.ID,
			Name:      interest,
		})
	}

	return interests
}

func int32Ptr(val int32) *int32 {
	return &val
}
