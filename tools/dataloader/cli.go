package dataloader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli/v2"

	"accounts/domain"
)

const (
	flagConn = "conn"

	timestampLayout = "2006-01-02 15:04:05"
	nullTime        = "0000-00-00 00:00:00"
)

var (
	errEmptyConn = errors.New("empty connection string")
	errEmptyArg  = errors.New("command line args doesn't exist")
)

func Run() error {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     flagConn,
			Usage:    "connection string",
			Required: true,
		},
	}

	app.Action = run

	return app.Run(os.Args)
}

func run(ctx *cli.Context) error {
	connStr := ctx.String(flagConn)
	if connStr == "" {
		return errEmptyConn
	}

	conn, err := sqlx.Connect("pgx", connStr)
	if err != nil {
		return err
	}

	defer conn.Close()

	accounts, err := readFile(ctx.Args().Slice())
	if err != nil {
		return err
	}

	if err = writeToDB(conn, accounts); err != nil {
		return err
	}

	return nil
}

func readFile(paths []string) ([]Account, error) {
	if len(paths) == 0 {
		return nil, errEmptyArg
	}

	result := make([]Account, 0)
	for _, path := range paths {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		var accounts Accounts
		if err = json.Unmarshal(data, &accounts); err != nil {
			return nil, err
		}

		result = append(result, accounts.Accounts...)
	}

	return result, nil
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
	countries = make(map[string]int32)
	cities = make(map[string]int32)

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
	for name, id := range countries {
		queriesCountry = append(queriesCountry, fmt.Sprintf(`(%d, '%s')`, id, name))
	}

	_, err = conn.Exec(fmt.Sprintf(queryTotal, "country", strings.Join(queriesCountry, ", ")))
	return
}

func writeAccountsAndPersons(conn *sqlx.DB, accs []Account, countries, cities map[string]int32) error {
	accounts := make([]domain.Account, 0, len(accs))
	persons := make([]domain.Person, 0, len(accs))

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

	queryAccountTotal := `INSERT INTO account(id, joined, status, prem_start, prem_end) VALUES %s;`
	queriesAccount := make([]string, 0, len(accounts))
	for _, a := range accounts {
		query := fmt.Sprintf(`(%d, %s, '%s', %s, %s)`, a.ID, nullableTimestamp(a.Joined), a.Status, nullableTimestamp(a.PremiumStart), nullableTimestamp(a.PremiumEnd))
		queriesAccount = append(queriesAccount, query)
	}

	query := strings.Join(queriesAccount, ", ")
	if _, err := conn.Exec(fmt.Sprintf(queryAccountTotal, query)); err != nil {
		log.Println(query)
		return err
	}

	queryPersonTotal := `INSERT INTO person(account_id, email, sex, birth, name, surname, phone, country_id, city_id) VALUES %s;`
	queriesPerson := make([]string, 0, len(persons))
	for _, p := range persons {
		queriesPerson = append(
			queriesPerson,
			fmt.Sprintf(
				`(%d, '%s', '%s', %s, %s, %s, %s, %s, %s)`,
				p.ID, p.Email, p.Sex,
				nullableTimestamp(p.Birth),
				nullableString(p.Name),
				nullableString(p.Surname),
				nullableString(p.Phone),
				nullableInt32(p.CountryID),
				nullableInt32(p.CityID),
			),
		)
	}

	query = strings.Join(queriesPerson, ", ")
	if _, err := conn.Exec(fmt.Sprintf(queryPersonTotal, query)); err != nil {
		log.Println(query)
		return err
	}

	return nil
}

func writeLikesAndInterests(conn *sqlx.DB, accs []Account) error {
	likes := make([]domain.Like, 0)
	interests := make([]domain.Interest, 0)

	for _, acc := range accs {
		likes = append(likes, newLikes(&acc)...)
		interests = append(interests, newInterests(&acc)...)
	}

	queryLikeTotal := `INSERT INTO likes(liker_id, likee_id, ts) VALUES %s;`
	queriesLike := make([]string, 0, len(likes))
	for _, like := range likes {
		queriesLike = append(queriesLike, fmt.Sprintf(`(%d, %d, %s)`, like.LikerID, like.LikeeID, nullableTimestamp(like.Timestamp)))
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

func newAccount(acc *Account) domain.Account {
	a := domain.Account{
		ID:     acc.ID,
		Joined: int64PtrToTimestamp(&acc.Joined),
		Status: acc.Status,
	}

	if acc.Premium != nil {
		a.PremiumStart = int64PtrToTimestamp(&acc.Premium.Start)
		a.PremiumEnd = int64PtrToTimestamp(&acc.Premium.End)
	} else {
		a.PremiumStart = nullTime
		a.PremiumEnd = nullTime
	}

	return a
}

func newPerson(acc *Account, countryID, cityID *int32) domain.Person {
	return domain.Person{
		ID:        acc.ID,
		Email:     acc.Email,
		Sex:       acc.Sex,
		Birth:     int64PtrToTimestamp(&acc.Birth),
		Name:      acc.Name,
		Surname:   acc.Surname,
		Phone:     acc.Phone,
		CountryID: countryID,
		CityID:    cityID,
	}
}

func newLikes(acc *Account) []domain.Like {
	likes := make([]domain.Like, 0, len(acc.Likes))

	for _, like := range acc.Likes {
		likes = append(likes, domain.Like{
			LikerID:   acc.ID,
			LikeeID:   like.UserID,
			Timestamp: int64PtrToTimestamp(&like.Timestamp),
		})
	}

	return likes
}

func newInterests(acc *Account) []domain.Interest {
	interests := make([]domain.Interest, 0, len(acc.Interests))

	for _, interest := range acc.Interests {
		interests = append(interests, domain.Interest{
			AccountID: acc.ID,
			Name:      interest,
		})
	}

	return interests
}

func int32Ptr(val int32) *int32 {
	return &val
}

func int64PtrToTimestamp(val *int64) string {
	if val == nil {
		return nullTime
	}

	return time.Unix(*val, 0).Format(timestampLayout)
}

func nullableTimestamp(timestamp string) string {
	if timestamp == nullTime {
		return "null"
	}

	return fmt.Sprintf(`'%s'::timestamp`, timestamp)
}

func nullableString(str *string) string {
	if str == nil {
		return "null"
	}

	return fmt.Sprintf(`'%s'`, *str)
}

func nullableInt32(val *int32) string {
	if val == nil {
		return "null"
	}

	return strconv.Itoa(int(*val))
}
