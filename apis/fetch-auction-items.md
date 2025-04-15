# fetchAuctionItems API

All examples show cURL and [HTTPie](https://httpie.io/cli) snippets.

## Initiate fetch auction items

```sh
curl http://localhost:8080/api/items

http GET :8080/api/items
```

The response should be a 200 OK with the following JSON body:

```json
[
    {
        "id": 1,
        "name": "Vintage Watch",
        "description": "Classic timepiece from 1950s",
        "startingBid": 100,
        "currentBid": 100,
        "bids": []
    },
    {
        "id": 2,
        "name": "Art Painting",
        "description": "Original Artwork in Skribble.io",
        "startingBid": 150,
        "currentBid": 150,
        "bids": []
    }
]
```
