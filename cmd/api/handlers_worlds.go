package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) ListWorlds(w http.ResponseWriter, r *http.Request) {
	rows, err := app.DB.QueryContext(r.Context(),
		`SELECT Id, CreatedAt, UpdatedAt, Title, Thumbnail, MapData FROM World ORDER BY Id`)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var worlds []World
	for rows.Next() {
		var wr World
		var mapData []byte
		if err := rows.Scan(&wr.Id, &wr.CreatedAt, &wr.UpdatedAt, &wr.Title, &wr.Thumbnail, &mapData); err != nil {
			app.errorJSON(w, err, http.StatusInternalServerError)
			return
		}
		wr.MapData = json.RawMessage(mapData)
		worlds = append(worlds, wr)
	}

	_ = app.writeJSON(w, http.StatusOK, worlds)
}

func (app *application) GetWorld(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.errorJSON(w, errors.New("invalid id"))
		return
	}

	var wr World
	var mapData []byte
	err = app.DB.QueryRowContext(r.Context(),
		`SELECT Id, CreatedAt, UpdatedAt, Title, Thumbnail, MapData FROM World WHERE Id = $1`, id).
		Scan(&wr.Id, &wr.CreatedAt, &wr.UpdatedAt, &wr.Title, &wr.Thumbnail, &mapData)
	if errors.Is(err, sql.ErrNoRows) {
		app.notFound(w)
		return
	}
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	wr.MapData = json.RawMessage(mapData)

	_ = app.writeJSON(w, http.StatusOK, wr)
}

func (app *application) CreateWorld(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title     string          `json:"title"`
		Thumbnail *string         `json:"thumbnail"`
		MapData   json.RawMessage `json:"map_data"`
	}
	if err := app.readJSON(r, &input); err != nil {
		app.errorJSON(w, err)
		return
	}

	if input.MapData == nil {
		input.MapData = json.RawMessage("{}")
	}

	var wr World
	var mapData []byte
	err := app.DB.QueryRowContext(r.Context(),
		`INSERT INTO World (Title, Thumbnail, MapData) VALUES ($1, $2, $3)
		 RETURNING Id, CreatedAt, UpdatedAt, Title, Thumbnail, MapData`,
		input.Title, input.Thumbnail, []byte(input.MapData)).
		Scan(&wr.Id, &wr.CreatedAt, &wr.UpdatedAt, &wr.Title, &wr.Thumbnail, &mapData)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	wr.MapData = json.RawMessage(mapData)

	_ = app.writeJSON(w, http.StatusCreated, wr)
}

func (app *application) UpdateWorld(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.errorJSON(w, errors.New("invalid id"))
		return
	}

	var input struct {
		Title     string          `json:"title"`
		Thumbnail *string         `json:"thumbnail"`
		MapData   json.RawMessage `json:"map_data"`
	}
	if err := app.readJSON(r, &input); err != nil {
		app.errorJSON(w, err)
		return
	}

	if input.MapData == nil {
		input.MapData = json.RawMessage("{}")
	}

	var wr World
	var mapData []byte
	err = app.DB.QueryRowContext(r.Context(),
		`UPDATE World SET Title = $1, Thumbnail = $2, MapData = $3, UpdatedAt = NOW()
		 WHERE Id = $4
		 RETURNING Id, CreatedAt, UpdatedAt, Title, Thumbnail, MapData`,
		input.Title, input.Thumbnail, []byte(input.MapData), id).
		Scan(&wr.Id, &wr.CreatedAt, &wr.UpdatedAt, &wr.Title, &wr.Thumbnail, &mapData)
	if errors.Is(err, sql.ErrNoRows) {
		app.notFound(w)
		return
	}
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	wr.MapData = json.RawMessage(mapData)

	_ = app.writeJSON(w, http.StatusOK, wr)
}

func (app *application) DeleteWorld(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.errorJSON(w, errors.New("invalid id"))
		return
	}

	result, err := app.DB.ExecContext(r.Context(),
		`DELETE FROM World WHERE Id = $1`, id)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	n, _ := result.RowsAffected()
	if n == 0 {
		app.notFound(w)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}
