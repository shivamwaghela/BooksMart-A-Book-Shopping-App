/*
	Product API in Go
	Riak KV
*/

package main

import (
	"fmt"
	"errors"
	"log"
	"net/http"
	"io/ioutil"
	"time"
	"strings"
	"encoding/json"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	//"github.com/satori/go.uuid"
)

/*
	Go Rest Client:
		Tutorial:	https://medium.com/@marcus.olsson/writing-a-go-client-for-your-restful-api-c193a2f4998c
		Reference:	https://golang.org/pkg/net/http/
*/

/* Riak REST Client */
var debug = true
var server3 = "http://54.153.111.238:8098"
var server2 = "http://54.241.188.5:8098"
var server1 = "http://13.57.205.174:8098" 
//var server4 = "http://35.166.183.128:8098"
//var server5 = "http://35.161.230.225:8098" 

type Client struct {
	Endpoint string
	*http.Client
}

var tr = &http.Transport{
	MaxIdleConns:       10,
	IdleConnTimeout:    30 * time.Second,
	DisableCompression: true,
}

func NewClient(server string) *Client {
	return &Client{
		Endpoint:  	server,
		Client: 	&http.Client{Transport: tr},
	}
}

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	mx := mux.NewRouter()
	initRoutes(mx, formatter)
	n.UseHandler(mx)
	return n
}

// Init Database Connections
func init() {

	// Riak KV Setup	
	c1 := NewClient(server1)
	msg, err := c1.Ping( )
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Riak Ping Server1: ", msg)		
	}

	c2 := NewClient(server2)
	msg, err = c2.Ping( )
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Riak Ping Server2: ", msg)		
	}

	c3 := NewClient(server3)
	msg, err = c3.Ping( )
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Riak Ping Server3: ", msg)		
	}
	
	/*
	c4 := NewClient(server4)
	msg, err = c4.Ping( )
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Riak Ping Server4: ", msg)		
	}

	c5 := NewClient(server5)
	msg, err = c5.Ping( )
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Riak Ping Server5: ", msg)		
	}
	*/
	
}

