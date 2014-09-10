<?php

class O2dbClient {

    /** @var resource */
    protected $socket;
    /** @var string */
    protected $address;
    /** @var int */
    protected $port;
    const DELIMITER = 0;

    /**
     * @param string $address
     * @param int $port
     */
    public function __construct($address, $port = 1333) {
        $this->address = $address;
        $this->port = $port;
        $this->connect();
    }

    public function __destruct() {
        socket_close($this->socket);
    }

    /**
     * Attempt to establish connection
     *
     * @throws Exception
     */
    protected function connect() {
        $this->socket = @socket_create(AF_INET, SOCK_STREAM, SOL_TCP);
        if (!$this->socket) {
            throw new Exception(socket_strerror(socket_last_error()));
        }
        if (!@socket_connect($this->socket, $this->address, $this->port)) {
            throw new Exception(socket_strerror(socket_last_error($this->socket)));
        }
    }

    /**
     * Sends message in JSON encoded format and returns raw response
     *
     * @param mixed $message
     *
     * @return string
     */
    public function send($message) {
        $msg = json_encode($message, JSON_PRETTY_PRINT) . chr(self::DELIMITER);
        socket_write($this->socket, $msg, strlen($msg));
        $incoming = '';
        while ($response = socket_read($this->socket, 1)) {
            if (ord($response) === self::DELIMITER) {
                break;
            } else {
                $incoming .= $response;
            }
        }
        return $incoming;
    }
}

$client = new O2dbClient('127.0.0.1');
$message = [
    'type'    => 1,
    'payload' => [
        'class' => 'Customer'
    ],
];
$response = $client->send($message);
echo $response, PHP_EOL;

$message = [
    'type'    => 2,
    'payload' => [
        'hello' => 'World!'
    ],
];
$response = $client->send($message);
echo $response, PHP_EOL;
