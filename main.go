package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/ogen-go/ogen/conv"
)

func handle(err error) {
	if os.Getenv("APP") == "local" {
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if os.Getenv("APP") == "production" {
		if err != nil {
			fmt.Println("error")
			os.Exit(1)
		}
	}
}

func main() {
	// baca flag dan input .csv
	file_input := flag.String("f", "", "input sebuah file berformat .csv, example:\"./csv_to_json -f ./file.csv\"")
	file_output := flag.String("o", "", "membuat file baru sebagai output, yang berformat .json example:\"./csv_to_json -f ./file.csv -o output.json\"")
	flag.Parse()

	// eror kalo flag -file dikosongin
	if *file_input == "" {
		fmt.Println("eror: tidak boleh kosong")
		flag.Usage() // print help
		os.Exit(1)
	}

	read_file_data, err := os.Open(*file_input)
	handle(err)

	read_csv := csv.NewReader(read_file_data) // read dari read_file_data buat read csv

	// mulai iterasi untuk looping csv ke json secara dynamic
	per_row_iteration := 0
	id_per_object := 0

	header := []string{}
	map_storage := make(map[string]interface{})

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
				if id_per_object < 10 {
					map_storage["0"+conv.IntToString(id_per_object)] = map_storage_sementara
				} else {
					map_storage[conv.IntToString(id_per_object)] = map_storage_sementara
				}
			}
			id_per_object++
		}

		handle(err)
		per_row_iteration++
	}

	json_result_byte, err := json.MarshalIndent(map_storage, "", "  ")
	handle(err)

	json_result := string(json_result_byte)

	if *file_output != "" {
		os.Create(string(*file_output))
		err := os.WriteFile(string(*file_output), json_result_byte, 0760)
		_ = os.Chmod(string(*file_output), 0760)
		handle(err)
	}

	fmt.Println(json_result)
}
