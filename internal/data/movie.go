package data

import (
	"time"

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
