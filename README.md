#
Video #1: https://www.youtube.com/watch?v=VzBGi_n65iU   
Introduction and setting up a http web server.  

Video #2: https://www.youtube.com/watch?v=hodOppKJm5Y  
Extract handlers, use channel to create blocking code that waits for interrupt or kill commands before gracefully shutting down.  

Video #3: https://www.youtube.com/watch?v=eBeqtmrvVpg  
REST, starting to create an "online coffee shop"  

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