package db

import (
	"a21hc3NpZ25tZW50/model"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type postgres struct{
    DB *sql.DB
}

func NewDatabase() *postgres {
    return &postgres{}
}

func(p*postgres) Connect(credential *model.Credential) (*sql.DB,error) {
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
        "password=%s dbname=%s sslmode=disable",
        credential.Host, credential.Port, credential.Username, credential.Password, credential.DatabaseName)

    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        return nil, err
    }
    

    p.DB = db
    return db, nil
}

func(p*postgres)Migrate()error{
    if p.DB == nil{
        return fmt.Errorf("database connection is not initialized")
    }

    _,err := p.DB.Exec(
        `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) NOT NULL,
		email VARCHAR(100) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL)
		`)
    if err!=nil{
        return fmt.Errorf("failed to create tabel users: %w", err)
    }
    return nil
}
