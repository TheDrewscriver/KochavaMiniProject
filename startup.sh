#/bin/sh

#Start zookeeper
kafka/bin/zookeeper-server-start.sh kafka/config/zookeeper.properties &
sleep 10
#Start kafka
kafka/bin/kafka-server-start.sh kafka/config/server.properties &
sleep 10

#start Php server
cd php
php -S localhost:8000 
