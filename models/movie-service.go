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

	queryMovie := `
	SELECT id, title, description, year, release_date, runtime, rating, mpaa_rating, created_at, updated_at
	FROM public.movie_entity WHERE id = $1;
	`
	row := db.DB.QueryRow(ctx, queryMovie, id)

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
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			fmt.Println("No movie found with id: ", id)
			return nil, nil
		}
		return nil, err
	}

	queryGenre := `
	SELECT mg.id, mg.movie_id, mg.genre_id, g.genre_name
	FROM public.movies_genres mg
	LEFT JOIN genres g on (mg.genre_id = g.id)
	WHERE mg.movie_id = $1;
	`
	rows, _ := db.DB.Query(ctx, queryGenre, id)
	defer rows.Close()

	genres := make(map[int]string)
	for rows.Next() {
		var mg MovieGenre
		err := rows.Scan(
			&mg.ID,
			&mg.MovieID,
			&mg.GenreID,
			&mg.Genre.GenreName,
		)

		if err != nil {
			return nil, err
		}

		genres[mg.ID] = mg.Genre.GenreName
	}

	movie.MovieGenres = genres

	return &movie, nil
}

func (db *DBModel) GetAllMovie() ([]*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	queryMovie := `
	SELECT id, title, description, year, release_date, runtime, rating, mpaa_rating, created_at, updated_at
	FROM public.movie_entity ORDER BY title;
	`
	queryRows, err := db.DB.Query(ctx, queryMovie)

	if err != nil {
		return nil, err
	}

	// var movies []*Movie
	movies := make([]*Movie, 0) // This is the same as the above line but JSON marshalling will consider this as empty slice as its not a nil pointer unlike the above line.

	for queryRows.Next() {
		var movie Movie

		err := queryRows.Scan(
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
		)

		if err != nil {
			return nil, err
		}

		queryGenre := `
		SELECT mg.id, mg.movie_id, mg.genre_id, g.genre_name
		FROM public.movies_genres mg
		LEFT JOIN genres g on (mg.genre_id = g.id)
		WHERE mg.movie_id = $1;
		`

		genrRows, _ := db.DB.Query(ctx, queryGenre, movie.ID)

		genres := make(map[int]string)
		for genrRows.Next() {
			var mg MovieGenre
			err := genrRows.Scan(
				&mg.ID,
				&mg.MovieID,
				&mg.GenreID,
				&mg.Genre.GenreName,
			)
			if err != nil {
				return nil, err
			}
			genres[mg.ID] = mg.Genre.GenreName
		}
		genrRows.Close()
		movie.MovieGenres = genres
		movies = append(movies, &movie)
	}

	return movies, nil
}
