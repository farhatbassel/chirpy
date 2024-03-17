package database

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	Body string `json:"body"`
	Id   int    `json:"Id"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	err := db.ensureDB()

	if err != nil {
		if err := db.writeDB(DBStructure{}); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	chirps, err := db.GetChirps()

	if err != nil {
		return Chirp{}, err
	}

	ids := make([]int, len(chirps))

	for index, chirp := range chirps {
		ids[index] = chirp.Id
	}

	newChirp := Chirp{
		Id:   getNewID(ids),
		Body: body,
	}

	chirps = append(chirps, newChirp)

	chirpsMap := make(map[int]Chirp)
	for _, chirp := range chirps {
		chirpsMap[chirp.Id] = chirp
	}

	dbStructure := DBStructure{
		Chirps: chirpsMap,
	}

	if err := db.writeDB(dbStructure); err != nil {
		return Chirp{}, err
	}

	return newChirp, nil
}

func getNewID(ids []int) int {
	if len(ids) == 0 {
		return 1
	}

	maxID := ids[0]
	for _, id := range ids {
		if id > maxID {
			maxID = id
		}
	}

	return maxID + 1
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return []Chirp{}, err
	}

	var chirps []Chirp

	for _, value := range dbStructure.Chirps {
		chirps = append(chirps, value)
	}

	return chirps, nil
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)

	if errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	var dBStructure DBStructure

	file, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	if len(file) == 0 {
		return DBStructure{
			Chirps: make(map[int]Chirp),
		}, nil
	}

	if err := json.Unmarshal(file, &dBStructure); err != nil {
		return DBStructure{}, err
	}

	return dBStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	data, err := json.Marshal(dbStructure)

	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, data, fs.FileMode(0664))

	if err != nil {
		return err
	}

	return nil
}
