package types

type Money int64

type PaymentCategory string

type PaymentStatus string

const (
	PaymentStatusOk PaymentStatus = "OK"
	PaymentStatusFail PaymentStatus = "FAIL"
	PaymentStatusInProgress PaymentStatus = "INPROGRESS"
)

type Payment struct {
	ID        string
	Amount    Money
	AccountID int64
	Category  PaymentCategory
	Status    PaymentStatus
}

type Favorite struct {
	ID        string
	Name	  string
	Amount    Money
	AccountID int64
	Category  PaymentCategory
}

type Phone string

type Account struct {
	ID int64
	Phone Phone
	Balance Money
}