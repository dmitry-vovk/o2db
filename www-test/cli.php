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
    $ent->val = 'hello';
    $resp = $client->save($ent);
    print_r($resp);
}
