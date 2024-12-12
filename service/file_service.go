package service

import (
	repository "a21hc3NpZ25tZW50/repository/fileRepository"
	"errors"
	"strings"
	"encoding/csv"
	"strconv"
	"fmt"
	"time"
)

type FileService struct {
	Repo *repository.FileRepository
}

func validateDateFormat(dates []string) error {
    if len(dates) == 0 {
        return errors.New("no dates to validate")
    }

    // Ambil bulan dari tanggal pertama sebagai referensi
    firstDateParts := strings.Split(dates[0], "-")
    referenceMonth := firstDateParts[1]

    // Validasi setiap tanggal
    for _, dateStr := range dates {
        // Validasi format tanggal
        _, err := time.Parse("2006-01-02", dateStr)
        if err != nil {
            return fmt.Errorf("invalid date format: %v", err)
        }

        // Pisahkan tahun, bulan, dan hari
        parts := strings.Split(dateStr, "-")
        if len(parts) != 3 {
            return errors.New("invalid date format")
        }

        // Konversi bulan ke integer untuk validasi rentang
        month, err := strconv.Atoi(parts[1])
        if err != nil {
            return errors.New("invalid month format")
        }

        // Validasi rentang bulan (1-12)
        if month < 1 || month > 12 {
            return errors.New("month must be between 1 and 12")
        }

        // Periksa konsistensi bulan
        if parts[1] != referenceMonth {
            return fmt.Errorf("inconsistent month: expected %s, found %s", referenceMonth, parts[1])
        }
    }

    return nil
}

func (s *FileService) ProcessFile(fileContent string) (map[string][]string , error) {

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

    // Temukan indeks kolom Date
    dateColIndex := -1
    for i, header := range headers {
        if header == "Date" {
            dateColIndex = i
            break
        }
    }

    // Jika kolom Date tidak ditemukan, kembalikan error
    if dateColIndex == -1 {
        return nil, errors.New("date column not found")
    }

    // Kumpulkan semua tanggal untuk validasi
    var dates []string
    for _, row := range rows[1:] {
        dates = append(dates, row[dateColIndex])
    }

    // Validasi tanggal
    err = validateDateFormat(dates)
    if err != nil {
        return nil, fmt.Errorf("date validation error: %v", err)
    }

    return data, nil
}
