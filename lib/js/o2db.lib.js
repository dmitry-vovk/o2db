var O2DB = (function () {
    'use strict';
    var client = function () {
        var socket;
        this.init.apply(this, arguments);
    };
    client.prototype.send = function (param) {
        console.log('sending type ' + param.type);
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
    client.prototype.init = function (options) {
        var self = this;
        this.socket = new WebSocket(options.host);
        this.socket.onclose = function (evt) {
            if (typeof options.onclose === 'function') {
                options.onclose(evt);
            }
        };
        this.socket.onopen = function (evt) {
            if (typeof options.onopen === 'function') {
                options.onopen.apply(self, evt);
            }
        };
        this.socket.onmessage = function (evt) {
            if (typeof options.onmessage === 'function') {
                options.onmessage(JSON.parse(evt.data));
            }
        };
        this.socket.onerror = function (evt) {
            if (typeof options.onerror === 'function') {
                options.onerror.apply(this, evt);
            }
        };
    };
    return client;
})();


var params = {
    host: 'ws://localhost:8088/',
    onopen: function (e) {
        var o = {
            'type': 0,
            'payload': {
                'name': 'root',
                'password': '12345'
            }
        };
        this.send(o);
        this.send({
            'type': 102,
            'payload': {
                'name': 'test_01'
            }
        });
        this.send({
            'type': 403,
            'payload': {
                'classes': ['Job']
            }
        });
    },
    onclose: function (e) { console.log(e) },
    onmessage: function (data) {
        console.log('Got message:');
        console.dir(data.response);
    }
};

var client = new O2DB(params);
