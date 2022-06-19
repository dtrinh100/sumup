package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessDepositInvalidInputs(t *testing.T) {
	mockClient := &Client{Total: 10.0, Available: 10.0, Held: 0.0, IsLocked: false}
	processDeposit(-10.0, mockClient)
	assert.Equal(t, 10.0, mockClient.Available, "Negative amount values should not change client data")
	assert.Equal(t, 10.0, mockClient.Total, "Negative amount values should not change client data")
}

func TestProcessDeposit(t *testing.T) {
	mockClient := &Client{Total: 10.0, Available: 10.0, Held: 0.0, IsLocked: false}
	processDeposit(1.1235, mockClient)
	assert.Equal(t, 11.1235, mockClient.Available, "Deposits should increase the available funds")
	assert.Equal(t, 11.1235, mockClient.Total, "Deposits should increase the total funds")
}

func TestProcessWithdrawalInvalidInputs(t *testing.T) {
	mockClient := &Client{Total: 10.0, Available: 10.0, Held: 0.0, IsLocked: false}
	isSuccessful := processWithdrawal(-10.0, mockClient)
	assert.Equal(t, 10.0, mockClient.Available, "Negative amount values should not change client data")
	assert.Equal(t, 10.0, mockClient.Total, "Negative amount values should not change client data")
	assert.False(t, isSuccessful, "The withdrawal should fail on negative values")

	isSuccessful = processWithdrawal(100.0, mockClient)
	assert.Equal(t, 10.0, mockClient.Available, "Clients are not allowed to withdraw more than their available funds")
	assert.Equal(t, 10.0, mockClient.Total, "Clients are not allowed to withdraw more than their total funds")
	assert.False(t, isSuccessful, "The withdrawal should fail if amount exceeds available funds")
}

func TestProcessWithdrawal(t *testing.T) {
	mockClient := &Client{Total: 10.0, Available: 10.0, Held: 0.0, IsLocked: false}
	isSuccessful := processWithdrawal(5.01, mockClient)
	assert.Equal(t, 4.99, mockClient.Available, "Withdrawals should decrease the available funds")
	assert.Equal(t, 4.99, mockClient.Total, "Withdrawals should decrease the total funds")
	assert.True(t, isSuccessful, "Withdrawal should be successful")

	mockClient = &Client{Total: 10.0, Available: 10.0, Held: 0.0, IsLocked: false}
	processWithdrawal(10.0, mockClient)
	assert.Equal(t, 0.0, mockClient.Available, "Withdrawals should decrease the available funds")
	assert.Equal(t, 0.0, mockClient.Total, "Withdrawals should decrease the total funds")
	assert.True(t, isSuccessful, "Withdrawal should be successful")
}

func TestProcessDisputeInvalidInputs(t *testing.T) {
	mockClient := &Client{Total: 10.0, Available: 10.0, Held: 0.0, IsLocked: false}
	mockTransactionMap := make(map[int]*Transaction)
	mockTransaction := &Transaction{Id: 1, ClientID: 1, TransactionType: "deposit", Amount: 10, IsDisputed: false}
	mockTransactionMap[1] = mockTransaction
	processDispute(2, mockTransactionMap, mockClient)
	assert.Equal(t, 10.0, mockClient.Available, "Available funds should not change on nonexisting transaction")
	assert.Equal(t, 0.0, mockClient.Held, "Held funds should not change on nonexisting transaction")
	assert.False(t, mockTransaction.IsDisputed, "IsDisputed should not change on nonexisting transaction")

	mockClient = &Client{Id: 2, Total: 10.0, Available: 10.0, Held: 0.0, IsLocked: false}
	mockTransactionMap = make(map[int]*Transaction)
	mockTransaction = &Transaction{Id: 1, ClientID: 1, TransactionType: "deposit", Amount: 10, IsDisputed: false}
	mockTransactionMap[1] = mockTransaction
	processDispute(1, mockTransactionMap, mockClient)
	assert.Equal(t, 10.0, mockClient.Available, "Available funds should not change on a transaction that does not belong to the client")
	assert.Equal(t, 0.0, mockClient.Held, "Held funds should not change on a transaction that does not belong to the client")
	assert.False(t, mockTransaction.IsDisputed, "IsDisputed should not change on a transaction that does not belong to the client")
}

