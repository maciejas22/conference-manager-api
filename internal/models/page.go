package models

type Page struct {
	Number int `json:"number"`
	Size   int `json:"size"`
}

type PageInfo struct {
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
	Number     int `json:"number"`
	Size       int `json:"size"`
}
