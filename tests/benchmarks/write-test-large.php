<?php
/**
 * This script writes 1000 instances of an object and notes the time taken
 *
 * @author Dmitry Vovk <dmitry.vovk@gmail.com>
 * @created 09/02/15 23:07
 */
// How many inserts to make
const TIMES = 1000;

require '../../bootstrap.php';
if ($client->authenticate(USERNAME, PASSWORD)) {
    if (!$client->openDatabase(DATABASE)) {
        $client->createDatabase(DATABASE);
        $client->openDatabase(DATABASE);
    }
    $client->createCollection(LargeEntity::class);
    $entities = [];
    for ($count = 0; $count < TIMES; $count++) {
        $entities[] = new LargeEntity;
    }
    echo 'Inserting ', TIMES, ' records...', PHP_EOL;
    $start = microtime(true);
    foreach ($entities as $ent) {
        $client->write($ent);
    }
    $end = microtime(true);
    echo 'Time taken: ', round($end - $start, 6), ' seconds', PHP_EOL;
    echo '   Average: ', round(($end - $start) / TIMES, 6), ' seconds per insert', PHP_EOL;
    $result = $client->getObjectVersions(LargeEntity::class, $ent->id);
    echo '  Versions: ', $result, PHP_EOL;
    $client->dropDatabase(DATABASE);
}
