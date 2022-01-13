package models

import (
	"time"

	"github.com/jackc/pgx/v4"
)

// Models is a wrapper for all the models in the database
type Models struct {
	DB DBModel
}

// NewModels returns models with db pool
func NewModels(db *pgx.Conn) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

type Movie struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Year        int       `json:"year"`
	ReleaseDate time.Time `json:"release_date"`
	RunTime     int       `json:"run_time"`
	Rating      int       `json:"rating"`
	MPAARating  string    `json:"mpaa_rating"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	// MovieGenres []MovieGenre `json:"genres"`
	MovieGenres map[int]string `json:"genres"`
}

type Genre struct {
	ID        int       `json:"id"`
	GenreName string    `json:"genre_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MovieGenre struct {
	ID        int       `json:"id"`
	MovieID   int       `json:"movie_id"`
	GenreID   int       `json:"genre_id"`
	Genre     Genre     `json:"genre"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
