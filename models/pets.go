package models

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type Store struct {
	DB *sql.DB
}

type Pet struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Kind string `json:"kind"`
	Age  int    `json:"age"`
}

func (s *Store) InsertPet(age int, name, kind string) (int, error) {
	var ID int
	query := sq.Insert("pets").
		Columns("pet_name", "pet_kind", "pet_age").
		Values(name, kind, age).
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
		err := rows.Scan(&pet.ID, &pet.Name, &pet.Kind, &pet.Age)
		if err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	}
	return pets, nil
}

func (s *Store) RemoveAllPets() error {
	_, err := sq.Delete("").From("pets").RunWith(s.DB).Query()
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) FilterPets(kind string, firstLetter string, underage, overage int) ([]*Pet, error) {
	query := sq.Select("*").
		From("pets").
		Where(sq.Lt{"pet_age": underage}).
		Where(sq.Gt{"pet_age": overage}).
		Where(sq.Eq{"pet_kind": kind}).
		Where(sq.Like{"pet_name": firstLetter + "%"}).
		PlaceholderFormat(sq.Dollar).
		RunWith(s.DB)
	rows, err := query.Query()
	if err != nil {
		return nil, err
	}

	pets := make([]*Pet, 0)
	for rows.Next() {
		pet := &Pet{}
		err := rows.Scan(&pet.ID, &pet.Name, &pet.Kind, &pet.Age)
		if err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	}
	return pets, nil
}
