# KochavaMiniProject

To start the php server

        php -S localhost:8000

To start kafka and zookeeper

        bin/zookeeper-server-start.sh config/zookeeper.properties &
        bin/kafka-server-start.sh config/server.properties &

To create topic
        bin/kafka-topics.sh --create --topic test --bootstrap-server localhost:9092
