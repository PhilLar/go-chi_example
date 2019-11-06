package models

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
)

type Store struct {
	DB *sql.DB
}

type Pet struct {
	ID		int    `json:"id"`
	Name 	string `json:"name"`
	Kind   	string `json:"kind"`
}


func (s *Store) InsertPet(name, kind string) (int, error) {
	var ID int
	query := sq.Insert("pets").
		Columns("pet_name", "pet_kind").
		Values(name, kind).
		Suffix("RETURNING \"id\"").
		RunWith(s.DB).
		PlaceholderFormat(sq.Dollar)

	err := query.QueryRow().Scan(&ID)
	if err != nil {
		return 0, err
	}
	return ID, nil
}

func (s *Store) ListPets() ([]*Pet, error) {
	query := sq.Select("*").From("pets").RunWith(s.DB)
	rows, err := query.Query()
	if err != nil {
		return nil, err
	}

	pets := make([]*Pet, 0)
	for rows.Next() {
		pet := &Pet{}
		err := rows.Scan(&pet.ID, &pet.Name, &pet.Kind)
		if err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	}
	return pets, nil
}

