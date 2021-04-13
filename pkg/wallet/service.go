package wallet

import (
	"errors"
	_ "github.com/google/uuid"
	"github.com/ripol92/wallet/pkg/types"
)

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrAccountNotFound = errors.New("account not found")

type Service struct {
	nextAccountId int64
	accounts []*types.Account
	payments []*types.Payment
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}

	s.nextAccountId++
	account := &types.Account{
		ID: s.nextAccountId,
		Phone: phone,
		Balance: 0,
	}

	s.accounts = append(s.accounts, account)

	return account, nil
}

func (s *Service) FindAccountById(accountId int64) (*types.Account, error)  {
	for _, account := range s.accounts {
		if account.ID == accountId {
			return account, nil
		}
	}

	return nil, ErrAccountNotFound
}