func (c *Client) Ping() (string, error) {
	resp, err := c.Get(c.Endpoint + "/ping" )
	if err != nil {
		fmt.Println("[RIAK DEBUG] " + err.Error())
		return "Ping Error!", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if debug { fmt.Println("[RIAK DEBUG] GET: " + c.Endpoint + "/ping => " + string(body)) }
	return string(body), nil
}

func (c *Client) AddProduct(key string, prd_inp product) (product, error) {
	
	var prd_nil = product { }
	
	reqbody := "{\"title_register\": \"" + 
		prd_inp.Title + 
		"\",\"author_register\": \"" +
		 prd_inp.Author + 
		 "\",\"image_URL_register\": \"" +
		 prd_inp.Image_URL + 
		 "\",\"price_register\": \"" +
		 prd_inp.Price + 
		 "\",\"quantity_register\": \"" +
		 prd_inp.Quantity + 
		 "\"}"
		 
	req, _  := http.NewRequest("PUT", c.Endpoint + "/buckets/products/keys/"+key+"?pw=2&returnbody=true", strings.NewReader(reqbody) )
	req.Header.Add("Content-Type", "application/json")
	
	resp, err := c.Do(req)	
	if err != nil {
		fmt.Println("[RIAK DEBUG] " + err.Error())
		return prd_nil, err
	}	
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	
	if debug { fmt.Println("[RIAK DEBUG] PUT: " + c.Endpoint + "/buckets/products/keys/"+key+"?pw=2&returnbody=true => " + string(body)) }
	if (strings.Contains(string(body), "PW-value unsatisfied: ")){
		fmt.Println("Inside error")
		err := errors.New("PW-value unsatisfied")
		return prd_nil, err
	}
	
	var prd = product { }
	if err := json.Unmarshal(body, &prd); err != nil {
		fmt.Println("[RIAK DEBUG] JSON unmarshaling failed: %s", err)
		return prd_nil, err
	}
	return prd, nil
}

func (c *Client) GetProduct(key string) (product, error) {
	var prd_nil = product {}
	resp, err := c.Get(c.Endpoint + "/buckets/products/keys/"+key )
	if err != nil {
		fmt.Println("[RIAK DEBUG] " + err.Error())
		return prd_nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if debug { fmt.Println("[RIAK DEBUG] GET: " + c.Endpoint + "/buckets/products/keys/"+key +" => " + string(body)) }
	var prd = product { }
	if err := json.Unmarshal([]byte(body), &prd); err != nil {
		fmt.Println("[RIAK DEBUG] JSON unmarshaling failed: %s", err)
		return prd_nil, err
	}
	//fmt.Println(prd)
	return prd, nil
}

func (c *Client) GetProducts() ([]product, error) {
	var prd_nil []product
	
	resp, err := c.Get(c.Endpoint + "/buckets/products/keys/allproducts" )
	if err != nil {
		fmt.Println("[RIAK DEBUG] " + err.Error())
		return prd_nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if debug { fmt.Println("[RIAK DEBUG] GET: " + c.Endpoint + "/buckets/products/keys/allproducts" + " => " + string(body)) }
	
	var prd_array []product
	if err := json.Unmarshal([]byte(body), &prd_array); err != nil {
		fmt.Println("[RIAK DEBUG] JSON unmarshaling failed: %s", err)
		return prd_nil, err
	}
	return prd_array, nil
}

func (c *Client) Updateproduct(key string, value string) (product, error) {

	resp, _ := c.Get(c.Endpoint + "/buckets/products/keys/"+key )
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	
	var prd_nil = product { }
	var prd1 = product { }
	json.Unmarshal([]byte(body), &prd1)
	
	reqbody := "{\"title_register\": \"" + 
		prd1.Title + 
		"\",\"author_register\": \"" +
		 prd1.Author + 
		 "\",\"image_URL_register\": \"" +
		 prd1.Image_URL + 
		 "\",\"price_register\": \"" +
		 prd1.Price + 
		 "\",\"quantity_register\": \"" +
		 value + 
		 "\"}"
	req, _  := http.NewRequest("PUT", c.Endpoint + "/buckets/products/keys/"+key+"?pw=2&returnbody=true", strings.NewReader(reqbody) )
	req.Header.Add("Content-Type", "application/json")
	//fmt.Println( req )
	resp, err := c.Do(req)	
	if err != nil {
		fmt.Println("[RIAK DEBUG] " + err.Error())
		return prd_nil, err
	}	
	defer resp.Body.Close()
	body1, err := ioutil.ReadAll(resp.Body)
	if debug { fmt.Println("[RIAK DEBUG] PUT: " + c.Endpoint + "/buckets/products/keys/"+key+"?pw=2&returnbody=true => " + string(body1)) }
	if (strings.Contains(string(body1), "PW-value unsatisfied: ")){
		fmt.Println("Inside error")
		err := errors.New("PW-value unsatisfied")
		return prd_nil, err
	}
	var prd = product { }
	if err := json.Unmarshal(body1, &prd); err != nil {
		fmt.Println("[RIAK DEBUG] JSON unmarshaling failed: %s", err)
		return prd_nil, err
	}
	return prd, nil
}

// API Routes
func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/ping", pingHandler(formatter)).Methods("GET")
	mx.HandleFunc("/addproduct/{id}", NewProductHandler(formatter)).Methods("POST")
	mx.HandleFunc("/products", productGet(formatter)).Methods("GET")
	mx.HandleFunc("/products/{id}", productGet(formatter)).Methods("GET")
	mx.HandleFunc("/products/{id}", productUpdateQuantityHandler(formatter)).Methods("POST")
}

// Helper Functions
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"API version 1.0 alive!"})
	}
}

// API Add New Product to DB
func NewProductHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		
		params := mux.Vars(req)
		var uuid string = params["id"]
		fmt.Println( "Product Params ID: ", uuid )
		var prd product
    	_ = json.NewDecoder(req.Body).Decode(&prd)		
    	
		if uuid == ""  {
			formatter.JSON(w, http.StatusBadRequest, "Invalid Request. Product ID Missing.")
		} else {
			c1 := NewClient(server1)
			prd1, err1 := c1.AddProduct(uuid, prd)
			
			if (err1.Error() == "PW-value unsatisfied") {
				c2 := NewClient(server2)
				prd2, err2 := c2.AddProduct(uuid, prd)
				
				if err2 != nil {
					if (err2.Error() == "PW-value unsatisfied") {
						c3 := NewClient(server3)
						prd3, err3 := c3.AddProduct(uuid, prd)
						if err3 != nil {
							log.Fatal(err3)
							formatter.JSON(w, http.StatusBadRequest, err3)
						} else {
							prd = prd3
							formatter.JSON(w, http.StatusOK, prd)
						}
					}
					log.Fatal(err2)
					formatter.JSON(w, http.StatusBadRequest, err2)
				} else {
					prd = prd2
					formatter.JSON(w, http.StatusOK, prd)
				}
			} else if err1 != nil {
				log.Fatal(err1)
				formatter.JSON(w, http.StatusBadRequest, err1)
			} else {
				prd = prd1
				formatter.JSON(w, http.StatusOK, prd)
			}
		}
	}
}

