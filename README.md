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

