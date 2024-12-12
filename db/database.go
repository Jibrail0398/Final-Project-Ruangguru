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

    _,err = p.DB.Exec(
        `CREATE TABLE IF NOT EXISTS report (
		id SERIAL PRIMARY KEY,
        date DATE NOT NULL UNIQUE,
        stringText TEXT NOT NULL,
		fk_id_user INTEGER NOT NULL,
		FOREIGN KEY (fk_id_user) REFERENCES users(id))
		`)
    if err!=nil{
        return fmt.Errorf("failed to create tabel report: %w", err)
    }

    _,err = p.DB.Exec(
        `CREATE TABLE IF NOT EXISTS analyze_report (
		id SERIAL PRIMARY KEY,
		date DATE NOT NULL,
        response TEXT NOT NULL,
        fk_id_user INTEGER NOT NULL,
        fk_report_id INTEGER NOT NULL,
        fk_query_id INTEGER NOT NULL,
        FOREIGN KEY (fk_id_user) REFERENCES users(id),
        FOREIGN KEY (fk_report_id) REFERENCES report(id),
        FOREIGN KEY (fk_query_id) REFERENCES queryAI(id))
		`)
    if err!=nil{
        return fmt.Errorf("failed to create tabel analyze_report: %w", err)
    }

    _,err = p.DB.Exec(
        `CREATE TABLE IF NOT EXISTS chat (
		id SERIAL PRIMARY KEY,
		date DATE NOT NULL,
        question TEXT NOT NULL,
        response TEXT NOT NULL,
        fk_id_user INTEGER NOT NULL,
        fk_report_id INTEGER NOT NULL,
        FOREIGN KEY (fk_id_user) REFERENCES users(id),
        FOREIGN KEY (fk_report_id) REFERENCES report(id)
        )
		`)
    if err!=nil{
        return fmt.Errorf("failed to create tabel chat: %w", err)
    }

    _,err = p.DB.Exec(
        `CREATE TABLE IF NOT EXISTS queryAI (
		id SERIAL PRIMARY KEY,
        query VARCHAR(255) NOT NULL UNIQUE)
		`)
    if err!=nil{
        return fmt.Errorf("failed to create tabel queryAI: %w", err)
    }


    return nil
}
