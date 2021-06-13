package repository

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
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
	conn *pgxpool.Pool
}

func New(conn *pgxpool.Pool) *Repository {
	return &Repository{
		conn: conn,
	}
}

func (r *Repository) FilterAccounts(ctx context.Context, f *Filter) (*domain.AccountsOut, error) {
	sql, values, err := buildAccountSearchQuery(f)
	if err != nil {
		return nil, err
	}

	log.Println(sql, values)

	rows, err := r.conn.Query(ctx, sql, values...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var acc domain.AccountOut
	scanFields := make([]interface{}, 0, len(f.Columns()))
	scanFields = append(scanFields, &acc.ID)
	scanFields = append(scanFields, &acc.Email)
	for column := range f.Columns() {
		switch column {
		case AccountID, AccountEmail:
			continue
		case AccountSex:
			scanFields = append(scanFields, &acc.Sex)
		case AccountStatus:
			scanFields = append(scanFields, &acc.Status)
		case AccountBirth:
			scanFields = append(scanFields, &acc.Birth)
		case AccountFirstname:
			scanFields = append(scanFields, &acc.Fname)
		case AccountSurname:
			scanFields = append(scanFields, &acc.Sname)
		case AccountPhone:
			scanFields = append(scanFields, &acc.Phone)
		case CountryName:
			scanFields = append(scanFields, &acc.Country)
		case CityName:
			scanFields = append(scanFields, &acc.City)
		case AccountPremStart:
			scanFields = append(scanFields, &acc.Premium)
		case InterestName, LikesLikeeID:
			continue
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

	return &domain.AccountsOut{Accounts: accounts}, nil
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

func (r *Repository) AddLikes(ctx context.Context, likes *domain.LikesInput) error {
	return r.tryInsertLikes(ctx, likes.LikeModels(), nil)
}

func (r *Repository) insertAccount(ctx context.Context, a *domain.AccountModel, tx pgx.Tx) error {
	if a == nil {
		return errNilModel
	}

	sql, values, err := buildAccountInsertQuery(*a)
	if err != nil {
		return err
	}

	log.Println(sql, values)

	return tx.QueryRow(ctx, sql, values...).Scan(&a.ID)
}

func (r *Repository) tryInsertCity(ctx context.Context, c *domain.CityModel, tx pgx.Tx) (id uuid.UUID, err error) {
	if c == nil {
		return
	}

	sql, values, err := buildCityInsertQuery(*c)
	if err != nil {
		return
	}

	log.Println(sql, values)

	err = tx.QueryRow(ctx, sql, values...).Scan(&id)
	return
}

func (r *Repository) tryInsertCountry(ctx context.Context, c *domain.CountryModel, tx pgx.Tx) (id uuid.UUID, err error) {
	if c == nil {
		return
	}

	sql, values, err := buildCountryInsertQuery(*c)
	if err != nil {
		return
	}

	log.Println(sql, values)

	err = tx.QueryRow(ctx, sql, values...).Scan(&id)
	return
}

// TODO: rework bulk insert
func (r *Repository) tryInsertLikes(ctx context.Context, likes []domain.LikeModel, tx pgx.Tx) (err error) {
	if likes == nil || len(likes) == 0 {
		return
	}

	if tx == nil {
		_, err = r.conn.CopyFrom(
			ctx,
			pgx.Identifier{TableLike},
			[]string{LikesLikerID, LikesLikeeID, LikesTimestamp},
			pgx.CopyFromSlice(len(likes), func(i int) ([]interface{}, error) {
				return []interface{}{likes[i].LikerID, likes[i].LikeeID, likes[i].Timestamp}, nil
			}),
		)

		return
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{TableLike},
		[]string{LikesLikerID, LikesLikeeID, LikesTimestamp},
		pgx.CopyFromSlice(len(likes), func(i int) ([]interface{}, error) {
			return []interface{}{likes[i].LikerID, likes[i].LikeeID, likes[i].Timestamp}, nil
		}),
	)

	return
}

// TODO: rework bulk insert
func (r *Repository) tryInsertInterests(ctx context.Context, interests []domain.InterestModel, tx pgx.Tx) (err error) {
	if interests == nil || len(interests) == 0 {
		return
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{TableInterest},
		[]string{InterestAccountID, InterestName},
		pgx.CopyFromSlice(len(interests), func(i int) ([]interface{}, error) {
			return []interface{}{interests[i].AccountID, interests[i].Name}, nil
		}),
	)

	return
}

func (r *Repository) updateAccount(ctx context.Context, a domain.AccountUpdate, cityID, countryID uuid.UUID, tx pgx.Tx) error {
	sql, values, err := buildAccountUpdateQuery(a, cityID, countryID)
	result, err := tx.Exec(ctx, sql, values...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errNotAffected
	}

	return nil
}
