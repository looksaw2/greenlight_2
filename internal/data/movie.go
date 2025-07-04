package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/looksaw/greenlight_2/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty,string"`
	Generes   []string  `json:"generes,omitempty"`
	Version   int32     `json:"version"`
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must be no longer than 500 words")

	v.Check(movie.Year != 0, "year", "year shouldn't be zero")
	v.Check(movie.Year >= 1888, "year", "year must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "year must be less than now")

	v.Check(movie.Runtime != 0, "runtime", "runtime must be provided")
	v.Check(movie.Runtime > 0, "runtime", "runtime must be positive")

}

type MovieModel struct {
	DB *sql.DB
}

func (m MovieModel) Insert(movie *Movie) error {
	query := `INSERT INTO movies (title,year,runtime,genres)
			  VALUES ($1,$2,$3,$4)
			  RETURNING id , create_at , version
			 `
	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Generes)}
	return m.DB.
		QueryRow(query, args...).
		Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	query := `
	SELECT id, create_at, title,year,runtime,genres,version
	FROM movies
	WHERE id = $1
	`
	var movie Movie

	err := m.DB.QueryRow(query, id).
		Scan(&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(movie.Generes),
			&movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &movie, nil
}

func (m MovieModel) Update(movie *Movie) error {
	return nil
}

func (m MovieModel) Delete(id int64) error {
	return nil
}

type MockMovieModel struct{}

func (m MockMovieModel) Insert(movie *Movie) error {
	return nil
}

func (m MockMovieModel) Get(id int64) (*Movie, error) {
	return nil, nil
}

func (m MockMovieModel) Update(movie *Movie) error {
	return nil
}

func (m MockMovieModel) Delete(id int64) error {
	return nil
}