func TestProcessDispute(t *testing.T) {
	mockClient := &Client{Id: 1, Total: 10.0, Available: 10.0, Held: 0.0, IsLocked: false}
	mockTransactionMap := make(map[int]*Transaction)
	mockTransaction := &Transaction{Id: 1, ClientID: 1, TransactionType: "deposit", Amount: 6.0, IsDisputed: false}
	mockTransactionMap[1] = mockTransaction
	processDispute(1, mockTransactionMap, mockClient)
	assert.Equal(t, 4.0, mockClient.Available, "Available funds should be subtracted from dispute amount")
	assert.Equal(t, 6.0, mockClient.Held, "Held funds should increase from dispute amount")
	assert.Equal(t, 10.0, mockClient.Total, "Total funds should not change")
	assert.True(t, mockTransaction.IsDisputed, "Transaction's isDisputed status should update to true")

	mockClient = &Client{Id: 1, Total: 5.0, Available: 5.0, Held: 0.0, IsLocked: false}
	mockTransactionMap = make(map[int]*Transaction)
	mockTransaction = &Transaction{Id: 1, ClientID: 1, TransactionType: "withdrawal", Amount: 6.0, IsDisputed: false}
	mockTransactionMap[1] = mockTransaction
	processDispute(1, mockTransactionMap, mockClient)
	assert.Equal(t, -1.0, mockClient.Available, "Available funds should be subtracted from dispute amount")
	assert.Equal(t, 6.0, mockClient.Held, "Held funds should increase from dispute amount")
	assert.Equal(t, 5.0, mockClient.Total, "Total funds should not change")
	assert.True(t, mockTransaction.IsDisputed, "Transaction's isDisputed status should update to true")
}

func TestProcessResolveInvalidInputs(t *testing.T) {
	mockClient := &Client{Total: 10.0, Available: 10.0, Held: 0.0, IsLocked: false}
	mockTransactionMap := make(map[int]*Transaction)
	mockTransaction := &Transaction{Id: 1, ClientID: 1, TransactionType: "deposit", Amount: 10, IsDisputed: true}
	mockTransactionMap[1] = mockTransaction
	processResolve(2, mockTransactionMap, mockClient)
	assert.Equal(t, 10.0, mockClient.Available, "Available funds should not change on nonexisting transaction")
	assert.Equal(t, 0.0, mockClient.Held, "Held funds should not change on nonexisting transaction")

	mockClient = &Client{Id: 2, Total: 10.0, Available: 10.0, Held: 0.0, IsLocked: false}
	mockTransactionMap = make(map[int]*Transaction)
	mockTransaction = &Transaction{Id: 1, ClientID: 1, TransactionType: "deposit", Amount: 10, IsDisputed: true}
	mockTransactionMap[1] = mockTransaction
	processResolve(1, mockTransactionMap, mockClient)
	assert.Equal(t, 10.0, mockClient.Available, "Available funds should not change on a transaction that does not belong to the client")
	assert.Equal(t, 0.0, mockClient.Held, "Held funds should not change on a transaction that does not belong to the client")

	mockClient = &Client{Id: 2, Total: 10.0, Available: 10.0, Held: 0.0, IsLocked: false}
	mockTransactionMap = make(map[int]*Transaction)
	mockTransaction = &Transaction{Id: 1, ClientID: 2, TransactionType: "deposit", Amount: 10, IsDisputed: false}
	mockTransactionMap[1] = mockTransaction
	processResolve(1, mockTransactionMap, mockClient)
	assert.Equal(t, 10.0, mockClient.Available, "Available funds should not change on a transaction that has not been disputed")
	assert.Equal(t, 0.0, mockClient.Held, "Held funds should not change on a transaction that has not been disputed")
}

