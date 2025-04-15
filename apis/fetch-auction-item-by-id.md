# fetchAuctionItemById API

All examples show cURL and [HTTPie](https://httpie.io/cli) snippets.

## Initiate fetch auction item by id

### Successful request

```sh
curl http://localhost:8080/api/items/1

http GET :8080/api/items/1
```

The response should be a 200 OK  with the following JSON body:

```json
{
    "id": 1,
    "name": "Vintage Watch",
    "description": "Classic timepiece from 1950s",
    "startingBid": 100,
    "currentBid": 100,
    "bids": []
}
```

### Failed request

```sh
curl http://localhost:8080/api/items/45

http GET :8080/api/items/45
```

The response should be a 404 NotFound with the following JSON body:

```json
{
    "error": "Item not found"
}
```
