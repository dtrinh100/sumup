package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type TransactionData struct {
	TransactionType string
	ClientId        int
	TransactionId   int
	Amount          float64
}

type Transaction struct {
	Id              int
	ClientID        int
	TransactionType string
	Amount          float64
	IsDisputed      bool
}

type Client struct {
	Id        int
	Available float64
	Held      float64
	Total     float64
	IsLocked  bool
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.FieldsPerRecord = -1

	// Skip first row
	if _, err := csvReader.Read(); err != nil {
		log.Fatal(err)
	}

	rows, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Map used to store client data based on the client's ID
	clientMap := make(map[int]*Client)

	// Map used to store individual transaction data based on the transaction's ID
	transactionMap := make(map[int]*Transaction)

	// Process the data and converts each row into a Transction struct for easier processing
	for _, row := range rows {
		clientId, err := strconv.Atoi(strings.TrimSpace(row[1]))

		if err != nil {
			log.Fatal(err)
		}

		transactionId, err := strconv.Atoi(strings.TrimSpace(row[2]))

		if err != nil {
			log.Fatal(err)
		}

		amount := 0.0
		if len(row) > 3 {
			amount, err = strconv.ParseFloat(strings.TrimSpace(row[3]), 64)

			if err != nil {
				log.Fatal(err)
			}
		}

		transaction := TransactionData{
			TransactionType: row[0],
			ClientId:        clientId,
			TransactionId:   transactionId,
			Amount:          amount,
		}

		ProcessTransaction(transaction, clientMap, transactionMap)
	}

	OutputTransactionsData(clientMap)
}

// ProcessTransaction processes a transaction. Records all deposit and withdraw transactions.
func ProcessTransaction(data TransactionData, clientMap map[int]*Client, transactionMap map[int]*Transaction) {
	// Check the client map to see if id exists, if not create a new client
	if _, ok := clientMap[data.ClientId]; !ok {
		client := &Client{Id: data.ClientId}
		clientMap[data.ClientId] = client
	}

	switch data.TransactionType {
	case "deposit":
		processDeposit(data.Amount, clientMap[data.ClientId])
		transactionMap[data.TransactionId] = &Transaction{Id: data.TransactionId, ClientID: data.ClientId, TransactionType: data.TransactionType, Amount: data.Amount, IsDisputed: false}
	case "withdrawal":
		isSuccessFul := processWithdrawal(data.Amount, clientMap[data.ClientId])
		if isSuccessFul {
			transactionMap[data.TransactionId] = &Transaction{Id: data.TransactionId, ClientID: data.ClientId, TransactionType: data.TransactionType, Amount: data.Amount, IsDisputed: false}
		}
	case "dispute":
		processDispute(data.TransactionId, transactionMap, clientMap[data.ClientId])
	case "resolve":
		processResolve(data.TransactionId, transactionMap, clientMap[data.ClientId])
	case "chargeback":
		processChargeback(data.TransactionId, transactionMap, clientMap[data.ClientId])
	}

}

// processDeposit is a helper function that adds the transaction amount to the overall client's current total and available funds.
func processDeposit(amount float64, client *Client) {
	if amount > 0 {
		client.Available += amount
		client.Total += amount
	}
}

// processWithdrawal is a helper function that subtracts the transaction amount from the overall client's current total and available funds
// if the amount is lower than the available funds and returns true, otherwise it returns false without doing anything
func processWithdrawal(amount float64, client *Client) bool {
	isSuccessFul := false
	if client.Available >= amount && amount > 0 {
		client.Available -= amount
		client.Total -= amount
		isSuccessFul = true
	}

	return isSuccessFul
}

// processDispute is a helper function that subtracts available funds from the amount disputed. Held funds increases
// by the amount disputed. Total funds does not change. If the transaction id does not exist within the system
// then simply ignore it.
func processDispute(transactionId int, transactionMap map[int]*Transaction, client *Client) {
	// Make sure the transaction exists AND the transaction actually belongs to the client
	if transaction, ok := transactionMap[transactionId]; ok && client.Id == transaction.ClientID {
		client.Available -= transaction.Amount
		client.Held += transaction.Amount
		transaction.IsDisputed = true
	}
}

// processResolve is a helper function that subtracts held funds from the amount previously disputed. Available funds increases
// by the amount previously disputed. Total funds does not change. If the transaction id does not exist within the system
// then simply ignore it.
func processResolve(transactionId int, transactionMap map[int]*Transaction, client *Client) {
	// Make sure the transaction exists, actually belongs to the client, and it has been disputed
	if transaction, ok := transactionMap[transactionId]; ok && client.Id == transaction.ClientID && transaction.IsDisputed {
		client.Held -= transaction.Amount
		client.Available += transaction.Amount
		transaction.IsDisputed = false
	}
}

// processChargeback is a helper function that subtracts held funds from the amount previously disputed. Total funds decreases
// by the amount previously disputed. Any chargeback will lock the account. If the transaction id does not exist
// within the system then simply ignore it.
func processChargeback(transactionId int, transactionMap map[int]*Transaction, client *Client) {
	// Make sure the transaction exists, actually belongs to the client, and it has been disputed
	if transaction, ok := transactionMap[transactionId]; ok && client.Id == transaction.ClientID && transaction.IsDisputed {
		client.Held -= transaction.Amount
		client.Total -= transaction.Amount
		client.IsLocked = true
	}
}

// OutputTransactionsData outputs transactions data to a CSV
func OutputTransactionsData(clientMap map[int]*Client) {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	row := []string{"client", "available", "held", "total", "locked"}

	err := w.Write(row)

	if err != nil {
		log.Fatal(err)
	}

	var data [][]string
	for _, client := range clientMap {
		row := []string{strconv.Itoa(client.Id), fmt.Sprintf("%.4f", client.Available), fmt.Sprintf("%.4f", client.Held), fmt.Sprintf("%.4f", client.Total), strconv.FormatBool(client.IsLocked)}
		data = append(data, row)
	}

	err = w.WriteAll(data)

	if err != nil {
		log.Fatal(err)
	}
}
