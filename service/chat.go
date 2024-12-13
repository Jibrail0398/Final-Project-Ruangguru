package service

import(
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/repository"
	
)

type ChatService interface{
	SaveChat(date string,question string, response string, fk_id_user int, fk_report_id int)error
	GetChatByReport(id_report int) ([]model.GetChat,error)
}

type chatService struct{
	chatRepository repository.ChatRepository
}

func NewChatService(chatRepository repository.ChatRepository) ChatService {
	return &chatService{chatRepository}
	
}

func(s *chatService) SaveChat(date string,question string, response string, fk_id_user int, fk_report_id int)error{

	err:= s.chatRepository.SaveChat(date,question, response, fk_id_user, fk_report_id)

	if err!= nil{
		return err
	}
	return nil
}

func (s *chatService)GetChatByReport(id_report int) ([]model.GetChat,error){

	data,err := s.chatRepository.GetChatByReport(id_report)

	if err!= nil{
		return nil,err
	}

	return data,nil
}
