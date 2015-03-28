<?php
/**
 * Collection acceptance tests
 *
 * @author Dmitry Vovk <dmitry.vovk@gmail.com>
 */
require '../../bootstrap.php';
$dbName = 'test_collection_db';
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

$result = $client->createCollection(Entity::class, ['id' => ['type' => 'int',], 'val' => ['type' => 'int',],]);
if (!$result) {
    throw new Exception('Create collection failed!');
} else {
    echo 'Database deleted', PHP_EOL;
}

$entity = new Entity;
$entity->id = 1;
$entity->val = 1;

$result = $client->write($entity);
if (!$result) {
    throw new Exception('Object write failed!');
} else {
    echo 'Object written', PHP_EOL;
}

$result = $client->getOne(Entity::class, 1);
if (!$result) {
    throw new Exception('Object read failed!');
} else {
    echo 'Object read: ';
    print_r($result);
    echo PHP_EOL;
}

$entity->val = 2;
$result = $client->write($entity);
if (!$result) {
    throw new Exception('Object write failed!');
} else {
    echo 'Object written', PHP_EOL;
}

$result = $client->getOne(Entity::class, 1);
if (!$result) {
    throw new Exception('Object read failed!');
} else {
    echo 'Object read: ';
    print_r($result);
    echo PHP_EOL;
}

$result = $client->getObjectVersions(Entity::class, 1);
if (!$result) {
    throw new Exception('Get object versions failed!');
} else {
    echo 'Object versions: ';
    print_r($result);
    echo PHP_EOL;
}

$result = $client->getObjectDiff(Entity::class, 1, 0, 1);
if (!$result) {
    throw new Exception('Get object diff failed!');
} else {
    echo 'Object diff: ';
    print_r($result);
    echo PHP_EOL;
}

$result = $client->dropDatabase($dbName);
if (!$result) {
    throw new Exception('Drop databases failed!');
} else {
    echo 'Database deleted', PHP_EOL;
}
