package models

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// User ...
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}

// Users ...
type Users struct {
	pool *pgxpool.Pool
	sync.Mutex
}

// NewUsers ...
func NewUsers(pool *pgxpool.Pool) *Users {
	return &Users{pool: pool}
}

// Register ...
func (u *Users) Register(ctx context.Context, user *User) {
	u.Lock()
	defer u.Unlock()

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Phone), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	hasher := md5.New()
	hasher.Write(hash)
	token := hex.EncodeToString(hasher.Sum(nil))

	user.ID = uuid.New().String()
	user.Token = token
	user.CreatedAt = time.Now()

	err = u.pool.QueryRow(ctx, `
		INSERT INTO users (id, name, phone, token, created_at) VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, phone, token, created_at;`,
		user.ID, user.Name, user.Phone, user.Token, user.CreatedAt).Scan(&user.ID, &user.Name, &user.Phone, &user.Token, &user.CreatedAt)

	if err != nil {
		log.Print(err)

		return
	}

	if errors.Is(err, pgx.ErrNoRows) {
		log.Print("No rows")
		return
	}

	if err != nil {
		log.Print(err)
		return
	}
}
