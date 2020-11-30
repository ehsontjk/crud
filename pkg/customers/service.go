package customers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)


var ErrNotFound = errors.New("item not found")

var ErrInternal = errors.New("internal error")


type Service struct {
	db *sql.DB
}


func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}


type Customer struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	Phone string `json:"phone"`
	Active bool `json:"active"`
	Created time.Time `json:"created"`
}


func (s *Service) All(ctx context.Context) (c []*Customer, err error) {

	sql := `select * from customers`

	rows, err := s.db.QueryContext(ctx, sql)
	if err != nil {
	return nil, err
	}
	defer rows.Close()

	for rows.Next() {
	item := &Customer{}
	err := rows.Scan(
	&item.ID,
	&item.Name,
	&item.Phone,
	&item.Active,
	&item.Created,
	)
	if err != nil {
	log.Println(err)
	}
	c = append(c, item)
	}

	return c, nil
}


func (s *Service) AllActive(ctx context.Context) (c []*Customer, err error) {

	sqlStatement := `select * from customers where active=true`

	rows, err := s.db.QueryContext(ctx, sqlStatement)
	if err != nil {
	return nil, err
	}
	defer rows.Close()

	for rows.Next() {
	item := &Customer{}
	err := rows.Scan(
	&item.ID,
	&item.Name,
	&item.Phone,
	&item.Active,
	&item.Created,
	)
	if err != nil {
	log.Println(err)
	}
	c = append(c, item)
	}

	return c, nil
}


func (s *Service) ByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	sqls := `select * from customers where id=$1`
	err := s.db.QueryRowContext(ctx, sqls, id).Scan(
	&item.ID,
	&item.Name,
	&item.Phone,
	&item.Active,
	&item.Created)

	if err == sql.ErrNoRows {
	return nil, ErrNotFound
	}
	if err != nil {
	log.Print(err)
	return nil, ErrInternal
	}
	return item, nil

}


func (s *Service) ChangeActive(ctx context.Context, id int64, active bool) (*Customer, error) {
	item := &Customer{}

	sqls := `update customers set active=$2 where id=$1 returning *`
	err := s.db.QueryRowContext(ctx, sqls, id, active).Scan(
	&item.ID,
	&item.Name,
	&item.Phone,
	&item.Active,
	&item.Created)

	if err == sql.ErrNoRows {
	return nil, ErrNotFound
	}
	if err != nil {
	log.Print(err)
	return nil, ErrInternal
	}
	return item, nil

}


func (s *Service) Delete(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	sqls := `delete from customers where id=$1 returning *`
	err := s.db.QueryRowContext(ctx, sqls, id).Scan(
	&item.ID,
	&item.Name,
	&item.Phone,
	&item.Active,
	&item.Created)

	if err == sql.ErrNoRows {
	return nil, ErrNotFound
	}
	if err != nil {
	log.Print(err)
	return nil, ErrInternal
	}
	return item, nil

}


func (s *Service) Save(ctx context.Context, customer *Customer) (c *Customer, err error) {

	item := &Customer{}

	if customer.ID == 0 {
	sqlStatement := `insert into customers(name, phone) values($1, $2) returning *`
	err = s.db.QueryRowContext(ctx, sqlStatement, customer.Name, customer.Phone).Scan(
	&item.ID,
	&item.Name,
	&item.Phone,
	&item.Active,
	&item.Created)
	} else {
	sqlStatement := `update customers set name=$1, phone=$2 where id=$3 returning *`
	err = s.db.QueryRowContext(ctx, sqlStatement, customer.Name, customer.Phone, customer.ID).Scan(
	&item.ID,
	&item.Name,
	&item.Phone,
	&item.Active,
	&item.Created)
	}

	if err != nil {
	log.Print(err)
	return nil, ErrInternal
	}
	return item, nil

}