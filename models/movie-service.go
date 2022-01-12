package models

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
)

type DBModel struct {
	DB *pgx.Conn
}

func (db *DBModel) GetMovieById(id int) (*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT * FROM movie_entity WHERE id = $1`
	row := db.DB.QueryRow(ctx, query, id)

	var movie Movie

	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.Year,
		&movie.ReleaseDate,
		&movie.RunTime,
		&movie.Rating,
		&movie.MPAARating,
		&movie.CreatedAt,
		&movie.UpdatedAt,
		&movie.MovieGenres,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			fmt.Println("No movie found with id: ", id)
			return nil, nil
		}
		return nil, err
	}

	return &movie, nil
}

func (db *DBModel) GetAllMovie() ([]*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT * FROM movie_entity`
	rows, err := db.DB.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	// var movies []*Movie
	movies := make([]*Movie, 0) // This is the same as the above line but JSON marshalling will consider this as empty slice as its not a nil pointer unlike the above line.

	for rows.Next() {
		var movie Movie

		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Description,
			&movie.Year,
			&movie.ReleaseDate,
			&movie.RunTime,
			&movie.Rating,
			&movie.MPAARating,
			&movie.CreatedAt,
			&movie.UpdatedAt,
			&movie.MovieGenres,
		)

		if err != nil {
			return nil, err
		}

		movies = append(movies, &movie)
	}

	return movies, nil
}
