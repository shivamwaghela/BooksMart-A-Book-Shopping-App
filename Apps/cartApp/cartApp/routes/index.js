var express = require('express');
var router = express.Router();
var axios = require('axios');
var waterfall = require('async-waterfall');
var orderID = "";

/* GET home page. */
//axios api, axios.get /axios.post
router.get('/', function(req, res, next) {
    var cartItems = [];
    /*var req = {
        "userId": "",
        "items": []
    };*/
    var req = {"userId":"user2",
    "items":[
        {
        "name":"starbucks",
        "count":1,
        "rate":3.95
        },
        {
        "name":"peets",
        "count":1,
        "rate":4.95
        }
        ]
    };
    waterfall([
        /*function(callback){
            axios.get('')
            .then(function(response){
                req.userId = response.userId;
            });
        },*/
        function(callback){
            console.log("in function 1");
            axios.post('http://127.0.0.1:3000/order', req)
            .then(function (res){
                //console.log('response',res.data);
                cartItems=res.data.items;
                orderID = res.data.id;
                console.log("in then", cartItems);
                callback(null, cartItems);
            });
            
        },
        function(cartItems, callback){
            console.log("in function 2");
            console.log(cartItems);
            res.render('cart', { items: cartItems });
        }
        ], function(){

        });
    /*
    axios.post('http://127.0.0.1:3000/order', req)
    .then(function (res){
        //console.log('response',res.data);
        cartItems=res.data.items;
        console.log("in then", cartItems);
    });
    console.log(cartItems);*/

  //res.render('cart', { items: cartItems });
});

router.get('/clearCart', function (req,res,next) {
  console.log("in delete cart");
  console.log("order id " + orderID);
  var req={
    id: orderID
  };
  waterfall([
    function(callback){
        console.log("in function 2");
        axios.delete('http://127.0.0.1:3000/clearCart', {data: {id: orderID}})
        .then(function(res){
            console.log("in then");
            callback(null);
        });
    },
    function(callback){
        console.log("in function 3");
        res.render('cart', { items: [] });
    }
    ]);
  //res.render('cart', { items: [] });
});

module.exports = router;
