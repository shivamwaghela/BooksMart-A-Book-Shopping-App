
/**

Payment Area of Shopping Cart
Team: Starburst
CMPE 281

**/

var paymentAPI = "http://localhost:4500";


var fs = require('fs');
var express = require('express');
var Client = require('node-rest-client').Client;
var alert = require('alert-node');
//var session = require('express-session');
//var payment = express.Router();
var bp = require('body-parser');

var app = express();
app.use(bp.json());
app.use(bp.urlencoded({ extended: false }));
//app.use(session({ secret: 'token'}));

var transactions = Array();

// Luis Relevant Code
var checkoutView = function(req,res,data,method) {
	var page_body = fs.readFileSync('./paymentApp/checkout.html');
	if (method == "post") {
		transactions.push(data.Id);
		alert("Transaction " + data.Id + ": " + data.Status);
		res.redirect("/");
		return
	}
	res.setHeader('Content-Type', 'text/html');
	res.writeHead(200);
	var body = "" + page_body;
	
	
	body = body.replace("{user}", "user1");
	body = body.replace(/{amount}/g, data.amount.toString());
	body = body.replace("{username}", data.user);
	return res.end(body);
}

var get_checkout_page = function(req,res,next) {
	var data = Object();
	//call cart
	data.amount = 55.67;
	data.user = "user1" //req.session.username;
	checkoutView(req,res,data,"get");
}

var submit_payment = function(req,res,next) {
	
	var data = Object();
	
	data.amount = parseFloat(req.body.amount)
	data.user = "user1"
	if (req.body.submit === "Cancel") {
		res.redirect("/");
		return
	}
	if(req.body.submit === "Logoff") {
		//req.session.destroy();
		res.redirect("/");
		return
	}
	
	
	if (req.body.card_name == "") {
		alert("Must input a name for the payment")
		checkoutView(req,res,data,"get")
		return
	}
	if (req.body.u_name == "") {
		alert("Must input a username")
		checkoutView(req,res,data,"get")
		return
	}
	if (req.body.pay_pass == "") {
		alert("Must input a password")
		checkoutView(req,res,data,"get")
		return
	}
	
	var client = new Client();
	var url = paymentAPI + "/transaction";
	var args = {
		data: { "PaymentType": req.body.pay_type, 
			"Name": req.body.card_name, 
			"UsernameId": req.body.u_name, 
			"Password": req.body.pay_pass, 
			"Amount": parseFloat(req.body.amount)
		},
		headers: { "Content-Type": "application/json"}
	};
	var send_post = client.post(url,args,function(data,response){
		var processed_data = JSON.parse(data)
		console.log(url + " request success!");
		processed_data.user = req.body.user;
		//post user history
		//get cart items
		//clear cart
		checkoutView(req,res,processed_data,"post");
	});
	
	send_post.on('error', function(err) {
		console.log(err.toString())
		alert("Payment did not process");
	});
	
}

/* routes */
var payment = express.Router();
payment.get('/payment', get_checkout_page);
payment.post('/payment', submit_payment);

module.exports = payment


/* Server Start */

//app.set('port', (process.env.PORT || 10000));
/*
app.listen(app.get('port'), function() {
	console.log('Multi-User shopping Cart app is running on port', app.get('port'));
});
*/


