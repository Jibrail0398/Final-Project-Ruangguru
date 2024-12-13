package handler

import (
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/service"
	"net/http"
	"strconv"
	"encoding/json"
	"github.com/gorilla/mux"
)


type ChatCRUDHandler struct {
	ChatService service.ChatService
}

func NewChatCRUDHandler(chatService service.ChatService) *ChatCRUDHandler {
	return &ChatCRUDHandler{
		ChatService: chatService,
	}
}


func(h* ChatCRUDHandler) SaveChat(w http.ResponseWriter, r *http.Request){

	date := r.FormValue("date")
	question := r.FormValue("question")
	response := r.FormValue("response")
	fk_id_user := r.FormValue("fk_id_user")
	fk_report_id := r.FormValue("fk_report_id")

	fk_id_user_int,_ := strconv.Atoi(fk_id_user)
	fk_report_id_int,_ := strconv.Atoi(fk_report_id)

	err:= h.ChatService.SaveChat(date ,question , response , fk_id_user_int , fk_report_id_int)

	if err!=nil{
		w.WriteHeader(400)
		response := model.Error{Error: "Failed to save the chat"}
		json.NewEncoder(w).Encode(response)
		return
	}

	responseSuccess := map[string]string{
		"Messages":"Save chat succeed",
		"date" :date, 
		"question":question,
		"response" : response,
		"fk_id_user":fk_id_user,
		"fk_report_id":fk_report_id,

	}
	w.WriteHeader(201)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseSuccess)
}

func(h *ChatCRUDHandler) GetChatByReport(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id_user := vars["id"]

	id_user_Int,_ := strconv.Atoi(id_user)
	data,err := h.ChatService.GetChatByReport(id_user_Int)

	if err!=nil{
		w.WriteHeader(500)
		response:= "Error while get chat data "+err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}