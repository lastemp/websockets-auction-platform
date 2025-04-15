# processAuctionBids websocket API

All examples show cURL and [HTTPie](https://httpie.io/cli) snippets.

## Initiate process auction bids through websocket

### Successful request

```
http method: http POST
Content-Type: application/json
Url: ws://localhost:8080/
Json: {"type": "NEW_BID", "itemId": 1, "bidAmount": 130, "bidder": "MLH"}
```

The response should contain the following JSON body:

```json
{
    "bid": {
        "id": 3,
        "itemId": 1,
        "bidder": "MLH",
        "amount": 130,
        "timeStamp": "2025-04-15T09:56:32Z"
    },
    "item": {
        "id": 1,
        "name": "Vintage Watch",
        "description": "Classic timepiece from 1950s",
        "startingBid": 100,
        "currentBid": 130,
        "bids": [
            {
                "id": 1,
                "itemId": 1,
                "bidder": "MLH",
                "amount": 110,
                "timeStamp": "2025-04-15T09:52:26Z"
            },
            {
                "id": 2,
                "itemId": 1,
                "bidder": "Jane",
                "amount": 120,
                "timeStamp": "2025-04-15T09:53:52Z"
            },
            {
                "id": 3,
                "itemId": 1,
                "bidder": "MLH",
                "amount": 130,
                "timeStamp": "2025-04-15T09:56:32Z"
            }
        ]
    },
    "type": "BID_UPDATE"
}
```

### Failed request

```
http method: http POST
Content-Type: application/json
Url: ws://localhost:8080/
Json: {"type": "NEW_BID", "itemId": 0, "bidAmount": 0, "bidder": "MLH"}
```

The response should contain the following JSON body:

```json
{
    "message": "Missing required fields",
    "type": "ERROR"
}
```