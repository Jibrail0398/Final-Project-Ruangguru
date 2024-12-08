package service

import (
	repository "a21hc3NpZ25tZW50/repository/fileRepository"
	"errors"
	"strings"
	"encoding/csv"
)

type FileService struct {
	Repo *repository.FileRepository
}


func (s *FileService) ProcessFile(fileContent string) (map[string][]string, error) {
	if strings.TrimSpace(fileContent) == "" {
		return nil, errors.New("file content is empty")
	}

	reader := csv.NewReader(strings.NewReader(fileContent))
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, errors.New("failed to parse CSV")
	}

	if len(rows) < 2 {
		return nil, errors.New("invalid CSV format")
	}

	headers := rows[0]
	data := make(map[string][]string)
	for _, header := range headers {
		data[header] = []string{}
	}

	for _, row := range rows[1:] {
		if len(row) != len(headers) {
			return nil, errors.New("row length does not match headers")
		}
		for i, value := range row {
			data[headers[i]] = append(data[headers[i]], value)
		}
	}

	return data, nil
}
