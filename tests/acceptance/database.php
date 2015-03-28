<?php
/**
 * Acceptance tests for database creation, listing, deleting
 *
 * @author Dmitry Vovk <dmitry.vovk@gmail.com>
 */
require '../../bootstrap.php';

$prefix = 'atdb_';
$result = $client->authenticate(USERNAME, PASSWORD);
if (!$result) {
    throw new Exception('Authentication failed!');
} else {
    echo 'Authenticated', PHP_EOL;
}

$result = $client->createDatabase($prefix . '1');
if (!$result) {
    throw new Exception('Create database failed!');
} else {
    echo 'Database 1 created', PHP_EOL;
}

$result = $client->createDatabase($prefix . '2');
if (!$result) {
    throw new Exception('Create database failed!');
} else {
    echo 'Database 2 created', PHP_EOL;
}

$result = $client->createDatabase($prefix . '3');
if (!$result) {
    throw new Exception('Create database failed!');
} else {
    echo 'Database 3 created', PHP_EOL;
}

$result = $client->listDatabases($prefix . '*');
if (!$result) {
    throw new Exception('List databases failed!');
} else {
    echo 'Databases listed: ';
    print_r($result);
    echo PHP_EOL;
}

$result = $client->dropDatabase($prefix . '1');
if (!$result) {
    throw new Exception('Drop databases failed!');
} else {
    echo 'Database 1 deleted', PHP_EOL;
}

$result = $client->dropDatabase($prefix . '2');
if (!$result) {
    throw new Exception('Drop databases failed!');
} else {
    echo 'Database 1 deleted', PHP_EOL;
}

$result = $client->dropDatabase($prefix . '3');
if (!$result) {
    throw new Exception('Drop databases failed!');
} else {
    echo 'Database 1 deleted', PHP_EOL;
}

$result = $client->listDatabases($prefix . '*');
if (!$result) {
    throw new Exception('List databases failed!');
} else {
    echo 'Databases listed: ';
    print_r($result);
    echo PHP_EOL;
}
