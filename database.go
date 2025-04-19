package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Inicializa la base de datos SQLite
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./series.db")
	if err != nil {
		log.Fatal(err)
	}
	createTable()
}

// Cierra la base de datos SQLite
func closeDB() {
	db.Close()
}

// createTable crea la tabla de series si no existe
func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS series (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		status TEXT NOT NULL,
		lastEpisodeWatched INTEGER NOT NULL,
		totalEpisodes INTEGER NOT NULL,
		ranking INTEGER NOT NULL
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

// InsertSeries inserta una nueva serie en la base de datos
func InsertSeries(series Series) error {
	stmt, err := db.Prepare(`INSERT INTO series (title, status, lastEpisodeWatched, totalEpisodes, ranking) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(series.Title, series.Status, series.LastEpisodeWatched, series.TotalEpisodes, series.Ranking)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %v", err)
	}

	return nil
}

// GetAllSeries obtiene todas las series de la base de datos
func GetSeriesWithFilters(search, status, sort string) ([]Series, error) {
	query := "SELECT id, title, status, lastEpisodeWatched, totalEpisodes, ranking FROM series WHERE 1=1"
	var args []interface{}

	if search != "" {
		query += " AND title LIKE ?"
		args = append(args, "%"+search+"%")
	}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	if sort == "asc" {
		query += " ORDER BY ranking ASC"
	} else if sort == "desc" {
		query += " ORDER BY ranking DESC"
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query series: %v", err)
	}
	defer rows.Close()

	var seriesList []Series
	for rows.Next() {
		var series Series
		if err := rows.Scan(&series.ID, &series.Title, &series.Status, &series.LastEpisodeWatched, &series.TotalEpisodes, &series.Ranking); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		seriesList = append(seriesList, series)
	}

	if seriesList == nil {
		seriesList = []Series{}
	}

	return seriesList, nil
}

// GetSeriesByID obtiene una serie por su ID
func GetSeriesByID(id int) (Series, error) {
	row := db.QueryRow("SELECT id, title, status, lastEpisodeWatched, totalEpisodes, ranking FROM series WHERE id = ?", id)

	var series Series
	err := row.Scan(&series.ID, &series.Title, &series.Status, &series.LastEpisodeWatched, &series.TotalEpisodes, &series.Ranking)
	if err != nil {
		if err == sql.ErrNoRows {
			return Series{}, fmt.Errorf("no series found with ID %d", id)
		}
		return Series{}, fmt.Errorf("failed to scan row: %v", err)
	}

	return series, nil
}

// UpdateSeries actualiza una serie existente en la base de datos
func UpdateSeries(series Series) error {
	stmt, err := db.Prepare(`UPDATE series SET title = ?, status = ?, lastEpisodeWatched = ?, totalEpisodes = ?, ranking = ? WHERE id = ?`)

	if err != nil {
		return fmt.Errorf("fallo al preparar la declaración: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(series.Title, series.Status, series.LastEpisodeWatched, series.TotalEpisodes, series.Ranking, series.ID)
	if err != nil {
		return fmt.Errorf("fallo al ejecutar la declaración: %v", err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("fallo al obtener filas modificadas: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no se encontró ninguna serie con ID %d", series.ID)
	}
	return nil
}

// DeleteSeries elimina una serie de la base de datos
func DeleteSeries(id int) error {
	stmt, err := db.Prepare("DELETE FROM series WHERE id = ?")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %v", err)
	}

	return nil
}

// UpdateStatus actualiza el estado de una serie
func UpdateStatus(id int, status string) error {
	stmt, err := db.Prepare("UPDATE series SET status = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(status, id)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %v", err)
	}

	return nil
}

// UpdateEpisode actualiza el último episodio visto de una serie
func UpdateEpisode(id int, episode int) error {
	stmt, err := db.Prepare("UPDATE series SET lastEpisodeWatched = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(episode, id)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %v", err)
	}

	return nil
}

// UpVote sube en el ranking a una serie
func UpVote(id int) error {
	stmt, err := db.Prepare("UPDATE series SET ranking = ranking - 1 WHERE id = ?")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %v", err)
	}

	return nil
}

// DownVote baja en el ranking a una serie
func DownVote(id int) error {
	stmt, err := db.Prepare("UPDATE series SET ranking = ranking + 1 WHERE id = ?")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %v", err)
	}

	return nil
}
