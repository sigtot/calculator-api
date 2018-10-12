# Calculator API
This API evaluates basic arithmetic expressions sent with a POST request.

## Usage
Clone the repo and enter the directory
### Start server
```bash
go run main.go
```

### Run tests
Unit tests for the internal functions relating to evaluation and file handling can be run like so:
```bash
go test
cd file_handling && go test
```

There's also a simple bash script made to test the API:
```bash
sh api_test.sh
```
## Available endpoints
### Calc
Post a new calculation
```
POST /api/calc
```
The endpoint accepts a json object with the single attribute `expression`. 

### History
Get global calculation history
```
GET /api/history/{num}
```
Num is not necessary and the API will default to 5 previous expressions. The API then returns a list of objects with the fields `expression` and `result`.  

## Example usage (with curl)
```bash
$ curl -X POST -H "Content-Type: application/json" -d '{"expression": "-1 * (2 * 6 / 3)"}' http://127.0.0.1:8000/api/calc/

{"result":-4}
```

```bash
$ curl http://127.0.0.1:8000/api/history/2
[{"expression":"3 / 5","result":0.6},{"expression":"3 / 5","result":0.6}]
```
