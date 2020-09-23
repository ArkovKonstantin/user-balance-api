package provider

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"user-balance-api/models"

	// driver
	_ "github.com/lib/pq"
)

// Provider handlers
type Provider interface {
	Open() error
	GetConn() (*sql.DB, error)
}

type provider struct {
	db        *sql.DB
	cs        string
	idlConns  int
	openConns int
	lifetime  time.Duration
}

// New new provider
func New(db *models.SQLDataBase) Provider {
	info := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s  sslmode=disable",
		db.Server, db.Port, db.User, db.Password, db.Database)
	fmt.Println(info)
	return &provider{
		cs:        info,
		idlConns:  db.MaxIdleConns,
		openConns: db.MaxOpenConns,
		lifetime:  time.Duration(db.ConnMaxLifetime),
	}
}

// Open connection
func (p *provider) Open() error {
	var err error
	p.db, err = sql.Open("postgres", p.cs)

	if err != nil {
		log.Println(err)
		return err
	}

	p.db.SetMaxIdleConns(p.idlConns)
	p.db.SetMaxOpenConns(p.openConns)
	p.db.SetConnMaxLifetime(p.lifetime)

	err = p.db.Ping()

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// GetConn return sql database
func (p *provider) GetConn() (*sql.DB, error) {
	return p.db, nil
}
