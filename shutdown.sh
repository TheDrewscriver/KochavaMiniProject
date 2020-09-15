#/bin/sh


#Kill Consumer
while IFS= read -r PID
do
   echo "$PID"
   kill -s 9 $PID
done < go_process.pid


#Kill Kafka
kafka/bin/zookeeper-server-stop.sh
sleep 20
kafka/bin/kafka-server-stop.sh
sleep 20




#Kill Zookeeper
