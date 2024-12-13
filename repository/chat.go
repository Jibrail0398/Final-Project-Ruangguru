package repository

import (
	// "a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/model"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)



type ChatRepository interface {
	SaveChat(date string,question string, response string, fk_id_user int, fk_report_id int)error
	GetChatByReport(id_report int)([]model.GetChat,error)
}

type chatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *chatRepository {
	return &chatRepository{db: db}
}

func(u *chatRepository) SaveChat(date string,question string, response string, fk_id_user int, fk_report_id int) error{

	query := `
		INSERT INTO chat(date,question,response,fk_id_user,fk_report_id)
		VALUES($1,$2,$3,$4,$5)
	`

	_,err := u.db.Exec(query,date,question, response, fk_id_user, fk_report_id)

	if err!= nil{
		return fmt.Errorf("error while save chat")
	}

	return nil
}

func(u *chatRepository) GetChatByReport(id_report int)([]model.GetChat,error){
	query := `
		SELECT id,date,question,response FROM chat WHERE fk_report_id = $1
	`

	row,err:=u.db.Query(query,id_report)

	if err!=nil{
		if err == sql.ErrNoRows{
			return nil, fmt.Errorf("no row found based on id:",id_report)
		}

		return nil, fmt.Errorf("error while querying chat data")
	}

	var chatData []model.GetChat

	for row.Next(){
		var data model.GetChat

		err:= row.Scan(&data.Id,&data.Date,&data.Question,&data.Response)

		if err!=nil{
			return nil,fmt.Errorf("Error while scanning data",err.Error())
		}

		chatData = append(chatData, data)

	}
	return chatData,nil

}