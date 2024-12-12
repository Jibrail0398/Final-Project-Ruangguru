package repository

import (
	// "a21hc3NpZ25tZW50/model"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

type AnalyzeReportRepository interface {
	SaveReport(id_user int,date string)
	
}

type analyzeReportRepository struct {
	db *sql.DB
}

func NewAnalyzeReportRepository(db *sql.DB) *analyzeReportRepository {
	return &analyzeReportRepository{db: db}
}

func(u *analyzeReportRepository) SaveReport(id_user int,date time.Time)error{

	query := `
		INSERT INTO report(date,fk_id_user) VALUES($1,$2)
	`
	_,err := u.db.Exec(query,date,id_user)

	if err!=nil{
		return fmt.Errorf("error save report from id user: %d",id_user)
	}

	return nil

}
// func(u *analyzeReportRepository) GetReportByUser(id_user int,date time.Time) ([]model.Report,error){

// 	var reports []model.Report

// 	return reports,nil
// }

// func(u *analyzeReportRepository) CheckDate(date string)bool{

// 	return false
// }
