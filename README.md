#
A video series by Nic Jackson covering the basics on how to create microservices in Go: https://www.youtube.com/playlist?list=PLmD8u-IFdreyh6EUfevBcbiuCKzFk0EW_  

Video #1: https://www.youtube.com/watch?v=VzBGi_n65iU   
Introduction and setting up a http web server.  

Video #2: https://www.youtube.com/watch?v=hodOppKJm5Y  
Extract handlers, use channel to create blocking code that waits for interrupt or kill commands before gracefully shutting down.  

Video #3: https://www.youtube.com/watch?v=eBeqtmrvVpg  
REST, starting to create an "online coffee shop". Implementing the products GET handler.

Video #4: https://www.youtube.com/watch?v=UZbHLVsjpF0  
More REST -> Update product, create product, etc.  

Video #5: https://www.youtube.com/watch?v=DD3JlT_u0DM  
Using Gorilla framework for rest services. Adding middleware for validation  

Video #6: https://www.youtube.com/watch?v=gE8_-8KoOLc  
JSON validation and an initial unit test.

Video #7: https://www.youtube.com/watch?v=07XhTqE-j8k  
Adding swagger to the rest api.

Video #8: https://www.youtube.com/watch?v=Zn4joNjqBFc  
Auto-generating HTTP clients from Swagger files. Kind of skipped this one.

Video #9: https://www.youtube.com/watch?v=RlYoy_RiYPw  
Setting up CORS using Gorilla.  

Video #10: https://www.youtube.com/watch?v=ctmhYJpGsgU  
File uploads etc.  

Video #11: https://www.youtube.com/watch?v=ctmhYJpGsgU  
Multipart uploads.  

Video #12: https://www.youtube.com/watch?v=GtSg1H7SU5Y  
Gzip compression for HTTP responses.  

Video #13: https://www.youtube.com/watch?v=pMgty_RYIOc  
gRPC and protocol buffers  

Video #14: https://www.youtube.com/watch?v=oTBcd5J0VYU  
gRPC client connections

Video #15 (Part 1): https://www.youtube.com/watch?v=Vl88R9acq-Y
Refactoring, retrieving exchange rates dynamically, first test.  

Video #15 (Part 2): https://www.youtube.com/watch?v=QBl8LZ0Rems  
Refactoring, cleanup product api.  

Video #15 (Part 3): https://www.youtube.com/watch?v=ARvOyAsuFog
Refactoring, fixed quite a few bugs and typos, completed exchange functionality.  

Video #16 (Part 1): https://www.youtube.com/watch?v=4ohwkWVgEZM  
Bi-directional streaming with gRPC.

Video $16 (Part 2): https://www.youtube.com/watch?v=MT5tXSKa-KY
Bi-directional streaming with gRPC continued. Added bi-directional stream between product api and currency service. Exchange rates are now kept in sync.

## Setup notes
Avoid firewall prompt on startup by explicitly setting localhost ip for `ListenAndServe`.

## 
Can override json output by using struct tags:  
```
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// Outputs
{"id":1,"name":"Latte","description":"Frothy milky coffee","price":2.45,"sku":"abc323"}
```

## Gorilla
https://www.gorillatoolkit.org/  
Install: `go get github.com/gorilla/mux`  


## Validator
https://github.com/go-playground/validator  


## Swagger   
https://goswagger.io/   
Install: `go get -u github.com/go-swagger/go-swagger/cmd/swagger`  
Generate: `make swagger` from product-api dir  

*Empty swagger yaml generated*  
This can happen if there is a blank line between the comment and the package declaration:
```
...
//     Consumes: 
//     - application/json
// swagger:meta

package handlers // NOTE: Line break above ^^^

import (
...
```

Change it to this:  
```
//     Consumes: 
//     - application/json
// swagger:meta
package handlers  // NOTE: Line break above removed

import (
```

## gRPC and Protobuffs  
Install protoc etc
```
$ brew install protobuf
$ protoc --version  # Ensure compiler version is 3+
```

Run following command to generate (or add to makefile)
```
protoc -I protos/ protos/currency.proto --go-grpc_out=protos/currency
```

### gRPC Curl
Install: `brew install grpcurl`  
List services: `grpcurl --plaintext localhost:9092 list`  
Describe service
```
grpcurl --plaintext localhost:9092 describe currency.Currency.GetRate
currency.Currency.GetRate is a method:
rpc GetRate ( .currency.RateRequest ) returns ( .currency.RateResponse );
```
Describe message
```
grpcurl --plaintext localhost:9092 describe currency.RateRequest 
currency.RateRequest is a message:
message RateRequest {
  string Base = 1;
  string Destination = 2;
}
```

Get a template  
```
grpcurl --plaintext --msg-template localhost:9092 describe currency.RateRequest 
currency.RateRequest is a message:
message RateRequest {
  .currency.Currencies Base = 1;
  .currency.Currencies Destination = 2;
}

Message template:
{
  "Base": "USD",
  "Destination": "JPY"
}
```

Send a request
```
grpcurl --plaintext -d '{"Ba
se": "GBP", "Destination": "USD" }' localhost:9092 currency.Currency.GetRate
{
  "Rate": 2.5
}
```

Stream response
```
grpcurl --plaintext --msg-template -d @ localhost:9092 currency.Currency/SubscribeRates
```

## Testing  
Run a test in currency project: `go test -v ./data`  
