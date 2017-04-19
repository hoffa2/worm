goog.require('proto.FromClient');
goog.require('proto.ToClient');
var Grapher = (function () {
    function Grapher() {
        var self = this;
        var button = document.getElementById('addSegmentfoo');
        var that = this;
        button.onclick = function () { self.addSegment(); };
    }
    Grapher.prototype.send = function (message) {
        this.websocket.send(message.serializeBinary());
    };
    Grapher.prototype.connectWebsocket = function (address, port) {
        var self = this;
        this.websocket = new WebSocket('ws://' + address + ":" + port + '/ws');
        this.websocket.onopen = function () {
            var message = new proto.FromClient;
            message.setGettarget(true);
            self.send(message);
            console.log("Message sent");
        };
        this.websocket.onmessage = function (e) {
            var reader = new FileReader();
            reader.onload = function () {
                var ToClient = proto.ToClient.deserializeBinary(this.result);
                switch (ToClient.getMsgCase()) {
                    case proto.ToClient.MsgCase.ADDNODE:
                        ToClient.getNodeid();
                        break;
                    case proto.ToClient.MsgCase.TARGET:
                        console.log("Jeeeei");
                        self.targetSegments = ToClient.getTarget();
                        break;
                    default:
                        console.log("No handler for messageType: " + ToClient.getMsgCase());
                        break;
                }
            };
            reader.readAsArrayBuffer(e.data);
        };
    };
    Grapher.prototype.lol = function () {
    };
    Grapher.prototype.addSegment = function () {
        var self = this;
        var message = new proto.FromClient;
        message.setChangetarget(self.targetSegments);
        self.send(message);
    };
    Grapher.prototype.killNode = function () {
        var message = new proto.message.FromClient;
    };
    return Grapher;
}());
var NodeVisualizer = (function () {
    function NodeVisualizer() {
        var container = document.getElementById('mynetwork');
        var data = {
            nodes: dataset,
            edges: edges
        };
        var options = {
            width: '900px',
            height: '900px',
            layout: {
                improvedLayout: true
            },
            physics: {
                stabilization: false
            },
            nodes: {
                shape: 'dot',
                size: 15,
                font: {
                    size: 12,
                    color: '#ffffff'
                },
                borderWidth: 2
            },
            edges: {
                width: 2
            }
        };
        network = new vis.Network(container, data, options);
    }
    NodeVisualizer.prototype.addNode = function (node) {
        dataset.add({ id: node.data.ID, label: message.data.ID });
        nextid = message.data.ID + "next";
        previd = message.data.ID + "prev";
        edges.update({ id: nextid, from: message.data.ID, to: message.data.Next, arrows: 'to' });
        edges.update({ id: previd, from: message.data.ID, to: message.data.Prev, color: { color: "red" }, arrows: 'to' });
        if (nodata) {
            init();
            console.log("Initializing network");
            nodata = false;
        }
    };
    NodeVisualizer.prototype.removeNode = function (node) {
        dataset.remove({ id: message.data.ID });
    };
    NodeVisualizer.prototype.updateNode = function (node) {
        nextid = message.data.ID + "next";
        previd = message.data.ID + "prev";
        edges.update({ id: nextid, from: message.data.ID, to: message.data.Next });
        edges.update({ id: previd, from: message.data.ID, to: message.data.Prev });
        console.log(message.data.successors);
    };
    NodeVisualizer.prototype.stabilize = function () {
        network.stabilize();
    };
    NodeVisualizer.prototype.disble_physics = function () {
        network.setOptions({ physics: false });
    };
    return NodeVisualizer;
}());
