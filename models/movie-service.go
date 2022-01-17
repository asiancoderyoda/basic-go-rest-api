package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

func (db *DBModel) GetMovieById(id int) (*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	queryMovie := `
	SELECT id, title, description, year, release_date, runtime, rating, mpaa_rating, created_at, updated_at
	FROM public.movie_entity WHERE id = $1;
	`
	row := db.DB.QueryRowContext(ctx, queryMovie, id)

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
		if err == sql.ErrNoRows {
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
	rows, _ := db.DB.QueryContext(ctx, queryGenre, id)
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

func (db *DBModel) GetAllMovie(genre_ids ...int) ([]*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	where_clause := ""
	if len(genre_ids) > 0 {
		where_clause = fmt.Sprintf("WHERE id in (SELECT movie_id FROM movies_genres WHERE genre_id = %d)", genre_ids[0])
	}

	queryMovie := fmt.Sprintf(`
	SELECT id, title, description, year, release_date, runtime, rating, mpaa_rating, created_at, updated_at
	FROM public.movie_entity %s ORDER BY title;
	`, where_clause)
	queryRows, err := db.DB.QueryContext(ctx, queryMovie)

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

		genrRows, err := db.DB.QueryContext(ctx, queryGenre, movie.ID)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

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

func (db *DBModel) GetAllGenres() ([]*Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, genre_name, created_at, updated_at
		FROM public.genres;
	`
	rows, _ := db.DB.QueryContext(ctx, query)

	genres := make([]*Genre, 0)
	for rows.Next() {
		var genre Genre
		err := rows.Scan(
			&genre.ID,
			&genre.GenreName,
			&genre.CreatedAt,
			&genre.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		genres = append(genres, &genre)
	}
	return genres, nil
}

func (db *DBModel) DeleteMovie(id int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	getMovieQuery := `
		SELECT id, title FROM public.movie_entity WHERE id = $1;
	`

	row := tx.QueryRowContext(ctx, getMovieQuery, id)

	var movie Movie

	err = row.Scan(
		&movie.ID,
		&movie.Title,
	)

	if err != nil {
		return false, err
	}

	queryDeleteMovie := `
		DELETE FROM public.movie_entity WHERE id = $1;
	`
	_, errorDeleteMovie := db.DB.ExecContext(ctx, queryDeleteMovie, id)

	if errorDeleteMovie != nil {
		return false, err
	}

	queryDeleteMovieGenre := `
		DELETE FROM public.movies_genres WHERE movie_id = $1;
	`

	_, errorDeleteMovieGenre := db.DB.ExecContext(ctx, queryDeleteMovieGenre, id)

	if errorDeleteMovieGenre != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}
