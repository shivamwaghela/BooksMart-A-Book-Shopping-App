
/**

Book Shopping Cart
Team: Starburst
CMPE 281

**/

var userAPI = ""
var paymentAPI = "localhost:4500"
var productsAPI = "http://localhost:3000"
var shoppingAPI = ""
var historyAPI = ""
var vendorsAPI = ""


var fs = require('fs');
var express = require('express');
var Client = require('node-rest-client').Client;
var alert = require('alert-node');
var session = require('express-session');

var app = express();
app.use(express.bodyParser());
app.use(session({ secret: 'token'}));

var transactions = Array();

//For product catalog
var products = [];
var XMLHttpRequest = require("xmlhttprequest").XMLHttpRequest;

/* Views */

var checkoutView = function(req,res,data,method) {
	var page_body = fs.readFileSync('./checkout.html');
	if (method == "post") {
		if (data.transactionStatus)
		{
			transactions.push(data.Id);
			alert("Transaction " + data.Id + ": " + data.Status);
			res.redirect("/");
			return
		}
		else
		{
			alert("Payment did not process");
			return
		}
	}
	res.setHeader('Content-Type', 'text/html');
	res.writeHead(200);
	var body = "" + page_body;
	
	var body = body.replace("{amount}", data.amount.toString());
	var body = body.replace("{username}", data.user);
	res.end(body);
}

var loginView = function(req,res,data,method) {
	var body = fs.readFileSync('./login.html');
	var login_body = "" + body;
	
	if(method == "post") {
		req.session.username = "Frank";
		res.redirect("/payment");
		return
	}
	res.setHeader('Content-Type', 'text/html');
	res.writeHead(200);
	
	res.end(login_body)
}

var productsView = function(req,res,data) {
	//TODO
}

var cartView = function(req,res,data) {
	//TODO
}

var userHistoryView = function(req,res,data) {
	//TODO
}


/* GET HANDLERS */

var get_checkout_page = function(req,res) {
	var data = Object();
	data.amount = 55.67;
	data.user = req.session.username;
	checkoutView(req,res,data,"get");
}

var get_products_page = function(request,response) {
	console.log("In GET Products ");
	var xmlhttp = new XMLHttpRequest();
	xmlhttp.open("GET", productsAPI+ "/products");  
	xmlhttp.setRequestHeader("Content-Type", "application/json");
	xmlhttp.send();  
	xmlhttp.onreadystatechange = function() {
		if (this.readyState === 4 && this.status === 200) {
			console.log("API call to product API successful. Status: " + this.status);
			console.log (this.responseText);
			var products_array = JSON.parse(this.responseText);
			console.log("products array")
			console.log(products_array);
			response.sendFile('/pages/product_catalog.html', {products: products_array});
		}
	}
}

var get_login_page = function(req,res) {
	//TODO
	
	loginView(req,res,"", "get"); //test
}

var get_cart_page = function(req,res) {
	//TODO
}

var get_user_history_page = function(req,res) {
	//TODO
}

/* POST HANDLERS */

var submit_payment = function(req,res) {
	var data = Object();
	
	/*
	var client = Client();
	var url = paymentAPI + "/transaction";
	var args = {
		data: { "PaymentType": res.body.type, 
			"Name": res.body.card_name, 
			"UsernameId": res.body.u_name, 
			"Password": res.body.pay_pass, 
			"Amount": res.body.amount
		},
		headers: { "Content-Type": "application/json"}
	};
	var send_post = client.post(url,args,function(data,response_raw){
		var response_data = JSON.parse(data)
		data.user = req.body.user;
		data.Id = req.body.u_name;
		data.Status = response_data.Status;
		data.transactionStatus = true;
		checkoutView(req,res,data,"post");
	});
	
	send_post.on('error', function(err) {
		data.transactionStatus = false;
		checkoutView(req,res,data,"post")
	});
	*/
	
	data.user = req.body.user;
	data.Id = req.body.u_name;
	data.Status = req.body.pay_pass;
	data.transactionStatus = true;
	checkoutView(req,res,data,"post");
}

var login_post = function (req, res) {
	var data = Object();
	data.username = req.body.username;
	data.password = req.body.password;
	data.user = "Frank";
	loginView(req,res,data,"post");
}

var products_post = function (req,res) {
	//TODO
}

var cart_post = function (req,res) {
	//TODO
}

var history_post = function (req,res) {
	//TODO
}

/* Routes */

app.post("/", login_post );
app.get("/", get_login_page );
app.get("/products", get_products_page);
app.post("/products", products_post);
app.get("/cart", get_cart_page);
app.post("/cart", cart_post);
app.get("/history", get_user_history_page);
app.post("/history", history_post);
app.get("/payment", get_checkout_page);
app.post("/payment", submit_payment);


/* Server Start */

app.set('port', (process.env.PORT || 5000));
app.listen(app.get('port'), function() {
  console.log('Books shopping Cart app is running on port', app.get('port'));
});


