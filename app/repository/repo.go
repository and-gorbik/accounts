package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"accounts/domain"
)

var (
	errNilModel     = errors.New("nil model (input model wasn't validated probably)")
	errNotAffected  = errors.New("not affected")
	errInvalidField = errors.New("invalid field")
)

type Repository struct {
	conn *pgxpool.Conn
}

func (r *Repository) FilterAccounts(ctx context.Context, filter Filter) ([]domain.AccountOut, error) {
	rows, err := r.conn.Query(ctx, buildAccountSearchQuery(filter))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var acc domain.AccountOut
	scanFields := make([]interface{}, 0, len(filter.Fields))
	for field := range filter.Fields {
		switch field {
		case "id":
			scanFields = append(scanFields, &acc.ID)
		case "email":
			scanFields = append(scanFields, &acc.Email)
		case "sex":
			scanFields = append(scanFields, &acc.Sex)
		case "status":
			scanFields = append(scanFields, &acc.Status)
		case "birth":
			scanFields = append(scanFields, &acc.Birth)
		case "fname":
			scanFields = append(scanFields, &acc.Fname)
		case "sname":
			scanFields = append(scanFields, &acc.Sname)
		case "phone":
			scanFields = append(scanFields, &acc.Phone)
		case "country":
			scanFields = append(scanFields, &acc.Country)
		case "city":
			scanFields = append(scanFields, &acc.City)
		case "premium":
			scanFields = append(scanFields, &acc.Premium)
		default:
			return nil, errInvalidField
		}
	}

	accounts := []domain.AccountOut{}
	for rows.Next() {
		if err := rows.Scan(scanFields...); err != nil {
			return nil, err
		}

		accounts = append(accounts, acc)
	}

	return accounts, nil
}

func (r *Repository) AddAccount(ctx context.Context, a domain.AccountInput) error {
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	cityID, err := r.tryInsertCity(ctx, a.CityModel(), tx)
	if err != nil {
		return err
	}

	countryID, err := r.tryInsertCountry(ctx, a.CountryModel(), tx)
	if err != nil {
		return err
	}

	if err = r.insertAccount(ctx, a.AccountModel(&cityID, &countryID), tx); err != nil {
		return err
	}

	if err = r.tryInsertLikes(ctx, a.LikeModels(), tx); err != nil {
		return err
	}

	if err = r.tryInsertInterests(ctx, a.InterestModels(), tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *Repository) insertAccount(ctx context.Context, a *domain.AccountModel, tx pgx.Tx) error {
	if a == nil {
		return errNilModel
	}

	return tx.QueryRow(
		ctx,
		`INSERT INTO account(id, status, email, sex, birth, name, surname, phone, country_id, city_id, joined, prem_start, prem_end)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING account_id`,
		a.ID, a.Status, a.Email, a.Sex, a.Birth, a.Name, a.Surname, a.Phone, a.CountryID, a.CityID, a.Joined, a.PremiumStart, a.PremiumEnd,
	).Scan(&a.ID)
}

func (r *Repository) tryInsertCity(ctx context.Context, c *domain.CityModel, tx pgx.Tx) (id int32, err error) {
	if c == nil {
		return
	}

	err = tx.QueryRow(
		ctx,
		`INSERT INTO city(name) VALUES($1) ON CONFLICT(name) DO NOTHING RETURNING id`,
		c.Name,
	).Scan(&id)
	return
}

func (r *Repository) tryInsertCountry(ctx context.Context, c *domain.CountryModel, tx pgx.Tx) (id int32, err error) {
	if c == nil {
		return
	}

	err = tx.QueryRow(
		ctx,
		`INSERT INTO country(name) VALUES($1) ON CONFLICT(name) DO NOTHING RETURNING id`,
		c.Name,
	).Scan(&id)
	return
}

func (r *Repository) tryInsertLikes(ctx context.Context, likes []domain.LikeModel, tx pgx.Tx) (err error) {
	if likes == nil || len(likes) == 0 {
		return
	}

	if tx == nil {
		_, err = r.conn.CopyFrom(
			ctx,
			pgx.Identifier{"likes"},
			[]string{"liker_id", "likee_id", "ts"},
			pgx.CopyFromSlice(len(likes), func(i int) ([]interface{}, error) {
				return []interface{}{likes[i].LikerID, likes[i].LikeeID, likes[i].Timestamp}, nil
			}),
		)

		return
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"likes"},
		[]string{"liker_id", "likee_id", "ts"},
		pgx.CopyFromSlice(len(likes), func(i int) ([]interface{}, error) {
			return []interface{}{likes[i].LikerID, likes[i].LikeeID, likes[i].Timestamp}, nil
		}),
	)

	return
}

func (r *Repository) tryInsertInterests(ctx context.Context, interests []domain.InterestModel, tx pgx.Tx) (err error) {
	if interests == nil || len(interests) == 0 {
		return
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"interests"},
		[]string{"account_id", "name"},
		pgx.CopyFromSlice(len(interests), func(i int) ([]interface{}, error) {
			return []interface{}{interests[i].AccountID, interests[i].Name}, nil
		}),
	)

	return
}

func (r *Repository) UpdateAccount(ctx context.Context, a domain.AccountUpdate) error {
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	cityID, err := r.tryInsertCity(ctx, a.CityModel(), tx)
	if err != nil {
		return err
	}

	countryID, err := r.tryInsertCountry(ctx, a.CountryModel(), tx)
	if err != nil {
		return err
	}

	if err = r.updateAccount(ctx, a, cityID, countryID, tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *Repository) updateAccount(ctx context.Context, a domain.AccountUpdate, cityID, countryID int32, tx pgx.Tx) error {
	sql, values := buildAccountUpdateQuery(a, cityID, countryID)
	result, err := tx.Exec(ctx, sql, values...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errNotAffected
	}

	return nil
}

func (r *Repository) AddLikes(ctx context.Context, likes *domain.LikesInput) error {
	return r.tryInsertLikes(ctx, likes.LikeModels(), nil)
}
