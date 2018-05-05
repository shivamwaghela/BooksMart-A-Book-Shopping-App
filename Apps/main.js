var cart = require("./cartApp/app")
var payment = require("./paymentApp/app")

cart.use(payment);
cart.set('port', process.env.PORT || '4000')
cart.listen(cart.get('port'), function() {
  
});

/*
payment.set('port', process.env.PORT || '10000')
payment.listen(payment.get('port'), function() {
	
});
*/

