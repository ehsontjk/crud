package security

import (
	"time"
	"encoding/hex"
	"crypto/rand"
	"golang.org/x/crypto/bcrypt"
	"github.com/jackc/pgx/v4"
	"errors"
	"context"
	"log"
    "github.com/jackc/pgx/v4/pgxpool"
)

var(

	ErrNoSuchUser = errors.New("no such user")
	
	ErrInvalidPassword = errors.New("invalid password")
	
	ErrInternal = errors.New("internal error")

	ErrExpireToken = errors.New("token expired")
)




type Service struct {
	db *pgxpool.Pool
}


func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}


func (s *Service) Auth(login, password string) bool {

	
	sqlStatement := `select login, password from managers where login=$1 and password=$2`

	
	err := s.db.QueryRow(context.Background(), sqlStatement, login, password).Scan(&login, &password)
	
	if err != nil {
		log.Print(err)
		return false
	}
	
	return true
}


func (s *Service) TokenForCustomer(ctx context.Context, phone, password string)(string, error){

	var hash string
	var id int64

	err := s.db.QueryRow(ctx, "select id, password from customers where phone = $1", phone).Scan(&id, &hash)

	if err == pgx.ErrNoRows{
		return "", ErrNoSuchUser
	}
	if err != nil{
		return "", ErrInternal
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil{
		return "", ErrInvalidPassword
	}

	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil{
		return "", ErrInternal
	}

	token := hex.EncodeToString(buffer)
	_, err = s.db.Exec(ctx, "insert into customers_tokens(token, customer_id) values($1, $2)", token, id)
	if err != nil{
		return "", ErrInternal
	}

	return token, nil

}


func (s *Service) AuthenticateCustomer(ctx context.Context, token string)(int64, error){
	var id int64
	var expire time.Time
	err := s.db.QueryRow(ctx, "select customer_id, expire from customers_tokens where token=$1", token).Scan(&id, &expire)
	if err == pgx.ErrNoRows{
		return 0, ErrNoSuchUser
	}
	if err != nil{
		return 0, ErrInternal
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	end := expire.Format("2006-01-02 15:04:05")

	if now > end {
		return 0, ErrExpireToken
	}

	return id, nil
}