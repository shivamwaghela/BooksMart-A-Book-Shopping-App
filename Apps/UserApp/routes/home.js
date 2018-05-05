var axios = require('axios');
const ROOT_URL = 'http://localhost:3000';
var ejs = require("ejs");
var mysql = require('./mysql');
//var id;

function signup(req,res) {

	ejs.renderFile('./views/signup.ejs',function(err, result) {
	   // render on success
		if (!err) {
			//console.log(result);
	            res.end(result);
	   }
	   // render or error
	   else {
	            res.end('An error occurred');
	            console.log(err);
	   }
   });
}


function afterSignUp(req,res)
{
	
	var jsonToSend = {
		"UserId": req.param("inputUsername"),
		"Name":  req.param("name"),
		"Email": req.param("inputPassword")
	};
	console.log(jsonToSend);
	axios.post(`${ROOT_URL}/user`, jsonToSend).then(response => {
		console.log('response', response.data.UserId);
		//response.render('login', { title:response.data.UserId ,accountCreated: 'Account created successfully, please login!' },errorHandle);	
		if (response.data.UserId = response.data.UserId)
			{
			console.log("hello world");
			}	
	 
	 
		// dispatch(receivedImages(response));
			//console.log("response ----" + response);	
			//console.log("response ----" + response.UserId);	
			//console.log("response ----" + response.Name);
			//console.log("response ----"+response.Email);
	}).catch(error => {
		console.log("send error");
	});
}

function login(req, res) {
	
	
	var id = parseInt(req.param("inputUsername"));
	/*
	var jsonToSend = {
		'UserId': req.param("inputUsername"),
		//'Name':  req.param("name"),
		//'Email': req.param("inputPassword")
	}; */
	axios.get(`${ROOT_URL}/user/${id}`).then(response => {
		// dispatch(receivedImages(response));
		//console.log("response ----" + response);	
		//console.log("response ----" + response.UserId);	
		//console.log("response ----" + response.Name);
		//console.log("response ----"+response.Email);
		//console.log("response -------" + res);
		 //response.redirect('/signup');
		
	}).catch(error => {
		console.log("send error");
	});
		
}


function getAllUsers(req,res)
{
	var getAllUsers = "select * from users";
	console.log("Query is:"+getAllUsers);

	mysql.fetchData(function(err,results){
		if(err){
			throw err;
		}
		else
		{
			if(results.length > 0){

				var rows = results;
				var jsonString = JSON.stringify(results);
				var jsonParse = JSON.parse(jsonString);

				console.log("Results Type: "+(typeof results));
				console.log("Result Element Type:"+(typeof rows[0].emailid));
				console.log("Results Stringify Type:"+(typeof jsonString));
				console.log("Results Parse Type:"+(typeof jsString));

				console.log("Results: "+(results));
				console.log("Result Element:"+(rows[0].emailid));
				console.log("Results Stringify:"+(jsonString));
				console.log("Results Parse:"+(jsonParse));

				ejs.renderFile('./views/successLogin.ejs',{data:jsonParse},function(err, result) {
			        // render on success
			        if (!err) {
			            res.end(result);
			        }
			        // render or error
			        else {
			            res.end('An error occurred');
			            console.log(err);
			        }
			    });
			}
			else {

				console.log("No users found in database");
				ejs.renderFile('./views/failLogin.ejs',function(err, result) {
			        // render on success
			        if (!err) {
			            res.end(result);
			        }
			        // render or error
			        else {
			            res.end('An error occurred');
			            console.log(err);
			        }
			    });
			}
		}
	},getAllUsers);
}


exports.signup=signup;
exports.afterSignUp=afterSignUp;
exports.getAllUsers = getAllUsers;
exports.login=login;
