package models

type Permission struct {
	ID     int    `json:"id,omitempty"`
	Name   string `json:"name" validate:"required"`
	Path   string `json:"path" validate:"required"`
	Method string `json:"method" validate:"required"`
}

type Role struct {
	ID         int          `json:"id,omitempty"`
	Name       string       `json:"name" validate:"required"`
	Permission []Permission `json:"permission" validate:"required"`
}

type GrantPermission struct {
	UserID       int   `json:"user_id,omitempty"`
	PermissionID []int `json:"permission_id,omitempty"`
	RoleID       []int `json:"role_id,omitempty"`
}
