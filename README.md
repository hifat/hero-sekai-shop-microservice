# Hero Sekai Shop Microservice

## Kafka

```sh
docker exec -it kafka-1 bash
```

- Create Topic

```sh
kafka-topics.sh --create --topic inventory --replication-factor 1 --partitions 1 --bootstrap-server localhost:9092
```

- List Topics

```sh
kafka-topics.sh --list --bootstrap-server localhost:9092
```

- Describe

```sh
kafka-topics.sh --describe --topic inventory --bootstrap-server localhost:9092
```