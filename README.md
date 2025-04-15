# websockets-auction-platform

This is a Websockets auction platform system that comprises of one sub system as indicated below;

## backend
This is a RESTful Gin Web API that processes incoming requests. It receives posted messages from client applications and then does processesing.

Currently this RESTful API supports: 
- Fetch auction items
- Fetch auction item by Id
- Fetch bid history
- Process auction bids
- Process auction bids through websocket

The above features are all rest api routes except the last feature which is processed on websocket connection.

## Usage

All the following commands assume that your current working directory is _this_ directory. I.e.:

```console
$ pwd
.../websockets-auction-platform
```
  
1. Create a `.env` file in the directory [backend](./backend/) :

   ```ini
   SERVER_ADDR=127.0.0.1:8080
   ```
   
1. Run the application in directory [backend](./backend/) :

   ```sh
   go run .
   ```
   
1. Using a different terminal send requests to the running server. For example, using [HTTPie]:

   ```sh
   http GET :8080/
   ```

   See [the API documentation pages](./apis/) for more info.