var cart = require("./cartApp/app");
var payment = require("./paymentApp/app");
var users = require("./usersApp/app");
var history = require("./userHistory/app");
var products = require("./products/app");

cart.use(payment);
cart.use(users);
cart.use(history);
cart.use(products);
cart.set('port', process.env.PORT || '4000');
cart.listen(cart.get('port'), function() {
  	console.log("Running on port 4000");
});

/*
payment.set('port', process.env.PORT || '10000')
payment.listen(payment.get('port'), function() {
	
});
*/

