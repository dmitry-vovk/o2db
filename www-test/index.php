<?php
/**
 * @author Dmitry Vovk <dmitry.vovk@gmail.com>
 * @created 09/02/15 22:01
 */
require_once 'common.php';
?><!doctype html>
<html>
<head>
    <title>O2DB tests</title>
    <script src="o2db.lib.js" type="text/javascript"></script>
</head>
<body>
<?php

if ($client->authenticate(USERNAME, PASSWORD)) {
    if (!$client->openDatabase(DATABASE)) {
        $client->createDatabase(DATABASE);
        $client->openDatabase(DATABASE);
    }
    /**
    $client->createCollection(Entity::class);
    $ent = new Entity;
    $ent->id = 12;
    $ent->val = 'hello';
    $resp = $client->save($ent);

     * */
    $resp = $client->getOne(Entity::class, 12);
    print_r($resp);
}

?>
</body>
</html>
