<?php

require 'O2dbClient.php';
$client = new O2dbClient('127.0.0.1');

// Authenticate
$message = [
    'type' => O2dbClient::TYPE_AUTHENTICATE,
    'payload' => [
        'name' => 'root',
        'password' => '12345',
    ],
];
$response = $client->send($message);
echo '<<<', $response, PHP_EOL;
// Create database

$message = [
    'type' => O2dbClient::TYPE_CREATE_DB,
    'payload' => [
        'name' => 'test_01'
    ],
];
$response = $client->send($message);
echo '<<<', $response, PHP_EOL;

// Open database

$message = [
    'type' => O2dbClient::TYPE_OPEN_DB,
    'payload' => [
        'name' => 'test_01'
    ],
];
$response = $client->send($message);
echo '<<<', $response, PHP_EOL;

// Create collection

$message = [
    'type' => O2dbClient::TYPE_CREATE_COLLECTION,
    'payload' => [
        'class' => 'Job',
        'fields' => [
            'id' => [
                'type' => 'int',
            ],
            'created' => [
                'type' => 'int',
            ],
            'payload' => [
                'type' => 'string',
            ],
            'price' => [
                'type' => 'float',
            ],
        ],
    ],
];
$response = $client->send($message);
echo '<<<', $response, PHP_EOL;

/*
$message = [
    'type'    => O2dbClient::TYPE_CREATE_COLLECTION,
    'payload' => [
        'class'  => 'Batch',
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
*/
// Drop collection
/*
$message = [
    'type'    => O2dbClient::TYPE_DROP_COLLECTION,
    'payload' => [
        'class' => 'Job',
    ],
];
$response = $client->send($message);
echo '<<<', $response, PHP_EOL;
*/
// List databases
/*
$message = [
    'type'    => O2dbClient::TYPE_LIST_DB,
    'payload' => [
        'mask' => '*',
    ],
];
$response = $client->send($message);
echo '<<<', $response, PHP_EOL;
*/
// Write object
/*
$message = [
    'type'    => O2dbClient::TYPE_OBJECT_WRITE,
    'payload' => [
        'class' => 'Job',
        'data'  => [
            'id'    => '5',
            'prop1' => 'val1',
            'prop2' => 'val2',
        ],
    ],
];
$response = $client->send($message);
echo '<<<', $response, PHP_EOL;
$response = $client->send($message);
echo '<<<', $response, PHP_EOL;
$message['payload']['data']['id'] = 13;
$response = $client->send($message);
echo '<<<', $response, PHP_EOL;

$response = $client->getOne('Job', 5);
print_r($response);
echo PHP_EOL;

// Drop database
/*
$message = [
    'type'    => O2dbClient::TYPE_DROP_DB,
    'payload' => [
        'name' => 'test_01',
    ],
];
$response = $client->send($message);
echo '<<<', $response, PHP_EOL;
*/
/*
$message = [
    'type' => O2dbClient::TYPE_OBJECT_WRITE,
    'payload' => [
        'class' => 'Job',
        'data' => [
            'id' => 9,
            'created' => 5,
            'payload' => 'hello there again',
            'price' => 24.1,
            'extra' => 'extra field!',
        ],
    ],
];
$response = $client->send($message);
echo '<<<', $response, PHP_EOL;
*/
/*
$message = [
    'type' => O2dbClient::TYPE_SELECT_OBJECTS,
    'payload' => [
        'class' => 'Job',
        'query' => [
            //'prop1' => [1, 2, 5],
            //'created' => ['<' => 5],
            'price' => 3.5,
            //'payload' => 'hello there',
            //'prop3' => ['<' => 2.5, '>=' => 1],
        ],
    ],
];
$response = $client->send($message);
echo '<<<', $response, PHP_EOL;
*/
/*
$message = [
    'type' => O2dbClient::TYPE_SELECT_OBJECTS,
    'payload' => [
        'class' => 'Job',
        'query' => [
            'XOR' => [
                'price' => 3.5,
                'created' => 5,
            ],
        ],
    ],
];
$response = $client->send($message);
echo '<<<', $response, PHP_EOL;
*/
$message = [
    'type' => O2dbClient::TYPE_ADD_SUBSCRIPTION,
    'payload' => [
        'class' => 'Job',
        'key' => '(subscription-key)',
        'query' => [
            'price' => 3.5,
        ],
    ],
];
$response = $client->send($message);
echo '<<<', $response, PHP_EOL;

$message = [
    'type' => O2dbClient::TYPE_SUBSCRIBE,
    'payload' => [
        'class' => 'Job',
        'key' => '(subscription-key)',
    ],
];
$response = $client->send($message);
echo '<<<', $response, PHP_EOL;

/*
$message = [
    'type' => O2dbClient::TYPE_OBJECT_GET,
    'payload' => [
        'class' => 'Job',
        'data' => [
            'id' => 5,
        ],
    ],
];

$response = $client->send($message);
echo '<<<', $response, PHP_EOL;
*/
/*
$message = [
    'type' => O2dbClient::TYPE_GET_OBJECT_VERSIONS,
    'payload' => [
        'class' => 'Job',
        'id' => 5,
    ],
];
$response = $client->send($message);
echo '<<<', $response, PHP_EOL;
*/
/*
$message = [
    'type' => O2dbClient::TYPE_GET_OBJECT_DIFF,
    'payload' => [
        'class' => 'Job',
        'id' => 5,
        'from' => 6,
        'to' => 5,
    ],
];
$response = $client->send($message);
echo '<<<', $response, PHP_EOL;
*/
