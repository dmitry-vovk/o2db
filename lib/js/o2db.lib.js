'use strict';

var socket = new WebSocket('ws://localhost:8088/');

socket.onclose = function () {
    console.log('Connection closed');
};

socket.onopen = function () {
    console.log('Connection opened');
    var o = {
        'type': 0,
        'payload': {
            'name': 'root',
            'password': '12345'
        }
    };
    this.send(JSON.stringify(o));
};

socket.onmessage = function (msgEvent) {
    var msg = JSON.parse(msgEvent.data);
    console.dir(msg);
};
