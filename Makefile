kafka-topic:
	docker run --net=host --rm confluentinc/cp-kafka:latest kafka-topics --create --topic login-succeed --bootstrap-server localhost:19091 --partitions 2 --replication-factor 1

kafka-cat1:
	kafkacat -C -b localhost:19091 -t login-succeed -p 0

kafka-publish:
	echo 'test' | kafkacat -P -b localhost:19091 -t login-succeed -p 0

.PHONY: kafka-topic kafka-cat1 kafka-publish