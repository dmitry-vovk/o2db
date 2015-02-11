<?php
/**
 * @author Dmitry Vovk <dmitry.vovk@gmail.com>
 * @created 09/02/15 22:01
 */
require_once 'common.php';
if ($client->authenticate(USERNAME, PASSWORD)) {
    error_log('Authenticated');
    if (!$client->openDatabase(DATABASE)) {
        error_log('Code ' . $client->getCode());
        error_log('Creating database');
        $client->createDatabase(DATABASE);
        error_log('Opening database');
        $client->openDatabase(DATABASE);
        error_log('Creating collection');
        $client->createCollection(Entity::class);
    }
    $resp = $client->createSubscription(Entity::class, SUBSCRIPTION_KEY, ['id' => 12]);
} else {
    error_log('Not authenticated');
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
<progress id="progress" max="10" value="0"></progress>
</body>
</html>
