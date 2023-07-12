# nft marketplace demo



## Dependencies
    - redis for caching id increment
    - postgres for transaction records and listings
    - docker for testing

## Caching
    - incremental id for transactions are cached, if cache miss/error trying to get id it is treated as error
    - lazy loading for listings

## Database
    - repo pattern powered by gorm/gorm postgres driver (repo pattern maybe overkill for simple domain but designed for granular testability and extensibility)

## Testing
    - docker for functional tests
    - go test for integration/unit tests ( TODO )

### To function test locally with docker compose
```
go install
docker compose up
```