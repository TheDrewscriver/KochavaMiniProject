#/bin/sh

export HOME_DIR=`pwd`


#Kill Producer
while IFS= read -r PID
do
   echo "$PID"
   kill -s 9 $PID
done < php_process.pid


#Kill Consumer

while IFS= read -r PID
do
   echo "$PID"
   kill -s 9 $PID
done < go_process.pid


#Kill Kafka
kafka/bin/kafka-server-stop.sh

kafka/bin/zookeeper-server-stop.sh


#Kill Zookeeper
