package main

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) ListUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := app.DB.QueryContext(r.Context(),
		`SELECT Id, CreatedAt, UpdatedAt, Username FROM Users ORDER BY Id`)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Id, &u.CreatedAt, &u.UpdatedAt, &u.Username); err != nil {
			app.errorJSON(w, err, http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	_ = app.writeJSON(w, http.StatusOK, users)
}

func (app *application) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.errorJSON(w, errors.New("invalid id"))
		return
	}

	var u User
	err = app.DB.QueryRowContext(r.Context(),
		`SELECT Id, CreatedAt, UpdatedAt, Username FROM Users WHERE Id = $1`, id).
		Scan(&u.Id, &u.CreatedAt, &u.UpdatedAt, &u.Username)
	if errors.Is(err, sql.ErrNoRows) {
		app.notFound(w)
		return
	}
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, u)
}

func (app *application) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
	}
	if err := app.readJSON(r, &input); err != nil {
		app.errorJSON(w, err)
		return
	}

	var u User
	err := app.DB.QueryRowContext(r.Context(),
		`INSERT INTO Users (Username) VALUES ($1)
		 RETURNING Id, CreatedAt, UpdatedAt, Username`,
		input.Username).
		Scan(&u.Id, &u.CreatedAt, &u.UpdatedAt, &u.Username)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	_ = app.writeJSON(w, http.StatusCreated, u)
}

func (app *application) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.errorJSON(w, errors.New("invalid id"))
		return
	}

	var input struct {
		Username string `json:"username"`
	}
	if err := app.readJSON(r, &input); err != nil {
		app.errorJSON(w, err)
		return
	}

	var u User
	err = app.DB.QueryRowContext(r.Context(),
		`UPDATE Users SET Username = $1, UpdatedAt = NOW()
		 WHERE Id = $2
		 RETURNING Id, CreatedAt, UpdatedAt, Username`,
		input.Username, id).
		Scan(&u.Id, &u.CreatedAt, &u.UpdatedAt, &u.Username)
	if errors.Is(err, sql.ErrNoRows) {
		app.notFound(w)
		return
	}
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, u)
}

func (app *application) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.errorJSON(w, errors.New("invalid id"))
		return
	}

	result, err := app.DB.ExecContext(r.Context(),
		`DELETE FROM Users WHERE Id = $1`, id)
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
