/**
		Payment API - Route Handler Module
		Author: Luis Otero
		
		Handles POST /transaction (using a Payment object to gather request)
		Handles GET /transactions (get all transactions)
		Handles GET /transactions/{id} (get transaction based on an ID)
		Handles POST /process (processes all transactions not processed yet)
		Handles PUT /update (update a payment based on an ID and amount)
		Handles DELETE /delete/{id} (deletes a payment transaction based on an ID)
**/

package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"encoding/json"
	"github.com/satori/go.uuid"
)

//JSON Render machine
var jsonRender = render.New(render.Options{ IndentJSON: true,})

// RabbitMQ queue for sending transactions to be processed
var queue_server = "192.168.99.100" 
var queue_port = "5672"
var queue_name = "payments"
var queue_user = "guest"
var queue_pass = "guest"

//Riak database details
var addresses = []string{"52.53.198.242:8087", "54.153.119.255:8087", "54.183.13.2:8087", "54.193.78.32:8087", "54.183.79.83:8087"}
var bucket_name = "payments"

/**
	POST /transactions
	
	1. Takes a payment
	2. Transforms it to a transaction
	3. Stores it into the database
	4. Pushes it to a queue for processing later
	5. Return a status object to the client
**/
func AddTransactionToQueue() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		
		var newPayment Payment;
		
		//decode payment request
		var decode_err = json.NewDecoder(request.Body).Decode(&newPayment)
		
		//check if payment parameters and object are valid
		if decode_err != nil {
			jsonRender.JSON(writer, 500, "Could not decode the payment request")
			return
		}
		
		//create a new transaction object with payment parameters
		var newTransaction = Transaction {
			TransactionId: uuid.NewV4().String(),
			PaymentType: newPayment.PaymentType,
			UserDetails: User {
				Name: newPayment.Name,
				Id: newPayment.UsernameId,
				Password: newPayment.Password,
			},
			Amount: FloatToString(newPayment.Amount),
			Status: "Payment Pending",
		}
		
		//connect to database
		database,db_error := RiakConnect(addresses)
		
		if db_error != nil {
			jsonRender.JSON(writer, 600, "Could Not Connect to Database")
			return
		}
		
		//connect to queue
		_,channel,queue_error := RabbitmqConnect(queue_server, queue_port, queue_name, queue_user, queue_pass)
		
		if queue_error != nil {
			jsonRender.JSON(writer, 600, "Could Not Connect to Queue")
			return
		}
		
		//insert transaction to database for record keeping
		RiakSet(database, bucket_name, newTransaction.TransactionId, newTransaction)
		
		//send message to queue
		Enqueue(channel, queue_name, newTransaction.TransactionId)
		
		type TransactionState struct {
			Id string
			Status string
		}
		
		order := TransactionState{ Id : newTransaction.TransactionId, Status : "Payment Placement Successful!"}
		
		//send a status object to the client
		jsonRender.JSON(writer, http.StatusOK, order)
	}
}

/**
	GET /transactions and /transactions/{id}
	
	1. Checks for which of the two get operations to perform
	2. Get the transaction/s from the database
	3. Return the full transaction to the client
**/
func SearchForTransaction() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		parameters := mux.Vars(request)
		
		//connect to database
		database, db_connect_error := RiakConnect(addresses)
		
		if db_connect_error != nil {
			jsonRender.JSON(writer, 600, db_connect_error.Error())
			return
		}
		
		//Option 1: get all transactions. Option 2: get transaction based on ID specified
		if parameters["id"] == "" {
			
			if set,error := RiakGetAll(database, bucket_name); error == nil {
				for _,value := range set {
					jsonRender.JSON(writer, http.StatusOK, value)
				}
			}else {
				jsonRender.JSON(writer, 500, error.Error())
			}
			
		}else {
			if transaction,err := RiakGet(database, bucket_name, parameters["id"]); err == nil {
				jsonRender.JSON(writer, http.StatusOK, transaction)
				return
			}else {
				jsonRender.JSON(writer, 500, err.Error())
				return
			}	
		}
		

	}
}

