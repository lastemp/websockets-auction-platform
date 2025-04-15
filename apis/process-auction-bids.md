# processAuctionBids API

All examples show cURL and [HTTPie](https://httpie.io/cli) snippets.

## Initiate process auction bids

### Successful request

```sh
curl -d '{"itemId": 1, "bidAmount": 110, "bidder": "MLH"}' -H 'Content-Type: application/json' http://localhost:8080/api/bids

http POST :8080/api/bids itemId=1 bidAmount=110 bidder="MLH"
```

The response should be a 201 Created with the following JSON body:

```json
{
    "id": 1,
    "itemId": 1,
    "bidder": "MLH",
    "amount": 110,
    "timeStamp": "2025-04-15T09:52:26Z"
}
```

### Failed request

```sh
curl -d '{"itemId": 0, "bidAmount": 0, "bidder": "MLH"}' -H 'Content-Type: application/json' http://localhost:8080/api/bids

http POST :8080/api/bids itemId=1 bidAmount=110 bidder="MLH"
```

The response should be a 400 BadRequest with the following JSON body:

```json
{
    "error": "Invalid input data!"
}
```