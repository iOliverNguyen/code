# Counter

Simple proof-of-concept counter service.

## Quick Start

```
go run counter.go -auto-create
curl -X POST http://localhost:8901/counter
```

## API

#### GET /counter

- Return the current counter value

#### POST /counter

- Generate new value

# License

- [MIT License](https://opensource.org/licenses/mit-license.php)
