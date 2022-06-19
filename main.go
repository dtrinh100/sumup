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
	Id              int
	Amount          float64
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

	// Process the data and converts each row into a Transction struct for easier processing
	for _, row := range rows {
		clientId, err := strconv.Atoi(strings.TrimSpace(row[1]))

		if err != nil {
			log.Fatal(err)
		}

		id, err := strconv.Atoi(strings.TrimSpace(row[2]))

		if err != nil {
			log.Fatal(err)
		}

		amount, err := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)

		if err != nil {
			log.Fatal(err)
		}

		transaction := TransactionData{
			TransactionType: row[0],
			ClientId:        clientId,
			Id:              id,
			Amount:          amount,
		}

		ProcessTransaction(transaction, clientMap)
	}

	OutputTransactionsData(clientMap)
}

// ProcessTransaction processes a transaction
func ProcessTransaction(data TransactionData, clientMap map[int]*Client) {
	// Check the client map to see if id exists, if not create a new client
	if _, ok := clientMap[data.ClientId]; !ok {
		client := &Client{Id: data.ClientId}
		clientMap[data.ClientId] = client
	}

	switch data.TransactionType {
	case "deposit":
		processDeposit(data.Amount, clientMap[data.ClientId])
	case "withdrawal":
		processWithdrawal(data.Amount, clientMap[data.ClientId])
	}
}

// processDeposit is a helper function that adds the transaction amount to the overall client's current total and available funds.
func processDeposit(amount float64, client *Client) {
	client.Available += amount
	client.Total += amount
}

// processWithdrawal is a helper function that subtracts the transaction amount from the overall client's current total and available funds
// if the amount is lower than the available funds, otherwise it returns without doing anything
func processWithdrawal(amount float64, client *Client) {
	if amount < client.Available {
		client.Available -= amount
		client.Total -= amount
	}
}

// OutputTransactionsData outputs transactions data to a CSV
func OutputTransactionsData(clientMap map[int]*Client) {
	file, err := os.Create("output.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	row := []string{"client", "available", "held", "total", "locked"}

	err = w.Write(row)

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
