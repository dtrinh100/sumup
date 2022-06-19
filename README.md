## Running the program
1. cd into this program's directory
2. Run `go mod download` to install dependencies (if you want to run unit tests).
3. Run `go run main.go inputfile.csv > outputfile.csv`

## Testing
To run unit tests, use the command: `go test`. I also created a sample CSV file that I used to test with as well called `sample-csv.csv`.

## Design Overview/Open Questions
For simplicity purposes, I kept all logic in the `main.go` file, but in a real program, I would seperate out the logic into their own package. For client and tranasction lookup, I chose to use the map data structure and using the id for lookup, since they are unique. This also allows me to look up data in constant time. I also made sure to validate input, like making sure the deposit/withdraw amount is more than 0.

Open questions for this project are: If the account is locked, can we still withdraw from the client's account? Can a locked account still get a deposit?