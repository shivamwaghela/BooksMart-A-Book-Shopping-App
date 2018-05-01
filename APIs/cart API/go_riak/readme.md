Welcome to the CMPE281 Team Starburst Hackathon Project wiki!

## Cart API

* cart is stored as nested struct containing userid, cartdid(id) and struct of items, containing Items.Name, Items.Count, Items.Rate and Items.Amount. It also stores the total of the cart in Total
* All the 5 IP addresses of the nodes are given in var. This will be changed to addresses of 2 different load balancers in 2 VPC's. 
* /ping used to check status of nodes
* /order handles new order request, and order added to cart 
* /view/{id} pass the order id, and get the order details.
* /history/{id} gets the keys of all objects stored
* /update updates an element in cart 
* /clearcart clears the cart, clears riak 

Not yet completly implemented. 
