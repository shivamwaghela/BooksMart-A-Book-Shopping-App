
/**

Mighty Gumball, Inc.
Version 5.0

- Refactored Previous REST Client Approach to Transaction Based REST API
- (i.e. instead of the Scaffolded REST API based on Domain Object Annotation) 
- Handlebars Page Templates
- Client State Validation using HMAC Key-Based Hash

NodeJS-Enabled Standing Gumball
Model# M102988
Serial# 1234998871109

**/

var gettransactionids = "http://54.219.165.171:8098/types/maps/buckets/usertransactions/datatypes/";

// added in v3: handlebars
// https://www.npmjs.org/package/express3-handlebars
// npm install express3-handlebars

// added in v2: crypto
// crypto functions:  http://nodejs.org/api/crypto.html


var crypto = require('crypto');
var fs = require('fs');
   var express = require('express');
    var favicon = require('serve-favicon');
    var logger = require('morgan');
    var methodOverride = require('method-override');
    var session = require('express-session');
    var bodyParser = require('body-parser');
    var multer = require('multer');
    var errorHandler = require('errorhandler');var Client = require('node-rest-client').Client;

var app = express();

app.use(express.bodyParser());
app.use("/images", express.static(__dirname + '/images'));
handlebars  = require('express3-handlebars');
hbs = handlebars.create();
app.engine('handlebars', hbs.engine);
app.set('view engine', 'handlebars');

var request = require('request-promise');

console.log('start');
var page = function( req, res, state, ts, status ) {
console.log('in page');
    var result = new Object() ;
    console.log( state ) ;

var key = "kartik1"
var key_req="http://54.219.165.171:3000/addtransaction/{"+key+"}"
request({
    "method":"GET", 
    "uri": key_req,

    "json": true,
    "headers": {
        "user_register":"laxmi", "transactionid":"123"
    }
  }).then(console.log, console.log);
   
    var client = new Client();
            client.get( gettransactionids, 
                function(data, response_raw){
                    console.log("Response begins");
                    console.log(data.value);
                    console.log("Response ends");
                   // console.log(jsdata);
        
                    var msg =   "hello" ;
                   

                 
            });
}

var handle_get = function (req, res, next) {
    console.log( "Get: ..." ) ;
    ts = new Date().getTime()
    console.log( ts )
    state = "no-coin" ;
    page( req, res, state, ts ) ;
}

app.get('/', handle_get ) ;



console.log( "Server running on Port 8080..." ) ;

app.listen(8080);

