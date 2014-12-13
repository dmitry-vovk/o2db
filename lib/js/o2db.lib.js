var O2DB = (function () {
    'use strict';
    var client = function () {
        var socket;
        this.init.apply(this, arguments);
    };
    client.prototype.send = function (param) {
        var sock = this.socket;
        this.waitSocket(this.socket, function () {
            sock.send(JSON.stringify(param));
        });
    };
    client.prototype.waitSocket = function (s, callback) {
        var self = this;
        setTimeout(function () {
            if (typeof s != 'undefined' && s.readyState === 1) {
                callback();
                return;
            } else {
                self.waitSocket(callback);
            }
        }, 5);
    };
    client.prototype.init = function (addr, listener) {
        this.socket = new WebSocket(addr);
        this.socket.onclose = function () {
            console.log('Connection closed');
        };
        this.socket.onopen = function () {
            console.log('Connection opened');
        };
        this.socket.onmessage = function (evt) {
            if (typeof listener === 'function') {
                listener(JSON.parse(evt.data))
            }
        };
    };
    return client;
})();

var client = new O2DB('ws://localhost:8088/', function (data) {
    console.log('Got message:');
    console.dir(data);
});

var o = {
    'type': 0,
    'payload': {
        'name': 'root',
        'password': '12345'
    }
};

client.send(o);
