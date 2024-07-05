inside an ec2 instance(which has docker and docker compose installed)
run the below commands:

sudo yum update -y

sudo yum install docker -y

sudo service docker start

sudo usermod -aG docker ec2-user

exit

and reconnect to ec2 again

docker pull wurstmeister/kafka

docker run -d --name zookeeper -p 2181:2181 wurstmeister/zookeeper

docker run -d --name kafka -p 9092:9092 --link zookeeper:zookeeper \
-e KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181 \
-e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092 \
-e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://54.204.131.58:9092 \
-e KAFKA_LISTENER_SECURITY_PROTOCOL_MAP=PLAINTEXT:PLAINTEXT \
-e KAFKA_INTER_BROKER_LISTENER_NAME=PLAINTEXT \
-e KAFKA_BROKER_ID=1 wurstmeister/kafka



(publicIp address: 3.90.212.230:9092 )


then test the microservice using this below endpoint
http://localhost:8000/send-message


2 containers should be running

--->in go code public address is fetched using running instance id
