<?php
/**
 * @author Dmitry Vovk <dmitry.vovk@gmail.com>
 * @created 09/02/15 22:01
 */
require_once 'common.php';
if ($client->authenticate(USERNAME, PASSWORD)) {
    if (!$client->openDatabase(DATABASE) && $client->getCode() === O2dbClient::RESP_DATABASE_DOES_NOT_EXIST) {
        $client->createDatabase(DATABASE);
        $client->openDatabase(DATABASE);
        $client->createCollection(Entity::class);
    }
    $resp = $client->createSubscription(Entity::class, SUBSCRIPTION_KEY, ['id' => 12]);
}
?><!doctype html>
<html>
<head>
    <title>O2DB tests</title>
    <script type="text/javascript">
        var db = '<?= DATABASE ?>';
        var collection = '<?= Entity::class ?>';
        var key = '<?= SUBSCRIPTION_KEY ?>';
    </script>
    <script src="o2db.lib.js" type="text/javascript"></script>
</head>
<body>
<pre>
<?php
$resp = $client->getOne(Entity::class, 12);
var_dump($resp);
?>
</pre>
</body>
</html>
