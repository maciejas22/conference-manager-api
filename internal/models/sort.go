package models

type Sort struct {
	Column string `json:"column"`
	Order  Order  `json:"order"`
}
