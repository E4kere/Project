package models

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Gun struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Manufacturer string    `json:"manufacturer"`
	Price        float64   `json:"price"`
	Damage       int       `json:"damage"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type GunModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m *GunModel) Insert(gun *Gun) error {
	query := `
    INSERT INTO guns (name, manufacturer, price, damage) 
    VALUES ($1, $2, $3, $4) 
    RETURNING id, created_at, updated_at
    `
	args := []interface{}{gun.Name, gun.Manufacturer, gun.Price, gun.Damage}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&gun.ID, &gun.CreatedAt, &gun.UpdatedAt)
}

func (m *GunModel) Get(id int) (*Gun, error) {
	query := `
    SELECT id, created_at, updated_at, name, manufacturer, price, damage
    FROM guns
    WHERE id = $1
    `
	var gun Gun
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&gun.ID, &gun.CreatedAt, &gun.UpdatedAt, &gun.Name, &gun.Manufacturer, &gun.Price, &gun.Damage)
	if err != nil {
		return nil, err
	}
	return &gun, nil
}

func (m *GunModel) Update(gun *Gun) error {
	query := `
    UPDATE guns
    SET name = $1, manufacturer = $2, price = $3, damage = $4
    WHERE id = $5
    RETURNING updated_at
    `
	args := []interface{}{gun.Name, gun.Manufacturer, gun.Price, gun.Damage, gun.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&gun.UpdatedAt)
}

func (m *GunModel) Delete(id int) error {
	query := `
    DELETE FROM guns
    WHERE id = $1
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
