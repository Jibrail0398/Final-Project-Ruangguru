package db

import(
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

func(p*postgres) Connect(host string, port int, user string, password string, dbname string) (*sql.DB,error) {
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
        "password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

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