/**
	POST /process
	
	1. Dequeue the transactions
	2. Get all transactions from database
	3. Process transactions from queue
	4. Update values in the history
	5. Add unprocessed transactions to queue
	6. Return processed transactions to client
**/
func ProcessAllTransactions() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		//connect to queue
		_,channel,queue_connect_error := RabbitmqConnect(queue_server, queue_port, queue_name, queue_user, queue_pass)
		
		if queue_connect_error != nil {
			jsonRender.JSON(writer, 600, queue_connect_error.Error())
			return
		}
		
		//dequeue all transactions from queue
		ids,dequeue_error := DequeueAll(channel, queue_name)
		
		if dequeue_error != nil {
			jsonRender.JSON(writer, 600, dequeue_error.Error())
			return
		}
		
		//connect to database
		database,db_error := RiakConnect(addresses)
		
		if db_error != nil {
			jsonRender.JSON(writer, 600, db_error.Error())
			return
		}
		
		//process all of the transactions from the queue
		for _,value := range ids {
			if transaction,get_error := RiakGet(database, bucket_name, value); get_error == nil {
				if transaction.Status != "Payment Process Success!" {
					transaction.Status = "Payment Process Success!"
					transaction.UserDetails.Id = ""
					transaction.UserDetails.Password = ""
					jsonRender.JSON(writer, http.StatusOK, transaction)
					RiakSet(database, bucket_name, transaction.TransactionId, transaction)
				}
			}
		}
		
		//get all of the transactions from the database
		set,get_all_error := RiakGetAll(database, bucket_name)
		
		if get_all_error != nil {
			jsonRender.JSON(writer, 600, get_all_error.Error())
			return
		}
		
		//push transactions to queue that have not been processed
		for _,value := range set {
			if value.Status == "Payment Pending" {
				push_error := Enqueue(channel, queue_name, value.TransactionId)
				
				if push_error != nil {
					jsonRender.JSON(writer, 600, push_error.Error())
					return
				}
			}
		}
		
	}

}

/**
	PUT /update
	
	1. Create update and status objects
	2. Attempt to decode update request
	3. Get the transaction from the database
	4. Check if the transaction can be updated
	5. Update the transaction
	6. Store the transaction in the database
	7. Return a update status to the client
**/
func UpdateTransaction() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		//update object struct
		type update_details struct {
			Id string
			Amount float64
		}
		
		//status object struct
		type update_status struct {
			Id string
			Status string
		}
		
		var new_update update_details
		
		//decode update request
		parse_error := json.NewDecoder(request.Body).Decode(&new_update)
		
		if parse_error != nil {
			jsonRender.JSON(writer, 500, parse_error.Error())
			return
		}
		
		//connect to the database
		database,db_error := RiakConnect(addresses)
		
		if db_error != nil {
			jsonRender.JSON(writer, 600, db_error.Error())
			return
		}
		
		//Get transaction from the database
		newTransaction, get_error := RiakGet(database, bucket_name, new_update.Id)
		
		if get_error != nil {
			jsonRender.JSON(writer, 500, get_error.Error())
			return
		}
		
		//Update transaction and store
		if newTransaction.Status == "Payment Pending" {
			newTransaction.Amount = FloatToString(new_update.Amount)
			store_error := RiakSet(database, bucket_name, newTransaction.TransactionId, newTransaction)
			
			if store_error != nil {
				jsonRender.JSON(writer, 600, store_error.Error())
				return
			}
			
			//send status to client
			jsonRender.JSON(writer, http.StatusOK, update_status{ Id : new_update.Id, Status : "Payment update successful", })
		}else {
			jsonRender.JSON(writer, 500, update_status{ Id : new_update.Id, Status : "Payment already processed. Cannot update payment amount"})
		}		
		
	}
}

/**
	DELETE /delete/{id}
	
	1. Checks if an id was given
	2. Delete transaction from database
	3. Return status of deletion to client
**/
func DeleteTransaction() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		parameters := mux.Vars(request)
		
		//connect to the databse
		database,connect_err := RiakConnect(addresses)
		
		if connect_err != nil {
			jsonRender.JSON(writer, 600, "Could not connect to the database")
			return
		}
		
		//check if an ID has been provided
		if parameters["id"] == "" {
				jsonRender.JSON(writer, 404, "Must have a key id to delete. /delete/{id}")
			
		}else {
			if error := RiakDelete(database, bucket_name, parameters["id"]); error == nil {
				//return status to client
				jsonRender.JSON(writer, http.StatusOK, parameters["id"] + " deleted!")
			}else {
				jsonRender.JSON(writer, 500, error.Error())
			}
			
		}
		

	}
}

/**
	GET /ping and /
	
	1. Returns a json status of the API to the client
**/
func PingResponseProvider() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		type ping struct{
			API string
			Status string
		}
		
		jsonRender.JSON(writer, http.StatusOK, ping{ API: "Payment API v1", Status: "Server is running on port 4500",})
	}
}