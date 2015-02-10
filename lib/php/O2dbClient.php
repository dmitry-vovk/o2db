<?php

/**
 * @author Dmytro Vovk <dmitry.vovk@gmail.com>
 */
class O2dbClient
{

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
    const TYPE_GET_OBJECT_VERSIONS = 303;
    const TYPE_GET_OBJECT_DIFF = 304;
    const TYPE_SELECT_OBJECTS = 305;
    const TYPE_SUBSCRIBE = 400;
    const TYPE_ADD_SUBSCRIPTION = 401;
    const TYPE_CANCEL_SUBSCRIPTION = 402;
    const TYPE_LIST_SUBSCRIPTIONS = 403;
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
    // TODO collection
    // TODO objects
    // TODO data
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
    public function __construct($address, $port = 1333)
    {
        $this->address = $address;
        $this->port = $port;
        $this->connect();
    }

    public function __destruct()
    {
        socket_close($this->socket);
    }

    /**
     * Attempt to establish connection
     *
     * @throws Exception
     */
    protected function connect()
    {
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
    public function send($message)
    {
        $msg = json_encode($message, JSON_PRETTY_PRINT) . chr(self::DELIMITER);
        //echo '>>>', var_export($msg, true), PHP_EOL;
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

    /**
     * @param string $result
     * @return bool
     */
    protected function parseResult($result)
    {
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
    public function success()
    {
        return $this->lastResult;
    }

    /**
     * @return int
     */
    public function getCode()
    {
        return $this->lastStatusCode;
    }

    /**
     * @return mixed
     */
    public function getResponse()
    {
        return $this->lastResponse;
    }

    /**
     * @param string $username
     * @param string $password
     * @return bool
     */
    public function authenticate($username, $password)
    {
        $message = [
            'type' => O2dbClient::TYPE_AUTHENTICATE,
            'payload' => [
                'name' => $username,
                'password' => $password,
            ],
        ];
        if ($this->parseResult($this->send($message))) {
            return $this->lastResult;
        }
        return false;
    }

    /**
     * @param string $dbName
     * @return bool
     */
    public function createDatabase($dbName)
    {
        $message = [
            'type' => O2dbClient::TYPE_CREATE_DB,
            'payload' => [
                'name' => $dbName,
            ],
        ];
        if ($this->parseResult($this->send($message))) {
            return $this->lastResult;
        }
        return false;
    }

    /**
     * @param string $dbName
     * @return bool
     */
    public function openDatabase($dbName)
    {
        $message = [
            'type' => O2dbClient::TYPE_OPEN_DB,
            'payload' => [
                'name' => $dbName,
            ],
        ];
        if ($this->parseResult($this->send($message))) {
            return $this->lastResult;
        }
        return false;
    }

    /**
     * @param string $collectionName
     * @param array $indices
     * @return bool
     */
    public function createCollection($collectionName, array $indices = ['id' => ['type' => 'int']])
    {
        $message = [
            'type' => O2dbClient::TYPE_CREATE_COLLECTION,
            'payload' => [
                'class' => $collectionName,
                'fields' => $indices,
            ],
        ];
        if ($this->parseResult($this->send($message))) {
            return $this->lastResult;
        }
        return false;
    }

    /**
     * @param $object
     * @return bool
     */
    public function write($object)
    {
        $message = [
            'type' => self::TYPE_OBJECT_WRITE,
            'payload' => [
                'class' => get_class($object),
                'data' => $object,
            ],
        ];
        if ($this->parseResult($this->send($message))) {
            return $this->lastResult;
        }
        return false;
    }

    /**
     * @param string $class
     * @param int $id
     * @return bool
     */
    public function getOne($class, $id)
    {
        $message = [
            'type' => O2dbClient::TYPE_OBJECT_GET,
            'payload' => [
                'class' => $class,
                'data' => [
                    'id' => $id,
                ],
            ],
        ];
        if ($this->parseResult($this->send($message))) {
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
     * @param string $key
     * @param array $mask
     * @return bool
     */
    public function createSubscription($class, $key, array $mask)
    {
        $message = [
            'type' => self::TYPE_ADD_SUBSCRIPTION,
            'payload' => [
                'class' => $class,
                'key' => $key,
                'query' => $mask,
            ],
        ];
        if ($this->parseResult($this->send($message))) {
            return $this->lastResult;
        }
        return false;
    }
}
