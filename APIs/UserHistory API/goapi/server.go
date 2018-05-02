package main

import (
"fmt"
"log"
"net/http"
"io/ioutil"
"time"
//	"strings"
//"encoding/json"
"github.com/codegangsta/negroni"
"github.com/gorilla/mux"
"github.com/unrolled/render"
//"github.com/satori/go.uuid"

"strings"
"encoding/json"
)


var debug=true
var server1 = "http://54.219.165.171:8098"
//var server2 = "http://52.52.124.207:8098"
//var server3 = "http://54.176.35.69:8098"
//var server4 = "http://54.219.105.255:8098"
//var server5 = "http://54.241.239.217:8098"

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

	/*c2 := NewClient(server2)
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
	}*/

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

func (c *Client) AddUserTransactions(key string, user_transaction UserTransaction) (UserTransaction, error) {

	var ut_nil= UserTransaction{}

	/*reqbody := "{\"update\": {\"transactionIds_set\": {\"add_all\": [\""+
		user_transaction.TransactionId+
		"\"]}}}"*/
	fmt.Println("AddUserTransactions:User Name is: "+user_transaction.UserName)

	fmt.Println("AddUserTransactions:User Transaction Id is: "+user_transaction.TransactionId)
	reqbody :=  "{\"update\": {\"user_register\": \""+
		key +
		"\",\"usertransactionids_set\": {\"add_all\": [\""+
		user_transaction.TransactionId+
		"\"]}}}"


	fmt.Println(reqbody + "key is "+key)
	req, _ := http.NewRequest("POST", c.Endpoint+"/types/maps/buckets/testmay9/datatypes/"+key+"?returnbody=true", strings.NewReader(reqbody))
	req.Header.Add("Content-Type", "application/json")
	fmt.Println("Request is: ");
	fmt.Println(req )
	fmt.Println("End of Request ");
	resp, err := c.Do(req)
	fmt.Println("1")
	fmt.Println(debug)

	if err != nil {
		fmt.Println("[RIAK DEBUG] " + err.Error())
		return ut_nil, err
	}
	fmt.Println("2")

	fmt.Println(debug)

	defer resp.Body.Close()
	fmt.Println("3")

	fmt.Println(debug)

	body, err := ioutil.ReadAll(resp.Body)
	if debug {
		fmt.Println("[RIAK DEBUG] PUT: " + c.Endpoint + "/types/maps/buckets/testmay9/datatypes/" + key + "?returnbody=true => " + string(body))
	}

	var ut= UserTransaction{}
	if err := json.Unmarshal(body, &ut); err != nil {
		fmt.Println("[RIAK DEBUG] JSON unmarshaling failed: %s", err)
		return ut_nil, err
	}
	return ut, nil
}







// API Routes
func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/ping", pingHandler(formatter)).Methods("GET")
	mx.HandleFunc("/addtransaction/{id}", addUserTransactionHandler(formatter)).Methods("POST")
	mx.HandleFunc("/getusertransactions/{id}", getUserTransactionsHandler(formatter)).Methods("GET")
	//mx.HandleFunc("/getTransactionDetailsHandler/{id}", GetTrasactionDetails(formatter)).Methods("GET")

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
//


func addUserTransactionHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		fmt.Println("In addUserTransactionHandler");

		params := mux.Vars(req)

		var uname string = params["id"]

		fmt.Println( "addUserTransactionHandler:UserHistory Params ID: ", uname )

		var utransaction UserTransaction
		_ = json.NewDecoder(req.Body).Decode(&utransaction)
		fmt.Println("User Name is: "+utransaction.UserName)

		fmt.Println("addUserTransactionHandler:User Transaction Id is: "+utransaction.TransactionId)
		if uname == ""  {
			formatter.JSON(w, http.StatusBadRequest, "Invalid Request. User ID Missing.")
		} else {
			c1 := NewClient(server1)
			prd, err := c1.AddUserTransactions(uname, utransaction)
			if err != nil {
				log.Fatal(err)
				formatter.JSON(w, http.StatusBadRequest, err)
			} else {
				formatter.JSON(w, http.StatusOK, prd)
			}
		}
		log.Println("end of user transaction")

		fmt.Println("End of addUserTransactionHandler");

	}



}



func getUserTransactionsHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("In  getUserTransactionsHandler")
		params := mux.Vars(req)
		var uuid string = params["id"]
		fmt.Println( "User Params ID: ", uuid )

		if uuid == "" {
			formatter.JSON(w, http.StatusBadRequest, "Invalid Request. User ID Missing.")
		} else {

			c := NewClient(server1)

			transactions:=AllUserTransactions{}
			transactions = c.GetTransactionIds(uuid)
			fmt.Println("After  GetTransactionIds")
			fmt.Println("Your transactions are here: ", transactions)


			if transactions.Dtype  == "" {
				formatter.JSON(w, http.StatusBadRequest, "")
			} else {
				fmt.Println("Your transactions are in statusok: ", transactions)
				formatter.JSON(w, http.StatusOK ,transactions)
				fmt.Println(&w)
				fmt.Println(w)
p:=&w
fmt.Println(*p)

			}
		}
		fmt.Println("End of  getUserTransactionsHandler")

	}
}

func (c *Client) GetTransactionIds(key string) AllUserTransactions {
	fmt.Println("In  GetTransactionIds")

	var tid_nil = AllUserTransactions{}
	resp, err := c.Get(c.Endpoint + "/types/maps/buckets/testmay9/datatypes/" + key + "?include_context=false")

	if err != nil {
		fmt.Println("[RIAK DEBUG] " + err.Error())
		return tid_nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("Start of get user body")
	fmt.Println(string(body))
	fmt.Println("End of get user body")

	if debug {
		fmt.Println("[RIAK DEBUG] GET: " + c.Endpoint + "/types/maps/buckets/testmay9/datatypes/" + key + " => " + string(body))
	}

	var tid = AllUserTransactions{}
	if err := json.Unmarshal((body), &tid); err != nil {
		fmt.Println("RIAK DEBUG] JSON unmarshaling failed: %s", err)
		return tid_nil
	}
	fmt.Println("After Json unmarshall")
	fmt.Println("1:" + tid.Dtype)
	//fmt.Println( &tid.Value.TransactionIds)
	//fmt.Println( tid.Value.TransactionIds)

	fmt.Println("End of  GetTransactionIds")
	fmt.Println("begin of try")

	fmt.Println(tid.Value)
	fmt.Println(tid)







	fmt.Println("end of try")

	return tid
}

