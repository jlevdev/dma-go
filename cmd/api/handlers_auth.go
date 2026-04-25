package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/idtoken"
)

type googleAuthRequest struct {
	Credential string `json:"credential"`
}

func (app *application) googleAuth(w http.ResponseWriter, r *http.Request) {
	var req googleAuthRequest
	if err := app.readJSON(r, &req); err != nil {
		app.errorJSON(w, err)
		return
	}

	payload, err := idtoken.Validate(r.Context(), req.Credential, app.GoogleClientID)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credential"), http.StatusUnauthorized)
		return
	}

	googleID := payload.Subject
	email, _ := payload.Claims["email"].(string)
	displayName, _ := payload.Claims["name"].(string)

	user, err := app.upsertUserByGoogleID(googleID, email, displayName)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	token, err := app.generateJWT(user.Id)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (app *application) upsertUserByGoogleID(googleID, email, displayName string) (*User, error) {
	var user User
	err := app.DB.QueryRow(`
		INSERT INTO Users (Username, GoogleId, Email, DisplayName)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (GoogleId) DO UPDATE
			SET Email = EXCLUDED.Email,
			    DisplayName = EXCLUDED.DisplayName,
			    UpdatedAt = NOW()
		RETURNING Id, CreatedAt, UpdatedAt, Username, GoogleId, Email, DisplayName
	`, email, googleID, email, displayName).Scan(
		&user.Id, &user.CreatedAt, &user.UpdatedAt,
		&user.Username, &user.GoogleId, &user.Email, &user.DisplayName,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (app *application) generateJWT(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(app.JWTSecret))
}