// API Get product Status - Concurrent - Get One
func productGet(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		
		params := mux.Vars(req)
		var uuid string = params["id"]
		fmt.Println( "product Params ID: ", uuid )

		if uuid == ""  {
		
			c1 := make(chan []product)
    		c2 := make(chan []product)
    		c3 := make(chan []product)
    		//c4 := make(chan []product)
    		//c5 := make(chan []product)
    	
			go GetProductsServer1(c1)
			go GetProductsServer2(c2) 
			go GetProductsServer3(c3) 
			//go GetProductsServer4(c4) 
			//go GetProductsServer5(c5) 
			
			var prds []product
		  	select {
			    case prds = <-c1:
			        fmt.Println("Received Server1: ", prds)
			    case prds = <-c2:
			        fmt.Println("Received Server2: ", prds)
			        	/*
			    case prds = <-c3:
			        fmt.Println("Received Server3: ", prds)
			    case prds = <-c4:
			        fmt.Println("Received Server4: ", prds)
			    case prds = <-c5:
			        fmt.Println("Received Server5: ", prds)
			        */
		    }
		    
			if prds[0] == (product{}) {
				formatter.JSON(w, http.StatusBadRequest, "")
			} else {
				fmt.Println( "products: ", prds )
				formatter.JSON(w, http.StatusOK, prds)
			}
			
		} else {
		
			c1 := make(chan product)
	    	c2 := make(chan product)
	    	c3 := make(chan product)
			//c4 := make(chan product)
			//c5 := make(chan product)
	
			go GetProductServer1(uuid, c1) 
			go GetProductServer2(uuid, c2) 
			go GetProductServer3(uuid, c3) 
			//go GetProductServer4(uuid, c4) 
			//go GetProductServer5(uuid, c5) 
			
			var prd, prd1, prd2, prd3 product
			
			for i := 0; i < 3; i++ {
			  	select {
				    case prd1 = <-c1:
				        fmt.Println("Received Server1: ", prd1)
				    case prd2 = <-c2:
				        fmt.Println("Received Server2: ", prd2)
				    case prd3 = <-c3:
				        fmt.Println("Received Server3: ", prd3)
			    }
		    }
			
			/*
		  	select {
			    case prd = <-c1:
			        fmt.Println("Received Server1: ", prd)
			    case prd = <-c2:
			        fmt.Println("Received Server2: ", prd)
			    case prd = <-c3:
			        fmt.Println("Received Server3: ", prd)
			    case prd = <-c4:
			        fmt.Println("Received Server4: ", prd)
			    case prd = <-c5:
			        fmt.Println("Received Server5: ", prd)
			    
		    }
		    */

			if prd1.Quantity == prd2.Quantity{
		    	prd = prd1
		    } else if prd1.Quantity == prd3.Quantity {
		    	prd = prd1
		    } else if prd2.Quantity == prd3.Quantity {
		    	prd = prd2
		    }
		    
			if prd == (product{}) {
				formatter.JSON(w, http.StatusBadRequest, "")
			} else {
				fmt.Println( "product: ", prd )
				formatter.JSON(w, http.StatusOK, prd)
			}
		}
	}
}

func GetProductServer1(uuid string, chn chan<- product) {
	var prd_nil = product {}
	c := NewClient(server1)
	prd, err := c.GetProduct(uuid)
	if err != nil {
		chn <- prd_nil
	} else {
		fmt.Println( "Server1: ", prd)
		chn <- prd
	}
}

