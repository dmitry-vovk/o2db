<?php

/**
 * @author Dmytro Vovk <dmitry.vovk@gmail.com>
 */
class O2dbClient {

    /** @var resource */
    protected $socket;
    /** @var string */
    protected $address;
    /** @var int */
    protected $port;
    const DELIMITER = 0;
    const TYPE_AUTHENTICATE = 0;
    const TYPE_CREATE_DB = 100;
    const TYPE_DROP_DB = 101;
    const TYPE_OPEN_DB = 102;
    const TYPE_LIST_DB = 103;
    const TYPE_CREATE_COLLECTION = 200;
    const TYPE_DROP_COLLECTION = 201;
    const TYPE_LIST_COLLECTIONS = 202;
    const TYPE_OBJECT_GET = 300;
    const TYPE_OBJECT_WRITE = 301;
    const TYPE_OBJECT_DROP = 302;
    /*
     0 TypeAuth
     1 TypeCreateDatabase
     2 TypeDropDatabase
     3 TypeCreateCollection
     4 TypeDropCollection
     5 TypeOpenDatabase
     6 TypeListDatabases
     7 TypeListCollections
     8 TypeObjectInsert
     9 TypeObjectUpdate
    10 TypeObjectDelete
    11 TypeObjectSelect
    12 TypeTransactionStart
    13 TypeTransactionCommit
    14 TypeTransactionAbort
    */
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
        echo '>>>', var_export($msg, true), PHP_EOL;
        socket_write($this->socket, $msg, strlen($msg));
        $incoming = '';
        while (($response = socket_read($this->socket, 1)) !== false) {
            if (ord($response) === self::DELIMITER) {
                break;
            } else {
                $incoming .= $response;
            }
        }
        return $incoming;
    }

    public function createDatabase($dbName) {
        $message = [
            'type'    => O2dbClient::TYPE_CREATE_DB,
            'payload' => [
                'name' => $dbName,
            ],
        ];
        return $this->send($message);
    }
}
