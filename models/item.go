package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
)

type Item struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"_"`
	UpdatedAt    time.Time `json:"_"`
	Title        string    `json:"title"`
	Notes        string    `json:"note"`
	SellerID     uuid.UUID `json:"seller"`
	PriceInCents int64     `json:"price_in_cents"`
}

func (i *Item) Create(conn *pgx.Conn, userID string) error {
	i.Title = strings.Trim(i.Title, " ")
	if len(i.Title) < 1 {
		return fmt.Errorf("Title must not be empty")
	}
	if i.PriceInCents < 0 {
		i.PriceInCents = 0
	}

	now := time.Now()
	row := conn.QueryRow(context.Background(),
		"INSERT INTO item (title, note, seller_id, price_in_cents, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, seller_id",
		i.Title, i.Notes, userID, i.PriceInCents, now, now)

	err := row.Scan(&i.ID, &i.SellerID)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("There was an error creating the item")
	}

	return nil
}

func GetAllItems(conn *pgx.Conn) ([]Item, error) {
	rows, err := conn.Query(context.Background(), "SELECT id, title, note, seller_id, price_in_cents FROM item")
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("Error getting items")
	}
	var items []Item
	for rows.Next() {
		item := Item{}
		err = rows.Scan(&item.ID, &item.Title, &item.Notes, &item.SellerID, &item.PriceInCents)
		if err != nil {
			fmt.Println(err)
			continue
		}
		items = append(items, item)
	}
	return items, nil
}