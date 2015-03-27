<?php
/**
 * Basic bootstrap file
 *
 * @author Dmitry Vovk <dmitry.vovk@gmail.com>
 * @created 09/02/15 23:12
 */
error_reporting(E_ALL);
define('BASE_DIR', __DIR__);
ini_set('display_errors', true);

require BASE_DIR . '/../lib/php/O2dbClient.php';
require BASE_DIR . '/classes/entity.php';
require BASE_DIR . '/classes/large-entity.php';

define('DATABASE', 'benchmark-db');
define('USERNAME', 'root');
define('PASSWORD', '12345');
// Regenerate subscription key every time
define('SUBSCRIPTION_KEY', uniqid("sk-", true));

// Create instance of the DBMS client
$client = new O2dbClient('localhost');
