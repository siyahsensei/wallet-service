package transaction

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	Deposit    TransactionType = "DEPOSIT"
	Withdrawal TransactionType = "WITHDRAWAL"
	Transfer   TransactionType = "TRANSFER"
	Buy        TransactionType = "BUY"
	Sell       TransactionType = "SELL"
	Dividend   TransactionType = "DIVIDEND"
	Interest   TransactionType = "INTEREST"
	Fee        TransactionType = "FEE"
	Income     TransactionType = "INCOME"
	Expense    TransactionType = "EXPENSE"
	Tax        TransactionType = "TAX"
	Rebalance  TransactionType = "REBALANCE"
	Split      TransactionType = "SPLIT"
	Merger     TransactionType = "MERGER"
	Staking    TransactionType = "STAKING"
	Mining     TransactionType = "MINING"
	Airdrop    TransactionType = "AIRDROP"
	Lending    TransactionType = "LENDING"
	Borrowing  TransactionType = "BORROWING"
	Repayment  TransactionType = "REPAYMENT"
)

type Transaction struct {
	ID              uuid.UUID       `json:"id" db:"id"`
	UserID          uuid.UUID       `json:"userId" db:"user_id"`
	AccountID       uuid.UUID       `json:"accountId" db:"account_id"`
	AssetID         *uuid.UUID      `json:"assetId,omitempty" db:"asset_id"`
	Type            TransactionType `json:"type" db:"type"`
	Amount          float64         `json:"amount" db:"amount"`
	Quantity        float64         `json:"quantity" db:"quantity"`
	Price           float64         `json:"price" db:"price"`
	Fee             float64         `json:"fee" db:"fee"`
	Currency        string          `json:"currency" db:"currency"`
	Description     string          `json:"description" db:"description"`
	Category        string          `json:"category" db:"category"`
	Date            time.Time       `json:"date" db:"date"`
	ToAccountID     *uuid.UUID      `json:"toAccountId,omitempty" db:"to_account_id"`
	TransactionHash string          `json:"transactionHash,omitempty" db:"transaction_hash"`
	CreatedAt       time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time       `json:"updatedAt" db:"updated_at"`
}

func NewTransaction(
	userID uuid.UUID,
	accountID uuid.UUID,
	assetID *uuid.UUID,
	transactionType TransactionType,
	amount float64,
	quantity float64,
	price float64,
	fee float64,
	currency string,
	description string,
	category string,
	date time.Time,
	toAccountID *uuid.UUID,
	transactionHash string,
) *Transaction {
	now := time.Now()
	return &Transaction{
		ID:              uuid.New(),
		UserID:          userID,
		AccountID:       accountID,
		AssetID:         assetID,
		Type:            transactionType,
		Amount:          amount,
		Quantity:        quantity,
		Price:           price,
		Fee:             fee,
		Currency:        currency,
		Description:     description,
		Category:        category,
		Date:            date,
		ToAccountID:     toAccountID,
		TransactionHash: transactionHash,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func (t *Transaction) TotalAmount() float64 {
	if t.IsDebit() {
		return t.Amount + t.Fee
	}
	return t.Amount - t.Fee
}

func (t *Transaction) IsDebit() bool {
	return t.Type == Withdrawal || t.Type == Buy || t.Type == Transfer ||
		t.Type == Fee || t.Type == Expense || t.Type == Tax ||
		t.Type == Borrowing || t.Type == Repayment || t.Type == Lending
}

func (t *Transaction) IsCredit() bool {
	return t.Type == Deposit || t.Type == Sell || t.Type == Dividend ||
		t.Type == Interest || t.Type == Income || t.Type == Staking ||
		t.Type == Mining || t.Type == Airdrop || t.Type == Borrowing
}

func (t *Transaction) IsTransfer() bool {
	return t.Type == Transfer
}

func (t *Transaction) IsAssetTransaction() bool {
	return t.Type == Buy || t.Type == Sell || t.Type == Dividend ||
		t.Type == Split || t.Type == Merger || t.Type == Staking ||
		t.Type == Mining || t.Type == Airdrop
}
