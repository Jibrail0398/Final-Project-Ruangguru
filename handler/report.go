package handler

import (
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)


type ReportHandler struct {
	ReportService service.ReportService
}

func NewReportHandler(reportService service.ReportService) *ReportHandler {
	return &ReportHandler{
		ReportService: reportService,
	}
}

func (h * ReportHandler) Upload(w http.ResponseWriter, r *http.Request){

	var fileService = &service.FileService{}

	err := r.ParseMultipartForm(1024)
	if err != nil {
		w.WriteHeader(500)
		errorResponse := model.Error{Error: err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(400)
		errorResponse := model.Error{Error: err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(500)
		errorResponse := model.Error{Error: err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	resultFile, err := fileService.ProcessFile(string(content))
	if err != nil {
		w.WriteHeader(500)
		errorResponse := model.Error{Error: err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	resultJson,_:=json.Marshal(resultFile)
	resultString := string(resultJson)

	id_user := r.FormValue("id_user")
	id_userInt,_ := strconv.Atoi(id_user)
	dateFile := string(resultFile["Date"][0])

	err = h.ReportService.Upload(id_userInt,dateFile,resultString)
	if err!=nil{
		w.WriteHeader(400)
		errorResponse := model.Error{Error: err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	response := map[string]string{
		"Message":"success upload file",
		"id_report":id_user,
		"stringText":resultString,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	
}

func(h *ReportHandler) GetReportByUser(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id_user := vars["id"] 
	id_userInt,_ := strconv.Atoi(id_user)

	report,err := h.ReportService.GetReportByUser(id_userInt)
	log.Println("terjadi error di reportService")
	if err!=nil{
		w.WriteHeader(400)
		errorResponse := model.Error{Error: err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)

}

func(h *ReportHandler) Delete(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id_user := vars["id"] 
	id_userInt,_ := strconv.Atoi(id_user)

	err := h.ReportService.Delete(id_userInt)

	if err!=nil{
		w.WriteHeader(500)
		errorResponse := model.Error{Error: err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	response := map[string]string{
		"Message" : "Delete report succed",
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}