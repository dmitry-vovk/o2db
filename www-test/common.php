<?php
/**
 * @author Dmitry Vovk <dmitry.vovk@gmail.com>
 * @created 09/02/15 23:12
 */
error_reporting(E_ALL);
ini_set('display_errors', true);
ini_set('include_path', ini_get('include_path') . PATH_SEPARATOR . realpath(__DIR__ . '/../lib/php/'));
require 'O2dbClient.php';
require 'entity.php';
$client = new O2dbClient('localhost');

define('DATABASE', 'benchmark-db');
define('USERNAME', 'root');
define('PASSWORD', '12345');
define('SUBSCRIPTION_KEY', '93756028743');
