

var http = require('http')
var url = require('url')
var fs = require('fs')
var port = process.env.PORT || 3000;
var querystring = require('querystring')
var index = fs.readFileSync('index.html')
const PORT=8080
var Promise = require('bluebird')
var dispatcher = require('httpdispatcher');
nodes = [];

function handleRequest(request, response){
    try {
        //log the request on console
         console.log(request.url);
        //Disptach
        dispatcher.dispatch(request, response);
    } catch(err) {
        console.log(err);
    }
}
var server = http.createServer(handleRequest);

server.listen(PORT, function(){
    //Callback triggered when server is successfully listening. Hurray!
        console.log("Server listening on: http://localhost:%s", PORT);
});

var io = require('socket.io').listen(server);

dispatcher.setStatic('resources');

//var everyone = nowjs.initialize(server)
//A sample GET request
dispatcher.onGet("/", function(req, res) {
    res.writeHead(200, {'Content-Type': 'text/html'});
    res.end(index);
});

//A sample POST request
dispatcher.onPost("/add", function(req, res) {
    var obj = JSON.parse( req.body );

    console.log(obj)
    io.emit('add', {data: obj})
    res.writeHead(200, {'Content-Type': 'text/plain'});
    res.end('Got Post Data');
});

dispatcher.onPost("/remove", function(req, res) {
    var obj = JSON.parse( req.body );

    io.emit('remove', {data: obj})
    res.writeHead(200, {'Content-Type': 'text/plain'});
    res.end()
});


dispatcher.onPost("/update", function(req, res) {
    var obj = JSON.parse( req.body );

    io.emit('update', {data: obj})
    res.writeHead(200, {'Content-Type': 'text/plain'});
    res.end()
});