func GetProductServer2(uuid string, chn chan<- product) {
	var prd_nil = product {}
	c := NewClient(server2)
	prd, err := c.GetProduct(uuid)
	if err != nil {
		chn <- prd_nil
	} else {
		fmt.Println( "Server2: ", prd)
		chn <- prd
	}
}

func GetProductServer3(uuid string, chn chan<- product) {
	var prd_nil = product {}
	c := NewClient(server3)
	prd, err := c.GetProduct(uuid)
	if err != nil {
		chn <- prd_nil
	} else {
		fmt.Println( "Server3: ", prd)
		chn <- prd
	}
}

/*
func GetProductServer4(uuid string, chn chan<- product) {
	var prd_nil = product {}
	c := NewClient(server4)
	prd, err := c.GetProduct(uuid)
	if err != nil {
		chn <- prd_nil
	} else {
		fmt.Println( "Server4: ", prd)
		chn <- prd
	}
}

func GetProductServer5(uuid string, chn chan<- product) {
	var prd_nil = product {}
	c := NewClient(server5)
	prd, err := c.GetProduct(uuid)
	if err != nil {
		chn <- prd_nil
	} else {
		fmt.Println( "Server5: ", prd)
		chn <- prd
	}
}
*/

func GetProductsServer1(chn chan<- []product) {
	
	var prds_nil []product
	c := NewClient(server1)
	prds, err := c.GetProducts()
	if err != nil {
		chn <- prds_nil
	} else {
		fmt.Println( "Server1: ", prds)
		chn <- prds
	}
}

func GetProductsServer2(chn chan<- []product) {
	
	var prds_nil []product
	c := NewClient(server2)
	prds, err := c.GetProducts()
	if err != nil {
		chn <- prds_nil
	} else {
		fmt.Println( "Server2: ", prds)
		chn <- prds
	}
}

func GetProductsServer3(chn chan<- []product) {
	
	var prds_nil []product
	c := NewClient(server3)
	prds, err := c.GetProducts()
	if err != nil {
		chn <- prds_nil
	} else {
		fmt.Println( "Server3: ", prds)
		chn <- prds
	}
}

/*
func GetProductsServer4(chn chan<- []product) {
	
	var prds_nil []product
	c := NewClient(server4)
	prds, err := c.GetProducts()
	if err != nil {
		chn <- prds_nil
	} else {
		fmt.Println( "Server4: ", prds)
		chn <- prds
	}
}

func GetProductsServer5(chn chan<- []product) {
	
	var prds_nil []product
	c := NewClient(server5)
	prds, err := c.GetProducts()
	if err != nil {
		chn <- prds_nil
	} else {
		fmt.Println( "Server5: ", prds)
		chn <- prds
	}
}
*/

// API Update Quantity
func productUpdateQuantityHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		
		params := mux.Vars(req)
		var uuid string = params["id"]
		fmt.Println( "Product Params ID: ", uuid )
		var prd product
    	_ = json.NewDecoder(req.Body).Decode(&prd)		
    	fmt.Println("Update Product Quantity: ", prd.Quantity)

		if uuid == ""  {
			formatter.JSON(w, http.StatusBadRequest, "Invalid Request. Product ID Missing.")
		} else {
			c1 := NewClient(server1)
			prd1, err1 := c1.Updateproduct(uuid, prd.Quantity)
			
			if (err1.Error() == "PW-value unsatisfied") {
				c2 := NewClient(server2)
				prd2, err2 := c2.Updateproduct(uuid, prd.Quantity)
				
				if err2 != nil {
					if (err2.Error() == "PW-value unsatisfied") {
						c3 := NewClient(server3)
						prd3, err3 := c3.Updateproduct(uuid, prd.Quantity)
						if err3 != nil {
							log.Fatal(err3)
							formatter.JSON(w, http.StatusBadRequest, err3)
						} else {
							prd = prd3
							formatter.JSON(w, http.StatusOK, prd)
						}
					}
					log.Fatal(err2)
					formatter.JSON(w, http.StatusBadRequest, err2)
				} else {
					prd = prd2
					formatter.JSON(w, http.StatusOK, prd)
				}
			} else if err1 != nil {
				log.Fatal(err1)
				formatter.JSON(w, http.StatusBadRequest, err1)
			} else {
				formatter.JSON(w, http.StatusOK, prd1)
			}
		}
	}
}