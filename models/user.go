package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID              uuid.UUID `json:"id"`
	CreatedAt       time.Time `json:"_"`
	UpdatedAt       time.Time `json:"_"`
	Email           string    `json:"email"`
	PasswordHash    string    `json:"-"`
	Password        string    `json:"password"`
	PasswordConfirm string    `json:"password_confirm"`
}

func (u *User) Register(conn *pgx.Conn) error {
	//validate password and email
	if len(u.Password) < 4 || len(u.PasswordConfirm) < 4 {
		return fmt.Errorf("Password must be as least 4 character long.")
	}

	if u.Password != u.PasswordConfirm {
		return fmt.Errorf("Passwords do not match.")
	}

	if len(u.Email) < 4 {
		return fmt.Errorf("Email must be as least 4 character long.")
	}

	u.Email = strings.ToLower(u.Email)
	row := conn.QueryRow(context.Background(), "SELECT id FROM user_account WHERE email= $1", u.Email)
	userLookup := User{}
	err := row.Scan(&userLookup)
	if err != pgx.ErrNoRows {
		return fmt.Errorf("A user with that email already exists")
	}

	//Hash password
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("There was an error creating your account")
	}
	u.PasswordHash = string(pwdHash)

	//Save object into DB
	//now := time.Now()
	//_, err = conn.Exec(context.Background(), "INSERT INTO user_account (created_at, updated_at, email, password_hash) VALUES ($1, $2, $3, $4)", now, now, u.Email, u.PasswordHash)
	return err
}

//GetAuthToken returns the auth token to be used
func (u *User) GetAuthToken() (authToken string, err error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = u.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	authToken, err = token.SignedString([]byte("secretfortoken"))
	return

}
