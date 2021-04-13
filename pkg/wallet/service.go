package wallet

import (
	"errors"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"github.com/ripol92/wallet/pkg/types"
)

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrNotEnoughBalance = errors.New("not enough balance")

type Service struct {
	nextAccountId int64
	accounts      []*types.Account
	payments      []*types.Payment
	favourites    []*types.Favorite
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}

	s.nextAccountId++
	account := &types.Account{
		ID:      s.nextAccountId,
		Phone:   phone,
		Balance: 0,
	}

	s.accounts = append(s.accounts, account)

	return account, nil
}

func (s *Service) FindAccountByID(accountId int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountId {
			return account, nil
		}
	}

	return nil, ErrAccountNotFound
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}

	return nil, ErrPaymentNotFound
}

func (s *Service) FindFavoriteById(favoriteID string) (*types.Favorite, error) {
	for _, favorite := range s.favourites {
		if favorite.ID == favoriteID {
			return favorite, nil
		}
	}

	return nil, ErrPaymentNotFound
}

func (s *Service) Reject(paymentID string) error {
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			payment.Status = types.PaymentStatusFail
			for _, account := range s.accounts {
				if account.ID == payment.AccountID {
					account.Balance += payment.Amount
				}
			}
			return nil
		}
	}

	return ErrPaymentNotFound
}

func (s *Service) Pay(accountId int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountId {
			account = acc
			break
		}
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount
	paymentId := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentId,
		AccountID: accountId,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}

	s.payments = append(s.payments, payment)
	return payment, nil
}

func (s *Service) Deposit(accountId int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountId {
			account = acc
			break
		}
	}

	if account == nil {
		return ErrAccountNotFound
	}

	account.Balance += amount
	return nil
}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	repeatPayment := &types.Payment{
		ID:        uuid.New().String(),
		Amount:    payment.Amount,
		AccountID: payment.AccountID,
		Category:  payment.Category,
		Status:    payment.Status,
	}

	repeatPayment, err = s.Pay(repeatPayment.AccountID, repeatPayment.Amount, repeatPayment.Category)
	if err != nil {
		return nil, err
	}

	return repeatPayment, nil
}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	favorite := &types.Favorite{
		ID:        uuid.New().String(),
		Name:      name,
		Amount:    payment.Amount,
		AccountID: payment.AccountID,
		Category:  payment.Category,
	}

	s.favourites = append(s.favourites, favorite)

	return favorite, err
}

func (s *Service) PayFromFavorite(favoriteId string) (*types.Payment, error) {
	favorite, err := s.FindFavoriteById(favoriteId)
	if err != nil {
		return nil, err
	}


	payment, err := s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

