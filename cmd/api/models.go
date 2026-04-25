package main

import (
	"encoding/json"
	"time"
)

type User struct {
	Id          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Username    string    `json:"username"`
	GoogleId    string    `json:"google_id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
}

type World struct {
	Id        int64           `json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Title     string          `json:"title"`
	Thumbnail *string         `json:"thumbnail"`
	MapData   json.RawMessage `json:"map_data"`
}

type MapAsset struct {
	Id        int64           `json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	WorldId   *int64          `json:"world_id"`
	Title     string          `json:"title"`
	Type      string          `json:"type"`
	Data      json.RawMessage `json:"data"`
}
