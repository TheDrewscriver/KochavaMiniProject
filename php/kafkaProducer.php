<?php
require './vendor/autoload.php';
date_default_timezone_set('PRC');
use Monolog\Logger;
use Monolog\Handler\StreamHandler;

class KafkaProducer
{
  public function sendKafkaMessage($json){
    $logger = new Logger('my_logger');
    // Now add some handlers
    $logger->pushHandler(new StreamHandler('php://stdout'));

    $config = \Kafka\ProducerConfig::getInstance();
    $config->setMetadataRefreshIntervalMs(1000);
    $config->setMetadataBrokerList('127.0.0.1:9092');
    $config->setBrokerVersion('1.0.0');
    $config->setRequiredAck(1);
    $config->setIsAsyn(false);
    $config->setProduceInterval(500);
    $producer = new \Kafka\Producer();
    $producer->setLogger($logger);
    $producer->send([
      [
        'topic' => 'kochavaPostback',
        'value' => $json,
        'key' => '',
      ],
      ]);
    }
}

?>
