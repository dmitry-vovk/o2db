<?php

require 'O2dbClient.php';
$client = new O2dbClient('127.0.0.1');
// Authenticate
$message = [
    'type'    => O2dbClient::TYPE_AUTHENTICATE,
    'payload' => [
        'name'     => 'root',
        'password' => '12345',
    ],
];
$response = $client->send($message);
echo '<<<', $response, PHP_EOL;
// Create database
$message = [
    'type'    => O2dbClient::TYPE_CREATE_DB,
    'payload' => [
        'name' => 'test_01'
    ],
];
$response = $client->send($message);
echo '<<<', $response, PHP_EOL;
// Open database
$message = [
    'type'    => O2dbClient::TYPE_OPEN_DB,
    'payload' => [
        'name' => 'test_01'
    ],
];
$response = $client->send($message);
echo '<<<', $response, PHP_EOL;
// Create collection
$message = [
    'type'    => O2dbClient::TYPE_CREATE_COLLECTION,
    'payload' => [
        'class'  => 'Job',
        'fields' => [
            'id'      => [
                'type'  => 'int',
                'index' => 'primary',
            ],
            'created' => [
                'type'  => 'datetime',
                'index' => 'secondary',
            ],
            'payload' => [
                'type' => 'string',
            ],
        ],
    ],
];
$response = $client->send($message);
echo '<<<', $response, PHP_EOL;
// Drop collection
$message = [
    'type'    => O2dbClient::TYPE_DROP_COLLECTION,
    'payload' => [
        'class' => 'Job',
    ],
];
$response = $client->send($message);
echo '<<<', $response, PHP_EOL;
