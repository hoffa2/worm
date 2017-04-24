goog.require('proto.FromClient')
goog.require('proto.ToClient')

class Grapher {

    private websocket: any;

    private targetSegments: any;

    private nodeviz: NodeVisualizer;

    private send (message) {
        this.websocket.send(message.serializeBinary());
    }

    constructor() {
        let self = this;

        var button = document.getElementById('addSegmentfoo');
        let that = this;
        button.onclick = function() {self.addSegment()}
    }

    public connectWebsocket(address: string, port: number) {
        let self = this;

        this.websocket = new WebSocket('ws://' + address + ":" + port + '/ws');

        this.websocket.onopen = function () {
            var message = new proto.FromClient;
            message.setGettarget(true);
            self.send(message);
            console.log("Message sent")
        }

        this.websocket.onmessage = function (e) {

            let reader = new FileReader();

            reader.onload = function ()
            {
                var ToClient = proto.ToClient.deserializeBinary(this.result)

                switch (ToClient.getMsgCase())
                {
                    case proto.ToClient.MsgCase.ADDNODE:
                        self.nodeviz.addNode(ToClient.getNodeid());
                        break;
                    case proto.ToClient.MsgCase.TARGET:
                        console.log("Jeeeei")
                        self.targetSegments = ToClient.getTarget();
                        break;
                    default:
                        console.log("No handler for messageType: " + ToClient.getMsgCase())
                        break;
                }
            }
            reader.readAsArrayBuffer(e.data);
        }

    }

    public addSegment() {
        let self = this;
        var message = new proto.FromClient;
        message.setChangetarget(self.targetSegments);
        self.send(message);
    }

    public killNode() {
        let message = new proto.message.FromClient;
    }
}

class NodeVisualizer {
        private dataset: any;
        private network: any;
        private nodes: any;
        private edges: any;

        public addNode(nodeid) {
            this.dataset.add({id: nodeid, label: nodeid});
            //edges.update({id: nextid, from: message.data.ID, to: message.data.Next, arrows:'to'});
            //edges.update({id: previd, from: message.data.ID, to: message.data.Prev, color:{color:"red"}, arrows:'to'});
        }

        public removeNode(node) {
            this.dataset.remove({id: node});
        }

        public updateNode(node) {
            //nextid = message.data.ID + "next"
            //previd = message.data.ID + "prev"
            //edges.update({id: nextid, from: message.data.ID, to: message.data.Next});
            //edges.update({id: previd, from: message.data.ID, to: message.data.Prev})
            //console.log(message.data.successors)
        }

        public stabilize() {
            this.network.stabilize();
        }

        public disble_physics() {
            this.network.setOptions({physics: false});
        }

        constructor() {
            var container = document.getElementById('mynetwork')
                var data = {
                    nodes: this.dataset,
                    edges: this.edges
                };
            var options = {
                width: '900px',
                height: '900px',
                layout: {
                    improvedLayout: true,
                },
                physics: {
                    stabilization: false,
                },
                nodes: {
                    shape: 'dot',
                    size: 15,
                    font: {
                        size: 12,
                        color: '#ffffff'
                    },
                    borderWidth: 2,
                },
                edges: {
                    width: 2
                }
            };
            this.network = new vis.Network(container, data, options)
        }
}

