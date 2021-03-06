package dataloader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
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

	if err = writeAccounts(conn, accounts, countries, cities); err != nil {
		return err
	}

	return writeLikesAndInterests(conn, accounts)
}

func writeCountriesAndCities(conn *sqlx.DB, accs []Account) (countries, cities map[string]uuid.UUID, err error) {
	countries = make(map[string]uuid.UUID)
	cities = make(map[string]uuid.UUID)

	for _, acc := range accs {
		if acc.Country != nil {
			countries[*acc.Country] = uuid.New()
		}

		if acc.City != nil {
			cities[*acc.City] = uuid.New()
		}
	}

	queryTotal := `INSERT INTO %s(id, name) VALUES %s;`

	queriesCity := make([]string, 0, len(cities))
	for name, id := range cities {
		queriesCity = append(queriesCity, fmt.Sprintf(`('%s'::uuid, '%s')`, id, name))
	}

	if _, err = conn.Exec(fmt.Sprintf(queryTotal, "city", strings.Join(queriesCity, ", "))); err != nil {
		return
	}

	queriesCountry := make([]string, 0, len(countries))
	for name, id := range countries {
		queriesCountry = append(queriesCountry, fmt.Sprintf(`('%s'::uuid, '%s')`, id, name))
	}

	_, err = conn.Exec(fmt.Sprintf(queryTotal, "country", strings.Join(queriesCountry, ", ")))
	return
}

func writeAccounts(conn *sqlx.DB, accs []Account, countries, cities map[string]uuid.UUID) error {
	accounts := make([]domain.AccountModel, 0, len(accs))

	for _, acc := range accs {

		var countryID, cityID *uuid.UUID
		if acc.Country != nil {
			countryID = ptrUUID(countries[*acc.Country])
		}

		if acc.City != nil {
			cityID = ptrUUID(cities[*acc.City])
		}

		accounts = append(accounts, newAccount(&acc, countryID, cityID))
	}

	queryAccountTotal := `INSERT INTO account(id, email, sex, status, birth, name, surname, phone, country_id, city_id, joined, prem_start, prem_end) VALUES %s;`
	queriesAccount := make([]string, 0, len(accounts))
	for _, a := range accounts {
		query := fmt.Sprintf(
			`(%d, '%s', '%s', '%s', %s, %s, %s, %s, %s, %s, %s, %s, %s)`,
			a.ID, a.Email, a.Sex, a.Status,
			nullableTimestamp(&a.Birth),
			nullableString(a.Name),
			nullableString(a.Surname),
			nullableString(a.Phone),
			nullableUUID(a.CountryID),
			nullableUUID(a.CityID),
			nullableTimestamp(&a.Joined),
			nullableTimestamp(a.PremiumStart),
			nullableTimestamp(a.PremiumEnd),
		)

		queriesAccount = append(queriesAccount, query)
	}

	query := strings.Join(queriesAccount, ", ")
	if _, err := conn.Exec(fmt.Sprintf(queryAccountTotal, query)); err != nil {
		log.Println(query)
		return err
	}

	return nil
}

func writeLikesAndInterests(conn *sqlx.DB, accs []Account) error {
	likes := make([]domain.LikeModel, 0)
	interests := make([]domain.InterestModel, 0)

	for _, acc := range accs {
		likes = append(likes, newLikes(&acc)...)
		interests = append(interests, newInterests(&acc)...)
	}

	queryLikeTotal := `INSERT INTO likes(liker_id, likee_id, ts) VALUES %s;`
	queriesLike := make([]string, 0, len(likes))
	for _, like := range likes {
		queriesLike = append(queriesLike, fmt.Sprintf(`(%d, %d, %s)`, like.LikerID, like.LikeeID, nullableTimestamp(&like.Timestamp)))
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

func newAccount(a *Account, countryID, cityID *uuid.UUID) domain.AccountModel {
	account := domain.AccountModel{
		ID:        a.ID,
		Status:    a.Status,
		Email:     a.Email,
		Sex:       a.Sex,
		Birth:     int64PtrToTimestamp(&a.Birth),
		Name:      a.Name,
		Surname:   a.Surname,
		Phone:     a.Phone,
		CountryID: countryID,
		CityID:    cityID,
		Joined:    int64PtrToTimestamp(&a.Joined),
	}

	if a.Premium != nil {
		start := int64PtrToTimestamp(&a.Premium.Start)
		account.PremiumStart = &start
		end := int64PtrToTimestamp(&a.Premium.End)
		account.PremiumEnd = &end
	}

	return account
}

func newLikes(acc *Account) []domain.LikeModel {
	likes := make([]domain.LikeModel, 0, len(acc.Likes))

	for _, like := range acc.Likes {
		likes = append(likes, domain.LikeModel{
			LikerID:   acc.ID,
			LikeeID:   like.UserID,
			Timestamp: int64PtrToTimestamp(&like.Timestamp),
		})
	}

	return likes
}

func newInterests(acc *Account) []domain.InterestModel {
	interests := make([]domain.InterestModel, 0, len(acc.Interests))

	for _, interest := range acc.Interests {
		interests = append(interests, domain.InterestModel{
			AccountID: acc.ID,
			Name:      interest,
		})
	}

	return interests
}

func ptrInt32(val int32) *int32 {
	return &val
}

func ptrUUID(val uuid.UUID) *uuid.UUID {
	return &val
}

func int64PtrToTimestamp(val *int64) time.Time {
	if val == nil {
		return time.Time{}
	}

	return time.Unix(*val, 0)
}

func nullableTimestamp(ts *time.Time) string {
	if ts == nil || *ts == (time.Time{}) {
		return "null"
	}

	return fmt.Sprintf(`'%s'::timestamp`, ts.Format(timestampLayout))
}

func nullableString(str *string) string {
	if str == nil {
		return "null"
	}

	return fmt.Sprintf(`'%s'`, *str)
}

func nullableUUID(val *uuid.UUID) string {
	if val == nil {
		return "null"
	}

	return fmt.Sprintf(`'%s'::uuid`, val.String())
}
