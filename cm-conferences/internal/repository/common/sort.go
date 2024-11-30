package common

type Order string

const (
	ASC  Order = "asc"
	DESC Order = "desc"
)

type Sort struct {
	Column string `json:"column"`
	Order  Order  `json:"order"`
}
