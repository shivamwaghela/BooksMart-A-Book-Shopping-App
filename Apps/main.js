var cart = require("./cartApp/app");
var payment = require("./paymentApp/app");
var users = require("./UserApp/app");
var history = require("./User History/app");
var products = require("./ProductsApp/app");

cart.use(payment);
cart.use(users);
cart.use(history);
cart.use(products);
cart.set('port', process.env.PORT || '4000');

cart.listen(cart.get('port'), function() {
  	console.log("Running on port 4000");
});

