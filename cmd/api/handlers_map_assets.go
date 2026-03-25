package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) ListMapAssets(w http.ResponseWriter, r *http.Request) {
	rows, err := app.DB.QueryContext(r.Context(),
		`SELECT Id, CreatedAt, UpdatedAt, WorldId, Title, Type, Data FROM MapAsset ORDER BY Id`)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var assets []MapAsset
	for rows.Next() {
		var a MapAsset
		var data []byte
		if err := rows.Scan(&a.Id, &a.CreatedAt, &a.UpdatedAt, &a.WorldId, &a.Title, &a.Type, &data); err != nil {
			app.errorJSON(w, err, http.StatusInternalServerError)
			return
		}
		a.Data = json.RawMessage(data)
		assets = append(assets, a)
	}

	_ = app.writeJSON(w, http.StatusOK, assets)
}

func (app *application) GetMapAsset(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.errorJSON(w, errors.New("invalid id"))
		return
	}

	var a MapAsset
	var data []byte
	err = app.DB.QueryRowContext(r.Context(),
		`SELECT Id, CreatedAt, UpdatedAt, WorldId, Title, Type, Data FROM MapAsset WHERE Id = $1`, id).
		Scan(&a.Id, &a.CreatedAt, &a.UpdatedAt, &a.WorldId, &a.Title, &a.Type, &data)
	if errors.Is(err, sql.ErrNoRows) {
		app.notFound(w)
		return
	}
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	a.Data = json.RawMessage(data)

	_ = app.writeJSON(w, http.StatusOK, a)
}

func (app *application) CreateMapAsset(w http.ResponseWriter, r *http.Request) {
	var input struct {
		WorldId *int64          `json:"world_id"`
		Title   string          `json:"title"`
		Type    string          `json:"type"`
		Data    json.RawMessage `json:"data"`
	}
	if err := app.readJSON(r, &input); err != nil {
		app.errorJSON(w, err)
		return
	}

	if input.Data == nil {
		input.Data = json.RawMessage("{}")
	}

	var a MapAsset
	var data []byte
	err := app.DB.QueryRowContext(r.Context(),
		`INSERT INTO MapAsset (WorldId, Title, Type, Data) VALUES ($1, $2, $3, $4)
		 RETURNING Id, CreatedAt, UpdatedAt, WorldId, Title, Type, Data`,
		input.WorldId, input.Title, input.Type, []byte(input.Data)).
		Scan(&a.Id, &a.CreatedAt, &a.UpdatedAt, &a.WorldId, &a.Title, &a.Type, &data)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	a.Data = json.RawMessage(data)

	_ = app.writeJSON(w, http.StatusCreated, a)
}

func (app *application) UpdateMapAsset(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.errorJSON(w, errors.New("invalid id"))
		return
	}

	var input struct {
		WorldId *int64          `json:"world_id"`
		Title   string          `json:"title"`
		Type    string          `json:"type"`
		Data    json.RawMessage `json:"data"`
	}
	if err := app.readJSON(r, &input); err != nil {
		app.errorJSON(w, err)
		return
	}

	if input.Data == nil {
		input.Data = json.RawMessage("{}")
	}

	var a MapAsset
	var data []byte
	err = app.DB.QueryRowContext(r.Context(),
		`UPDATE MapAsset SET WorldId = $1, Title = $2, Type = $3, Data = $4, UpdatedAt = NOW()
		 WHERE Id = $5
		 RETURNING Id, CreatedAt, UpdatedAt, WorldId, Title, Type, Data`,
		input.WorldId, input.Title, input.Type, []byte(input.Data), id).
		Scan(&a.Id, &a.CreatedAt, &a.UpdatedAt, &a.WorldId, &a.Title, &a.Type, &data)
	if errors.Is(err, sql.ErrNoRows) {
		app.notFound(w)
		return
	}
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	a.Data = json.RawMessage(data)

	_ = app.writeJSON(w, http.StatusOK, a)
}

func (app *application) DeleteMapAsset(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.errorJSON(w, errors.New("invalid id"))
		return
	}

	result, err := app.DB.ExecContext(r.Context(),
		`DELETE FROM MapAsset WHERE Id = $1`, id)
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
