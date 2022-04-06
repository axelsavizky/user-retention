# User Retention

## Design decisions:
- Tried to divide it in domains as much as I could, but it's not a problem with a lot of domains.
- Tested everything with unit tests except where I open the input file.
- I thought the exercise to optimize it in terms of speed complexity. This decision might take more memory or might make it 
difficult to parallelize.
- I read the CSV only one line at a time because the file might not fit in memory. To not make it slower, I start processing
while I read the file.
- Many tests have table tests. I think is a good practice in Golang.
- The userRetention struct was planned to be immutable. That's why `ProcessRecords` returns a new struct.


## Libraries used
Besides using the native library from Golang, I used `github.com/stretchr/testify` to make tests easier.

## Commands

### Run
To run it, just run `go run main.go path/to/input/file`

### Tests
To run tests, from the root directory run `go test ./...`
