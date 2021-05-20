package repository

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"accounts/domain"
)

type Repository struct {
	conn *pgxpool.Conn
}

func (r *Repository) FilterAccounts(ctx context.Context, filter string) ([]domain.AccountOut, error) {
	return nil, nil
}

func (r *Repository) AddAccount(ctx context.Context, a domain.AccountInput) error {

	return nil
}

func (r *Repository) insertAccount(ctx context.Context, a *domain.AccountTable, tx pgx.Tx) error {
	return tx.QueryRow(
		ctx,
		`INSERT INTO account(id, joined, status, prem_start, prem_end) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		a.ID, a.Joined, a.Status, a.PremiumStart, a.PremiumEnd,
	).Scan(&a.ID)
}

func (r *Repository) insertCity(ctx context.Context, c *domain.CityTable, tx pgx.Tx) error {
	return tx.QueryRow(
		ctx,
		`INSERT INTO city(name) VALUES($1) ON CONFLICT(name) DO NOTHING RETURNING id`,
		c.Name,
	).Scan(&c.ID)
}

func (r *Repository) insertCountry(ctx context.Context, c *domain.CountryTable, tx pgx.Tx) error {
	return tx.QueryRow(
		ctx,
		`INSERT INTO country(name) VALUES($1) ON CONFLICT(name) DO NOTHING RETURNING id`,
		c.Name,
	).Scan(&c.ID)
}

func (r *Repository) insertPerson(ctx context.Context, p *domain.PersonTable, tx pgx.Tx) error {
	return tx.QueryRow(
		ctx,
		`INSERT INTO person(account_id, email, sex, birth, name, surname, phone, country_id, city_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING account_id`,
		p.ID, p.Email, p.Sex, p.Birth, p.Name, p.Surname, p.Phone, p.CountryID, p.CityID,
	).Scan(&p.ID)
}

func (r *Repository) insertLikes() error {
	return nil
}

func (r *Repository) insertInterests() error {
	return nil
}

func (r *Repository) UpdateAccount(ctx context.Context, a domain.AccountUpdate) error {
	return nil
}

func (r *Repository) AddLikes(ctx context.Context, likes []domain.LikeInput) error {
	return nil
}
