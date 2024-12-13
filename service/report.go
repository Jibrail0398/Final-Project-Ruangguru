package service

import(
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/repository"
	
)

type ReportService interface{
	Upload(id_user int,date string,stringText string)error
	GetReportByUser(id_user int) ([]model.Report,error)
	Delete(id_user int)error
}

type reportService struct{
	reportRepository repository.ReportRepository
}

func NewReportService(reportRepository repository.ReportRepository) ReportService {
	return &reportService{reportRepository}
	
}

func(s *reportService) Upload(id_user int,date string,stringText string) error {

	err := s.reportRepository.SaveReport(id_user,date,stringText)
	
	if err!=nil{
		return err
	}

	return nil

}

func(s *reportService) GetReportByUser(id_user int) ([]model.Report,error){

	report,err := s.reportRepository.GetReportByUser(id_user)

	if err!=nil{
		return nil,err
	}

	return report,nil
}

func(s *reportService) Delete(id_user int) error{
	err := s.reportRepository.Delete(id_user)

	if err!=nil{
		return err
	}
	return nil
}