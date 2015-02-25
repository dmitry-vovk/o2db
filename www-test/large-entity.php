<?php

/**
 * @author Dmitry Vovk <dmitry.vovk@gmail.com>
 * @created 09/02/15 23:43
 */
class LargeEntity
{
    public $id = 17;
    public $val = [];

    public function __construct() {
        for ($j = 0; $j < 64; $j++) {
            $row = [];
            for ($i = 0; $i < 64; $i++) {
                $row[] = rand(0, 255);
            }
            $this->val[] = $row;
        }
    }
}
