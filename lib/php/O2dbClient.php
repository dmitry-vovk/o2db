<?php

/**
 * O2DB PHP Client
 *
 * @author Dmytro Vovk <dmitry.vovk@gmail.com>
 */
class O2dbClient {

    /** @var resource */
    protected $socket;
    /** @var string */
    protected $address;
    /** @var int */
    protected $port;
    /** Message delimiter */
    const DELIMITER = 0;
    /** Query types */
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
    const TYPE_GET_OBJECT_VERSIONS = 303;
    const TYPE_GET_OBJECT_DIFF = 304;
    const TYPE_SELECT_OBJECTS = 305;
    const TYPE_SUBSCRIBE = 400;
    const TYPE_ADD_SUBSCRIPTION = 401;
    const TYPE_CANCEL_SUBSCRIPTION = 402;
    const TYPE_LIST_SUBSCRIPTIONS = 403;
    /** Response types */
    const RESP_NO_ERROR = 0;
    const RESP_AUTHENTICATED = 1;
    const RESP_NOT_AUTHENTICATED = 2;
    const RESP_AUTH_REQUIRED = 3;
    const RESP_DATABASE_CREATED = 4;
    const RESP_DATABASE_DELETED = 5;
    const RESP_DATABASE_OPENED = 6;
    const RESP_DATABASE_LIST = 7;
    const RESP_DATABASE_ALREADY_EXISTS = 8;
    const RESP_DATABASE_NOT_SELECTED = 9;
    const RESP_DATABASE_DOES_NOT_EXIST = 10;
    const RESP_COLLECTION_CREATED = 11;
    const RESP_COLLECTION_DELETED = 12;
    const RESP_COLLECTION_ALREADY_EXISTS = 13;
    const RESP_COLLECTION_DOES_NOT_EXIST = 14;
    const RESP_COLLECTION_LIST = 15;
    const RESP_OBJECT = 16;
    const RESP_OBJECT_WRITTEN = 17;
    const RESP_OBJECT_INVALID = 18;
    const RESP_OBJECT_ENCODE_ERROR = 19;
    const RESP_OBJECT_DECODE_ERROR = 20;
    const RESP_OBJECT_DOES_NOT_EXIST = 21;
    const RESP_OBJECT_NOT_FOUND = 22;
    const RESP_DATA_WRITE_ERROR = 23;
    const RESP_DATA_READ_ERROR = 24;
    const RESP_SUBSCRIBED = 25;
    const RESP_UNSUBSCRIBED = 26;
    const RESP_SUBSCRIPTION_INVALID_FORMAT = 27;
    const RESP_SUBSCRIPTION_CREATED = 28;
    const RESP_SUBSCRIPTION_CANCELLED = 29;
    const RESP_SUBSCRIPTION_ALREADY_EXISTS = 30;
    const RESP_SUBSCRIPTION_DOES_NOT_EXIST = 31;
    const RESP_SUBSCRIPTIONS_LIST = 32;
    /** @var bool */
    protected $lastResult = true;
    /** @var int */
    protected $lastStatusCode = 0;
    /** @var mixed */
    protected $lastResponse = '';

