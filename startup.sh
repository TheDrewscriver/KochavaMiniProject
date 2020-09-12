#/bin/sh

export HOME_DIR=`pwd`

#Start go consumer
golang/KafkaConsumer &
PID=$!
echo $PID > go_process.pid


#Start zookeeper
kafka/bin/zookeeper-server-start.sh kafka/config/zookeeper.properties &
sleep 10
#Start kafka
kafka/bin/kafka-server-start.sh kafka/config/server.properties &
sleep 10
#Create topic if it doesn't exist
kafka/bin/kafka-topics.sh --create --if-not-exists  --topic kochavaPostback --bootstrap-server localhost:9092


#start Php server
php -S localhost:8000 -t php &
PID=$!
echo $PID > php_process.pid
cd ..
