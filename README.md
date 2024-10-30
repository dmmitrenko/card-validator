# Card validator :walking:
This is a service for validating credit card numbers using the Luhn algorithm, expiration date validation, and third-party API checks for issuer information.

**Folder PATH listing**
```
card-validator
├───cards
│   └───mocks
├───cmd
├───domain
├───internal
│   ├───config
│   └───grpc
│       ├───handler
│       ├───middleware
│       └───proto
└───tests
    └───integration
```

## Running the Project
Build the Docker image:
```
docker build -t card-validator .
```

Run the Docker container:
```
docker run -p 9000:9000 card-validator
```

## API Documentation
- `ValidateCard(CardValidationRequest) returns (CardValidationResponse)`
  - `CardValidationRequest`: includes fields `CardNumber`, `ExpirationMonth`, `ExpirationYear`
  - `CardValidationResponse`: includes fields `Valid` (boolean) and `Error` (optional)

### Error Codes:
- `1`: **Invalid card number** - The card number is not valid. *_The card number length must be between 6 and 20 digits._
- `2`: **Invalid month** - The expiration month provided is incorrect.
- `3`: **Invalid year** - The expiration year provided is incorrect. *_There is no strict upper limit for the expiration year._
- `4`: **Card doesn't match Luhn's algorithm** - The card number fails the Luhn algorithm check.
- `5`: **Card is expired** - The card has passed its expiration date.
- `6`: **This BIN is not found** - The card's BIN (Bank Identification Number) is unknown or not found. *_Only the first 6 digits (BIN/IIN) are used for issuer identification._

## Usage
```
grpcurl -plaintext -d '{"CardNumber": "4111111111111111", "ExpirationMonth": 12, "ExpirationYear": 2025}' localhost:9000 proto.CardValidator/ValidateCard
```
## Dependencies
- Go 1.23.2
- Docker
- gRPC
