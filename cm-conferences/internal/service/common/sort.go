package service

type Order string

const (
	ASC  Order = "asc"
	DESC Order = "desc"
)

type Sort struct {
	Column string
	Order  Order
}
