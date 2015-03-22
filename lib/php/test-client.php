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
$response = $client->sendAndRead($message);
echo $response, PHP_EOL;
// Create database

$message = [
    'type' => O2dbClient::TYPE_CREATE_DB,
    'payload' => [
        'name' => 'test_01'
    ],
];
$response = $client->sendAndRead($message);
echo $response, PHP_EOL;

// Open database

$message = [
    'type' => O2dbClient::TYPE_OPEN_DB,
    'payload' => [
        'name' => 'test_01'
    ],
];
$response = $client->sendAndRead($message);
echo $response, PHP_EOL;

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
$response = $client->sendAndRead($message);
echo $response, PHP_EOL;

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

$message = [
    'type' => O2dbClient::TYPE_OBJECT_WRITE,
    'payload' => [
        'class' => 'Job',
        'data' => [
            'id' => 5,
            'payload' => 'payload here',
            'price' => 3.5,
            'hello' => 'there',
        ],
    ],
];
$response = $client->sendAndRead($message);
echo $response, PHP_EOL;
/*
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
            'OR' => [
                'price' => 3.5,
                'AND' => [
                    'payload' => 'empty',
                    'price' => 3.5,
                ]
            ]
        ],
    ],
];
$response = $client->sendAndRead($message);
echo $response, PHP_EOL;

/*
$message = [
    'type' => O2dbClient::TYPE_SUBSCRIBE,
    'payload' => [
        'class' => 'Job',
        'key' => '(subscription-key)',
    ],
];
$response = $client->send($message);
echo '<<<', $response, PHP_EOL;
*/
/*
$message = [
    'type' => O2dbClient::TYPE_LIST_SUBSCRIPTIONS,
    'payload' => [
        'classes' => ['Job'],
    ],
];
$response = $client->send($message);
echo '<<<', $response, PHP_EOL;
*/
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
