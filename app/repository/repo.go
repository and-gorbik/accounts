package repository

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"accounts/domain"
	"accounts/util"
)

var (
	errNilModel    = errors.New("nil model (input model wasn't validated probably)")
	errNotAffected = errors.New("not affected")
)

type Repository struct {
	conn *pgxpool.Conn
}

func (r *Repository) FilterAccounts(ctx context.Context, filter string) ([]domain.AccountOut, error) {
	return nil, nil
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

	if err = r.insertAccount(ctx, a.AccountModel(), tx); err != nil {
		return err
	}

	if err = r.insertPerson(ctx, a.PersonModel(&cityID, &countryID), tx); err != nil {
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
		`INSERT INTO account(id, joined, prem_start, prem_end) VALUES ($1, $2, $3, $4) RETURNING id`,
		a.ID, a.Joined, a.PremiumStart, a.PremiumEnd,
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

func (r *Repository) insertPerson(ctx context.Context, p *domain.PersonModel, tx pgx.Tx) error {
	if p == nil {
		return errNilModel
	}

	return tx.QueryRow(
		ctx,
		`INSERT INTO person(account_id, status, email, sex, birth, name, surname, phone, country_id, city_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING account_id`,
		p.ID, p.Status, p.Email, p.Sex, p.Birth, p.Name, p.Surname, p.Phone, p.CountryID, p.CityID,
	).Scan(&p.ID)
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

	if err = r.updatePerson(ctx, a, cityID, countryID, tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *Repository) updatePerson(ctx context.Context, a domain.AccountUpdate, cityID, countryID int32, tx pgx.Tx) error {
	fields := []string{}
	values := []interface{}{}

	if a.Email != nil {
		fields = append(fields, "email=$"+strconv.Itoa(len(fields)+1))
		values = append(values, string(*a.Email))
	}
	if a.Birth != nil {
		fields = append(fields, "birth=$"+strconv.Itoa(len(fields)+1))
		values = append(values, *util.TimestampToDatetime((*int64)(a.Birth)))
	}
	if a.Status != nil {
		fields = append(fields, "status=$"+strconv.Itoa(len(fields)+1))
		values = append(values, string(*a.Status))
	}
	fields = append(fields, "city_id=$"+strconv.Itoa(len(fields)+1))
	values = append(values, cityID)
	fields = append(fields, "country_id=$"+strconv.Itoa(len(fields)+1))
	values = append(values, countryID)
	values = append(values, a.ID)

	var builder strings.Builder
	builder.WriteString("UPDATE person SET ")
	builder.WriteString(strings.Join(fields, ", "))
	builder.WriteString(" WHERE id = $" + strconv.Itoa(len(fields)+1))

	result, err := tx.Exec(ctx, builder.String(), values...)
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
