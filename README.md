# Hero Sekai Shop Microservice

## Kafka

```sh
docker exec -it kafka-1 bash
```

- Create Topic

```sh
kafka-topics.sh --create --topic inventory --replication-factor 1 --partitions 1 --bootstrap-server localhost:9092
kafka-topics.sh --create --topic payment --replication-factor 1 --partitions 1 --bootstrap-server localhost:9092
kafka-topics.sh --create --topic player --replication-factor 1 --partitions 1 --bootstrap-server localhost:9092

```

- Add Topic Retention

```sh
kafka-configs.sh --bootstrap-server localhost:9092 --entity-type topics --entity-name inventory --alter --add-config retention.ms=180000
kafka-configs.sh --bootstrap-server localhost:9092 --entity-type topics --entity-name payment --alter --add-config retention.ms=180000
kafka-configs.sh --bootstrap-server localhost:9092 --entity-type topics --entity-name player --alter --add-config retention.ms=180000

```

- List Topics

```sh
kafka-topics.sh --list --bootstrap-server localhost:9092
```

- Describe

```sh
kafka-topics.sh --describe --topic inventory --bootstrap-server localhost:9092
```

- Write a message into the topic

```sh
kafka-console-producer.sh --topic inventory --bootstrap-server localhost:9092 
kafka-console-producer.sh --topic payment --bootstrap-server localhost:9092 
kafka-console-producer.sh --topic player --bootstrap-server localhost:9092
```

Write a message with key into the topic

```
--property "key.separator=:" --property "parse.key=true"
```

- Read a message on that topic

```sh
kafka-console-consumer.sh --topic inventory --from-beginning --bootstrap-server localhost:9092
kafka-console-consumer.sh --topic payment --from-beginning --bootstrap-server localhost:9092
kafka-console-consumer.sh --topic player --from-beginning --bootstrap-server localhost:9092
```

- Delete topic

```sh
kafka-topics.sh --delete --bootstrap-server localhost:9092 --topic inventory
kafka-topics.sh --delete --bootstrap-server localhost:9092 --topic payment
kafka-topics.sh --delete --bootstrap-server localhost:9092 --topic player
```

```sh
kafka-console-consumer --bootstrap-server localhost:9092 --topic payment --from-beginning
```

```sh
kafka-console-producer --bootstrap-server localhost:9092 --topic payment
```