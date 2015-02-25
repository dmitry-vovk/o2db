<?php
/**
 * This script writes 1000 instances of an object and notes the time taken
 *
 * @author Dmitry Vovk <dmitry.vovk@gmail.com>
 * @created 09/02/15 23:07
 */
// How many inserts to make
const TIMES = 1000;

require_once __DIR__ . '/../common.php';
if ($client->authenticate(USERNAME, PASSWORD)) {
    if (!$client->openDatabase(DATABASE)) {
        $client->createDatabase(DATABASE);
        $client->openDatabase(DATABASE);
    }
    $client->createCollection(Entity::class);
    $ent = new Entity;
    $ent->id = 12;
    echo 'Inserting ', TIMES, ' records...', PHP_EOL;
    $start = microtime(true);
    for ($ent->val = 0; $ent->val < TIMES; $ent->val++) {
        $client->write($ent);
    }
    $end = microtime(true);
    echo 'Time taken: ', round($end - $start, 6), ' seconds', PHP_EOL;
    echo '   Average: ', round(($end-$start)/TIMES, 6), ' seconds per insert', PHP_EOL;
}