func TestProcessResolve(t *testing.T) {
	mockClient := &Client{Id: 1, Total: 10.0, Available: 4.0, Held: 6.0, IsLocked: false}
	mockTransactionMap := make(map[int]*Transaction)
	mockTransaction := &Transaction{Id: 1, ClientID: 1, TransactionType: "deposit", Amount: 6.0, IsDisputed: true}
	mockTransactionMap[1] = mockTransaction
	processResolve(1, mockTransactionMap, mockClient)
	assert.Equal(t, 10.0, mockClient.Available, "Available funds should be increased from resolved amount")
	assert.Equal(t, 0.0, mockClient.Held, "Held funds should decrease from resolved amount")
	assert.Equal(t, 10.0, mockClient.Total, "Total funds should not change")
	assert.False(t, mockTransaction.IsDisputed, "Transaction's isDisputed status should update to false")
}

func TestProcessChargebackInvalidInputs(t *testing.T) {
	mockClient := &Client{Total: 10.0, Available: 10.0, Held: 0.0, IsLocked: false}
	mockTransactionMap := make(map[int]*Transaction)
	mockTransaction := &Transaction{Id: 1, ClientID: 1, TransactionType: "deposit", Amount: 10, IsDisputed: true}
	mockTransactionMap[1] = mockTransaction
	processChargeback(2, mockTransactionMap, mockClient)
	assert.Equal(t, 10.0, mockClient.Total, "Total funds should not change on nonexisting transaction")
	assert.Equal(t, 0.0, mockClient.Held, "Held funds should not change on nonexisting transaction")
	assert.False(t, mockClient.IsLocked, "Client account should not be locked on nonexisting transaction")

	mockClient = &Client{Id: 2, Total: 10.0, Available: 10.0, Held: 0.0, IsLocked: false}
	mockTransactionMap = make(map[int]*Transaction)
	mockTransaction = &Transaction{Id: 1, ClientID: 1, TransactionType: "deposit", Amount: 10, IsDisputed: true}
	mockTransactionMap[1] = mockTransaction
	processChargeback(1, mockTransactionMap, mockClient)
	assert.Equal(t, 10.0, mockClient.Total, "Total funds should not change on a transaction that does not belong to the client")
	assert.Equal(t, 0.0, mockClient.Held, "Held funds should not change on a transaction that does not belong to the client")
	assert.False(t, mockClient.IsLocked, "Client account should not be locked on a transaction that does not belong to the client")

	mockClient = &Client{Id: 2, Total: 10.0, Available: 10.0, Held: 0.0, IsLocked: false}
	mockTransactionMap = make(map[int]*Transaction)
	mockTransaction = &Transaction{Id: 1, ClientID: 2, TransactionType: "deposit", Amount: 10, IsDisputed: false}
	mockTransactionMap[1] = mockTransaction
	processChargeback(1, mockTransactionMap, mockClient)
	assert.Equal(t, 10.0, mockClient.Total, "Total funds should not change on a transaction that has not been disputed")
	assert.Equal(t, 0.0, mockClient.Held, "Held funds should not change on a transaction that has not been disputed")
	assert.False(t, mockClient.IsLocked, "Client account should not be locked on a transaction that has not been disputed")
}

func TestProcessChargeback(t *testing.T) {
	mockClient := &Client{Id: 1, Total: 10.0, Available: 4.0, Held: 6.0, IsLocked: false}
	mockTransactionMap := make(map[int]*Transaction)
	mockTransaction := &Transaction{Id: 1, ClientID: 1, TransactionType: "deposit", Amount: 6.0, IsDisputed: true}
	mockTransactionMap[1] = mockTransaction
	processChargeback(1, mockTransactionMap, mockClient)
	assert.Equal(t, 4.0, mockClient.Total, "Total funds should decrease by chargeback amount")
	assert.Equal(t, 0.0, mockClient.Held, "Held funds should decrease by chargeback amount")
	assert.True(t, mockClient.IsLocked, "Client account should be locked after a chargeback")
}
