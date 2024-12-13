package repository

import (
	"a21hc3NpZ25tZW50/model"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	
)

type ReportRepository interface {
	SaveReport(id_user int,date string,stringText string) error
	GetReportByUser(id_user int) ([]model.Report,error)
	Delete(id_user int)error
}

type reportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *reportRepository {
	return &reportRepository{db: db}
}

func(u *reportRepository) SaveReport(id_user int,date string,stringText string)error{

	query := `
		INSERT INTO report(date,stringText,fk_id_user) VALUES($1,$2,$3)
	`
	_,err := u.db.Exec(query,date,stringText,id_user)

	if err!=nil{
		return fmt.Errorf("error save report from id user: %d",id_user)
	}

	return nil

}
func(u *reportRepository) GetReportByUser(id_user int) ([]model.Report,error){

	var reports []model.Report

	query:=`
		SELECT * FROM report WHERE id = $1 
	`

	rows,err:=u.db.Query(query,id_user)
	if err!=nil{
		if err == sql.ErrNoRows{
			return nil,fmt.Errorf("report with id_user %d not found",id_user)
		}
		return nil, fmt.Errorf("error while querying report")
	}

	defer rows.Close()
	for rows.Next(){
		var report model.Report

		err:=rows.Scan(&report.Id,&report.Date,&report.StringText,&report.FK_ID_USER)
		if err!=nil{
			return nil,fmt.Errorf("error when scan report",err.Error())
		}
		reports = append(reports,report)
	}

	return reports,nil
}

func(u *reportRepository) Delete(id_user int)error{
	query := `
		DELETE FROM report WHERE id = $1;
	`
	_,err:=u.db.Exec(query,id_user)

	if err!=nil{
		return fmt.Errorf(err.Error())
	}

	return nil
}

