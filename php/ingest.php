<?php

include  'kafkaProducer.php';
//Get the data in post call
$json = file_get_contents('php://input');

// Converts it into a PHP object
$data = json_decode($json);

if(strtolower($data->endpoint->method) != "get" && strtolower($data->endpoint->method) != "post"){
  print("Only GET or POST supported for method field");
  http_response_code(400);
  return;
}


if (filter_var($data->endpoint->url, FILTER_VALIDATE_URL) == false){
  print("Invalid URL provided in endpoint field");
  http_response_code(400);
  return;
}

//Push this into kafka
$kafkaProducer=new KafkaProducer;
$kafkaProducer->sendKafkaMessage($json);
