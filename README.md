# Calculator API
This API evaluates basic arithmetic expressions sent with a POST request.

## Usage
Clone and enter the repository
### Start server
```bash
go run main.go
```

### Run tests
```bash
go test
```

## Example usage (with curl)
```bash
$ curl -X POST -H "Content-Type: application/json" -d '{"expression": "-1 * (2 * 6 / 3)"}' http://127.0.0.1:8000/api/calc/

{"result":-4}
```
