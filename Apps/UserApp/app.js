var express = require('express');
var routes = require('./routes');
var user = require('./routes/user');
var http = require('http');
var path = require('path');
var home = require('./routes/home');

var app = express();

// all environments
app.set('port', process.env.PORT || 3006);
app.set('views', path.join(__dirname, 'views'));
app.set('view engine', 'ejs');
app.use(express.favicon());
app.use(express.logger('dev'));
app.use(express.json());
app.use(express.urlencoded());
app.use(express.methodOverride());
app.use(app.router);
app.use(express.static(path.join(__dirname, 'public')));
//app.use(express.bodyParser());
app.use(express.cookieParser());

// development only
if ('development' == app.get('env')) {
  app.use(express.errorHandler());
}

app.get('/', home.signup);

app.get('/signup.ejs', home.signup);
app.post('/afterSignUp', home.afterSignUp);
app.get('/getAllUsers', home.getAllUsers);
app.get('/getAllUsers', home.getAllUsers);
app.post('/login', home.login);
app.get('/login', home.login);
http.createServer(app).listen(app.get('port'), function(){
  console.log('Express server listening on port ' + app.get('port'));
});
