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
    /** @var $ent Entity|bool */
    $ent = $client->getOne(Entity::class, 12);
    if ($ent) {
        $ent->val < 10
            ? $ent->val++
            : $ent->val = 0;
    } else {
        $ent = new Entity;
        $ent->id = 12;
        $ent->val = 0;
    }
    $client->write($ent);
    echo 'Written ';
    print_r($ent);
}
