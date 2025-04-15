# fetchBidHistory API

All examples show cURL and [HTTPie](https://httpie.io/cli) snippets.

## Initiate fetch bid history

### Successful request

```sh
curl http://localhost:8080/api/history

http GET :8080/api/history
```

The response should be a 200 OK with the following JSON body:

```json
[
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
    }
]
```

### Failed request

```sh
curl http://localhost:8080/api/history

http GET :8080/api/history
```

The response should be a 404 NotFound with the following JSON body:

```json
{
    "error": "No bid history"
}
```