    /**
     * @param string $address
     * @param int $port
     */
    public function __construct($address, $port = 1333) {
        assert('is_string($address)');
        assert('!empty($address)');
        assert('is_int($port)');
        assert('$port > 0');
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
    public function sendAndRead($message) {
        $this->send($message);
        return $this->read();
    }

    /**
     * @param mixed $message
     */
    public function send($message) {
        $msg = json_encode($message, JSON_PRETTY_PRINT) . chr(self::DELIMITER);
        socket_write($this->socket, $msg, strlen($msg));
    }

    /**
     * Reads a message from the socket
     *
     * @param bool $blocking
     *
     * @return string
     */
    public function read($blocking = true) {
        $incoming = '';
        $blocking
            ? socket_set_block($this->socket)
            : socket_set_nonblock($this->socket);
        while (($response = socket_read($this->socket, 1)) !== false) {
            if (ord($response) === self::DELIMITER) {
                break;
            } else {
                $incoming .= $response;
            }
        }
        return $incoming;
    }

    /**
     * @param string $result
     *
     * @return bool
     */
    protected function parseResult($result) {
        $parsed = json_decode($result, true);
        $jsonError = json_last_error();
        if ($jsonError !== JSON_ERROR_NONE) {
            $this->lastResult = false;
            $this->lastStatusCode = -1;
            $this->lastResponse = $jsonError;
            return false;
        } else {
            $this->lastResult = $parsed['result'];
            $this->lastStatusCode = $parsed['code'];
            $this->lastResponse = $parsed['response'];
        }
        return true;
    }

    /**
     * @return bool
     */
    public function success() {
        return $this->lastResult;
    }

    /**
     * @return int
     */
    public function getCode() {
        return $this->lastStatusCode;
    }

    /**
     * @return mixed
     */
    public function getResponse() {
        return $this->lastResponse;
    }

    /**
     * @param string $username
     * @param string $password
     *
     * @return bool
     */
    public function authenticate($username, $password) {
        $message = [
            'type'    => O2dbClient::TYPE_AUTHENTICATE,
            'payload' => [
                'name'     => $username,
                'password' => $password,
            ],
        ];
        if ($this->parseResult($this->sendAndRead($message))) {
            return $this->lastResult;
        }
        return false;
    }

    /**
     * @param string $dbName
     *
     * @return bool
     */
    public function createDatabase($dbName) {
        $message = [
            'type'    => O2dbClient::TYPE_CREATE_DB,
            'payload' => [
                'name' => $dbName,
            ],
        ];
        if ($this->parseResult($this->sendAndRead($message))) {
            return $this->lastResult;
        }
        return false;
    }

    /**
     * @param string $dbName
     *
     * @return bool
     */
    public function openDatabase($dbName) {
        $message = [
            'type'    => O2dbClient::TYPE_OPEN_DB,
            'payload' => [
                'name' => $dbName,
            ],
        ];
        if ($this->parseResult($this->sendAndRead($message))) {
            return $this->lastResult;
        }
        return false;
    }

    /**
     * @param string $mask
     *
     * @return bool
     */
    public function listDatabases($mask = '*') {
        $message = [
            'type'    => O2dbClient::TYPE_LIST_DB,
            'payload' => [
                'mask' => $mask,
            ],
        ];
        if ($this->parseResult($this->sendAndRead($message))) {
            return $this->lastResponse;
        }
        return false;
    }

    /**
     * @param string $dbName
     *
     * @return bool
     */
    public function dropDatabase($dbName) {
        $message = [
            'type'    => O2dbClient::TYPE_DROP_DB,
            'payload' => [
                'name' => $dbName,
            ],
        ];
        if ($this->parseResult($this->sendAndRead($message))) {
            return $this->lastResult;
        }
        return false;
    }

    /**
     * @param string $collectionName
     * @param array $indices
     *
     * @return bool
     */
    public function createCollection($collectionName, array $indices = ['id' => ['type' => 'int',],]) {
        $message = [
            'type'    => O2dbClient::TYPE_CREATE_COLLECTION,
            'payload' => [
                'class'  => $collectionName,
                'fields' => $indices,
            ],
        ];
        if ($this->parseResult($this->sendAndRead($message))) {
            return $this->lastResult;
        }
        return false;
    }

    /**
     * @param $object
     *
     * @return bool
     */
    public function write($object) {
        $message = [
            'type'    => self::TYPE_OBJECT_WRITE,
            'payload' => [
                'class' => get_class($object),
                'data'  => $object,
            ],
        ];
        if ($this->parseResult($this->sendAndRead($message))) {
            return $this->lastResult;
        }
        return false;
    }

    /**
     * @param string $class
     * @param int $id
     *
     * @return bool
     */
    public function getOne($class, $id) {
        $message = [
            'type'    => O2dbClient::TYPE_OBJECT_GET,
            'payload' => [
                'class' => $class,
                'data'  => [
                    'id' => $id,
                ],
            ],
        ];
        if ($this->parseResult($this->sendAndRead($message))) {
            if (is_array($this->lastResponse)) {
                $object = new $class;
                foreach ($this->lastResponse as $key => $value) {
                    $object->{$key} = $value;
                }
                return $object;
            }
        }
        return false;
    }

    /**
     * @param string $class
     * @param int $id
     *
     * @return bool|mixed
     */
    public function getObjectVersions($class, $id) {
        $message = [
            'type'    => self::TYPE_GET_OBJECT_VERSIONS,
            'payload' => [
                'class' => $class,
                'id'    => $id,
            ],
        ];
        if ($this->parseResult($this->sendAndRead($message))) {
            return $this->lastResponse;
        }
        return false;
    }

    /**
     * @param string $class
     * @param int $id
     * @param int $v1
     * @param int $v2
     *
     * @return bool|mixed
     */
    public function getObjectDiff($class, $id, $v1, $v2) {
        $message = [
            'type'    => self::TYPE_GET_OBJECT_DIFF,
            'payload' => [
                'class' => $class,
                'id'    => $id,
                'from'  => $v1,
                'to'    => $v2,
            ],
        ];
        if ($this->parseResult($this->sendAndRead($message))) {
            return $this->lastResponse;
        }
        return false;
    }

    /**
     * @param string $class
     * @param string $key
     * @param array $mask
     *
     * @return bool
     */
    public function createSubscription($class, $key, array $mask) {
        $message = [
            'type'    => self::TYPE_ADD_SUBSCRIPTION,
            'payload' => [
                'class' => $class,
                'key'   => $key,
                'query' => $mask,
            ],
        ];
        if ($this->parseResult($this->sendAndRead($message))) {
            return $this->lastResult;
        }
        return false;
    }

    /**
     * @param string $class
     * @param string $key
     *
     * @return bool
     */
    public function cancelSubscription($class, $key) {
        $message = [
            'type'    => self::TYPE_CANCEL_SUBSCRIPTION,
            'payload' => [
                'class' => $class,
                'key'   => $key,
            ],
        ];
        if ($this->parseResult($this->sendAndRead($message))) {
            return $this->lastResult;
        }
        return false;
    }

    /**
     * @param array $classes
     *
     * @return bool|mixed
     */
    public function listSubscriptions(array $classes) {
        $message = [
            'type'    => self::TYPE_CANCEL_SUBSCRIPTION,
            'payload' => [
                'classes' => $classes,
            ],
        ];
        if ($this->parseResult($this->sendAndRead($message))) {
            return $this->lastResponse;
        }
        return false;
    }
}
