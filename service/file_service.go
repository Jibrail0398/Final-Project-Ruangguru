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

	// Jika konten file (fileContent) hanya berisi spasi (atau kosong sepenuhnya), maka fungsi akan mengembalikan error dengan pesan "file content is empty".
	if strings.TrimSpace(fileContent) == "" {
		return nil, errors.New("file content is empty")
	}

	//Membaca file csv
	reader := csv.NewReader(strings.NewReader(fileContent))
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, errors.New("failed to parse CSV")
	}

	//Jika csv tidak memiliki data, Karena pada csv baris pertama adalah judul
	if len(rows) < 2 {
		return nil, errors.New("invalid CSV format")
	}

	headers := rows[0] //Mengambil header
	data := make(map[string][]string) //Membuat map kosong
	for _, header := range headers { //// Inisialisasi slice kosong untuk setiap header
		data[header] = []string{}
	}

	// Tambahkan data dari baris berikutnya
	for _, row := range rows[1:] { /// Iterasi dari baris kedua hingga akhir
		if len(row) != len(headers) { //Jika jumlah kolom tidak cocok, itu menandakan file CSV rusak atau tidak valid.
			return nil, errors.New("row length does not match headers")
		}
		//Menambahkan elemen pada setiap key dari map berdasarkan header
		for index, value := range row {
			data[headers[index]] = append(data[headers[index]], value)
		}
	}

	return data, nil
}
