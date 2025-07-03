package me

import "neptune/backend/models/user"

type MeResponse struct {
	UserID     string    `json:"UserId"`     // UUID string
	BinusianID string    `json:"BinusianId"` // UUID string
	Username   string    `json:"Username"`
	Name       string    `json:"Name"`
	Role       user.Role // Bind on fetch
}
