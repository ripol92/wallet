package wallet

import (
	"github.com/google/uuid"
	"github.com/ripol92/wallet/pkg/types"
	"reflect"
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

func TestService_Reject_success(t *testing.T) {
	svc := Service{}
	acc1, _ := svc.RegisterAccount("+9929888444444")
	acc2, _ := svc.RegisterAccount("+9929888444445")
	acc3, _ := svc.RegisterAccount("+9929888444446")

	_ = svc.Deposit(acc1.ID, types.Money(100))
	_ = svc.Deposit(acc2.ID, types.Money(100))
	_ = svc.Deposit(acc3.ID, types.Money(100))

	payment1, _ := svc.Pay(acc1.ID, types.Money(10), types.PaymentCategory("mobile"))
	svc.Pay(acc2.ID, types.Money(10), types.PaymentCategory("mobile"))
	svc.Pay(acc3.ID, types.Money(10), types.PaymentCategory("mobile"))

	rejectError := svc.Reject(payment1.ID)
	if rejectError != nil {
		t.Error(rejectError)
	}

	rejectedAccount, _ := svc.FindAccountByID(acc1.ID)
	if rejectedAccount.Balance != 100 {
		t.Error("Wrong balance")
	}

	rejectedPayment, _ := svc.FindPaymentByID(payment1.ID)
	if rejectedPayment.Status != types.PaymentStatusFail {
		t.Error("Wrong status")
	}
}

func TestService_Reject_fail(t *testing.T) {
	svc := Service{}
	acc1, _ := svc.RegisterAccount("+9929888444444")
	acc2, _ := svc.RegisterAccount("+9929888444445")
	acc3, _ := svc.RegisterAccount("+9929888444446")

	_ = svc.Deposit(acc1.ID, types.Money(100))
	_ = svc.Deposit(acc2.ID, types.Money(100))
	_ = svc.Deposit(acc3.ID, types.Money(100))

	svc.Pay(acc1.ID, types.Money(10), types.PaymentCategory("mobile"))
	svc.Pay(acc2.ID, types.Money(10), types.PaymentCategory("mobile"))
	svc.Pay(acc3.ID, types.Money(10), types.PaymentCategory("mobile"))

	rejectError := svc.Reject(uuid.New().String())
	if rejectError != ErrPaymentNotFound {
		t.Error("Payment must be not found")
	}
}

func TestService_Repeat(t *testing.T) {
	svc := Service{}
	acc1, _ := svc.RegisterAccount("+9929888444444")
	acc2, _ := svc.RegisterAccount("+9929888444445")
	acc3, _ := svc.RegisterAccount("+9929888444446")

	_ = svc.Deposit(acc1.ID, types.Money(100))
	_ = svc.Deposit(acc2.ID, types.Money(100))
	_ = svc.Deposit(acc3.ID, types.Money(100))

	payment1, _ := svc.Pay(acc1.ID, types.Money(10), types.PaymentCategory("mobile"))
	svc.Pay(acc2.ID, types.Money(10), types.PaymentCategory("mobile"))
	svc.Pay(acc3.ID, types.Money(10), types.PaymentCategory("mobile"))

	repeatedPayment, err := svc.Repeat(payment1.ID)
	if err != nil {
		t.Error(err)
	}

	repeatedPayment.ID = payment1.ID
	if !reflect.DeepEqual(repeatedPayment, payment1) {
		t.Error("Payment not repeated")
	}

	checkAccount, errAccount := svc.FindAccountByID(repeatedPayment.AccountID)
	if errAccount != nil || checkAccount == nil {
		t.Error(err)
	}
	if checkAccount.Balance != types.Money(80) {
		t.Error("Wrong balance")
	}
}
