package wallet

import (
	"testing"
)

func TestService_FindAccountById_success(t *testing.T) {
	svc := Service{}
	account1, _ := svc.RegisterAccount("+9929888444444")
	svc.RegisterAccount("+9929888444445")
	svc.RegisterAccount("+9929888444446")

	_, err1 := svc.FindAccountByID(account1.ID)

	if err1 != nil {
		t.Error("Account not found")
	}
}

func TestService_FindAccountById_notSuccess(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+9929888444444")
	svc.RegisterAccount("+9929888444445")
	svc.RegisterAccount("+9929888444446")

	_, err1 := svc.FindAccountByID(1777)

	if err1 == nil {
		t.Error(err1)
	}
}
