<?php
/**
 * @author Dmitry Vovk <dmitry.vovk@gmail.com>
 * @created 09/02/15 23:07
 */
require_once 'common.php';
if ($client->authenticate(USERNAME, PASSWORD)) {
    if (!$client->openDatabase(DATABASE)) {
        $client->createDatabase(DATABASE);
        $client->openDatabase(DATABASE);
    }
    $client->createCollection(Entity::class);
    $ent = new Entity;
    $ent->id = 12;
    for ($i = 0; $i < 10000; $i++) {
        $ent->val = $i;
        $client->write($ent);
    }
}
