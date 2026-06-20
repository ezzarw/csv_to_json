package main

import (
	// "bufio"
	// "errors"

	"encoding/csv"
	"encoding/json"
	"io"

	// "strings"

	// "encoding/json"
	"flag"
	"fmt"
	"os"
	// "strings"
)

func handle(err error) {
	if err != nil {
		fmt.Println("error")
		os.Exit(1)
	}
}

func main() {
	// baca flag dan input .csv
	file_input := flag.String("file", "", "masukkan sebuah file berformat .csv")
	flag.Parse()

	// eror kalo flag -file dikosongin
	if *file_input == "" {
		fmt.Println("eror: tidak boleh kosong")
		flag.Usage() // print help
		os.Exit(1)
	}

	// baca bytes dari file_input dan dijadikan agar read_csv bisa baca
	// read_file_data, err := os.ReadFile(*file_input)
	// handle(err)

	read_file_data, err := os.Open(*file_input)
	handle(err)

	read_csv := csv.NewReader(read_file_data) // read dari read_file_data buat read csv

	// mulai iterasi untuk looping csv ke json secara dynamic
	per_row_iteration := 0
	id_per_object := 0

	header := []string{}
	map_storage := make(map[int]interface{})

	for {
		row, err := read_csv.Read()

		if err == io.EOF {
			break
		}

		// header, deklarasi header agar sebagai key (column) dari setiap value (row)
		if per_row_iteration == 0 {
			header = row
		}

		if per_row_iteration > 0 {
			map_storage_sementara := make(map[string]interface{})
			for key_id := 0; key_id < len(header); key_id++ {
				map_storage_sementara[header[key_id]] = row[key_id]
				map_storage[id_per_object] = map_storage_sementara
			}
			id_per_object++
		}

		handle(err)
		per_row_iteration++
	}

	json_result, err := json.Marshal(map_storage)
	handle(err)

	fmt.Println(string(json_result))
}
