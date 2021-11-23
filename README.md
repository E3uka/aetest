# AEtest

## How to run
You can start the server by running the command from the top level directory of the project:

```sh
go run cmd/main.go
```

This server is listening on port `3000` by default. You can specify the required port you want it
to run on by using flags.

```sh
# run server on port 8080
go run affordability/cmd/main.go -http=:8080
```

## Simple order request

In another shell instance you can send a simple POST request to the `/submit-order` endpoint to get
an order summary with the total cost of the order. The structure of the JSON request that must be 
sent to the endpoint is found in the [`simple_order.json`](./examples/simple_order.json) file.
Below shows this request made with this same file:

```sh
# be sure to change the port if you are using a custom port
curl -X POST -H "Content-Type: application/json" -d @examples/simple_order.json localhost:3000/submit-order
```

By using a command line JSON processing tool like [jq](https://stedolan.github.io/jq/) you can
"pretty print" the output on your terminal as follows:

```sh
# be sure to change the port if you are using a custom port
curl -X POST -H "Content-Type: application/json" -d @examples/simple_order.json localhost:3000/submit-order | jq .
```

## Getting a single order

The response of making an order request is an object that contains a unique generated order id, a
summary of the order and the total cost of the order. Each successfull processed order is stored
internally and can be queried by making a POST request to the `/get-order` endpoint. The payload 
required is JSON format of the following structure.

```js
// replace the order_id below with one that has been generated from your order.
{"order_id":"36c9b2a4-a1eb-4c6a-9a55-7448898bc09c"}
```

If a supplied order id is invalid or not found the response will be an appropriate error. An
example request is shown below.

```sh
# replace [PAYLOAD] with JSON structure described above
curl -X POST -H "Content-Type: application/json" -d [PAYLOAD] localhost:3000/submit-order
```

## Getting all orders

A list of all orders can be obtained by making a GET request to the `/get-all-orders` endpoint. If
no orders have been previously made, this will return an empty list. An example request is shown
below:

```sh
# replace [PAYLOAD] with JSON structure described above
curl localhost:3000/get-all-orders
```

## Testing
There has been a series of test cases that have been produced. This can be found in the 
[`service_test`](service_test.go) file. You can run the tests with the below command:

```sh
# -cover includes code coverage to the result
go test ./... -cover
```
...

Kind regards,

Ebuka Agbanyim.
ebuka7@outlook.com
