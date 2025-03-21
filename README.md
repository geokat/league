## Building
```
go build ./...
```


## Running
```
go run ./...
```


## Running the test suite
```
go test -v ./...
```

## Testing the functionality

Testing the web (complete) API:
```
curl -F 'file=@/path/matrix.csv' "localhost:8080/echo"
curl -F 'file=@/path/matrix.csv' "localhost:8080/flatten"
curl -F 'file=@/path/matrix.csv' "localhost:8080/sum"
curl -F 'file=@/path/matrix.csv' "localhost:8080/multiply"
curl -F 'file=@/path/matrix.csv' "localhost:8080/invert"
```

Testing the stream (example) API:
```
curl -s -T '/path/matrix.csv' "localhost:8080/stream/echo"
```
