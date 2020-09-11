<?php

include  'kafkaProducer.php';
//Get the data in post call
$json = file_get_contents('php://input');

// Converts it into a PHP object
$data = json_decode($json);
//Push this into kafka
// sendKafkaMessage($json);
$kafkaProducer=new KafkaProducer;
$kafkaProducer->sendKafkaMessage($json);
