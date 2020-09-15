#!/bin/sh


rm *.pid

#Start zookeeper
kafka/bin/zookeeper-server-start.sh kafka/config/zookeeper.properties >/dev/null 2>&1 &
sleep 20
# #Start kafka
kafka/bin/kafka-server-start.sh kafka/config/server.properties  >/dev/null 2>&1  &
sleep 20
# #Create topic if it doesn't exist
kafka/bin/kafka-topics.sh --create --if-not-exists  --topic kochavaPostback --bootstrap-server localhost:9092 &


#Start go consumer
cd golang
./KafkaConsumer &
cd ..
PID=$!
echo $PID > go_process.pid


#start Php server
cd php
composer install
php -S 0.0.0.0:8000
