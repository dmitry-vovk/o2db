<?php
/**
 * Subscriptions acceptance tests
 *
 * @author Dmitry Vovk <dmitry.vovk@gmail.com>
 */
require '../../bootstrap.php';

$dbName = 'test_subscription_db';

$result = $client->authenticate(USERNAME, PASSWORD);
if (!$result) {
    throw new Exception('Authentication failed!');
} else {
    echo 'Authenticated', PHP_EOL;
}

$result = $client->createDatabase($dbName);
if (!$result) {
    throw new Exception('Create database failed!');
} else {
    echo 'Database created', PHP_EOL;
}

$result = $client->openDatabase($dbName);
if (!$result) {
    throw new Exception('Open database failed!');
} else {
    echo 'Database opened', PHP_EOL;
}

$result = $client->createCollection(
    Entity::class,
    [
        'id'  => [
            'type' => 'int',
        ],
        'val' => [
            'type' => 'int',
        ],
    ]
);
if (!$result) {
    throw new Exception('Create collection failed!');
} else {
    echo 'Collection created', PHP_EOL;
}

$result = $client->createSubscription(Entity::class, SUBSCRIPTION_KEY, ['id' => 1,]);
if (!$result) {
    throw new Exception('Create subscription failed!');
} else {
    echo 'Subscription created', PHP_EOL;
}

$result = $client->listSubscriptions([Entity::class]);
if (!$result) {
    throw new Exception('List subscriptions failed!');
} else {
    echo 'Subscriptions list: ';
    print_r($result);
    echo PHP_EOL;
}

$result = $client->cancelSubscription(Entity::class, SUBSCRIPTION_KEY);
if (!$result) {
    throw new Exception('Cancel subscription failed!');
} else {
    echo 'Subscription cancelled', PHP_EOL;
}

$result = $client->dropDatabase($dbName);
if (!$result) {
    throw new Exception('Drop databases failed!');
} else {
    echo 'Database deleted', PHP_EOL;
}
