var O2DB = (function () {
    'use strict';
    var client = function () {
        var socket;
        this.init.apply(this, arguments);
    };
    var statusConst = {
        0: 'Authenticate',
        100: 'CreateDatabase',
        101: 'DropDatabase',
        102: 'OpenDatabase',
        103: 'ListDatabases',
        200: 'CreateCollection',
        201: 'DropCollection',
        202: 'ListCollections',
        300: 'ObjectGet',
        301: 'ObjectWrite',
        302: 'ObjectDelete',
        303: 'GetObjectVersions',
        304: 'GetObjectDiff',
        305: 'SelectObjects',
        400: 'Subscribe',
        401: 'AddSubscription',
        402: 'CancelSubscription',
        403: 'ListSubscriptions'
    };
    // Set object immutable properties
    for (var i in statusConst) {
        Object.defineProperty(client.prototype, statusConst[i], {
            writable: false,
            value: parseInt(i)
        });
    }
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
            'type': client.Authenticate,
            'payload': {
                'name': 'root',
                'password': '12345'
            }
        };
        this.send(o);
        this.send({
            'type': client.OpenDatabase,
            'payload': {
                'name': 'test_01'
            }
        });
        this.send({
            'type': client.ListSubscriptions,
            'payload': {
                'classes': ['Job']
            }
        });
        this.send({
            'type': client.Subscribe,
            'payload': {
                'class': 'Job',
                'key': '(subscription-key)'
            }
        });
    },
    onclose: function (e) {
        console.log(e)
    },
    onmessage: function (data) {
        console.log('Got message:');
        console.dir(data.response);
    }
};

var client = new O2DB(params);
