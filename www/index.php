<?php
/**
 * @author Dmitry Vovk <dmitry.vovk@gmail.com>
 * @created 09/02/15 22:01
 */
require_once '../bootstrap.php';
/**
 * Preparation for front-end interaction:
 * 1: Authenticate back-end client
 * 2: Create a database
 * 3: Open the database
 * 4: Create collection
 * 5: Create subscription
 * 6: Provide front-end with database and collection name, and subscription key
 */
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
    // This subscription will filter all objects of class Entity, with id = 12
    $resp = $client->createSubscription(Entity::class, SUBSCRIPTION_KEY, ['id' => 12]);
} else {
    error_log('Not authenticated');
}
?><!doctype html>
<html>
<head>
    <title>O2DB tests</title>
    <script type="text/javascript">
        // Connection and subscription information
        var host = 'localhost:8088',
            db = '<?= DATABASE ?>',
            collection = '<?= Entity::class ?>',
            key = '<?= SUBSCRIPTION_KEY ?>';
    </script>
    <script src="o2db.lib.js" type="text/javascript"></script>
    <script>
        var params = {
            host: 'ws://' + host,
            /**
             * Function called after the connection is established
             */
            onopen: function () {
                this.send({
                    'type': client.Subscribe,
                    'payload': {
                        'database': db,
                        'class': collection,
                        'key': key
                    }
                });
            },
            onclose: function (e) {
                console.log(e)
            },
            /**
             * Function that will receive pushed messages
             * @param data
             */
            onmessage: function (data) {
                var bar = document.getElementById('progress');
                bar.setAttribute('value', data.response.val);
            }
        };
        var client = new O2DB(params);
    </script>
</head>
<body>
<pre>
<?php
// Get sample object
$resp = $client->getOne(Entity::class, 12);
// and display it
var_dump($resp);
?>
</pre>
<!-- Progress bar that will reflect objects change -->
<progress id="progress" max="10" value="0"></progress>
</body>
</html>
