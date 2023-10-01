package main

import (
	"encoding/json"
	"os"
)

func loadFile(filename string) ([]Movie, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var movies []Movie
	err = json.NewDecoder(file).Decode(&movies)
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func saveFile(filename string, movies []Movie) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	err = json.NewEncoder(file).Encode(movies)
	if err != nil {
		return err
	}

	return nil
}

func get(url string) ([]Movie, error) {
	return loadFile(url)
}

func post(url string, movie Movie) ([]Movie, error) {
	records, err := loadFile(url)
	if err != nil {
		return nil, err
	}

	records = append(records, movie)

	err = saveFile(url, records)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func patch(url string, data Movie) (*Movie, error) {
	records, err := loadFile(url)
	if err != nil {
		return nil, err
	}

	for i, record := range records {
		if record.Id == data.Id {
			records[i] = data
			break
		}
	}

	err = saveFile(url, records)
	if err != nil {
		return nil, err
	}

	return &data, nil

}

func del(url string, id string) ([]Movie, error) {
	records, err := loadFile(url)
	if err != nil {
		return nil, err
	}

	for i, record := range records {
		if record.Id == id {
			records = append(records[:i], records[i+1:]...)
			break
		}
	}

	err = saveFile(url, records)
	if err != nil {
		return nil, err
	}

	return records, nil
}
