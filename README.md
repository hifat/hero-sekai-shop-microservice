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

## K8S

```sh
kubectl create configmap <NAME> --from-file=<ENV_PATH>
```

```sh
kubectl delete configmap <NAME>
```

```sh
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.12.2/deploy/static/provider/cloud/deploy.yaml
```

```sh
kubectl get ingress
```

```sh
kubectl get deployment 
```

```sh
kubectl get pods 
```

```sh
kubectl logs -f <POD_NAME>
```

```sh
kubectl scale deployment/nginx-deployment --replicas=10
```

```sh
kubectl rollout restart deployment/<DEPLOYMENT_NAME>
```

- Check Ingress state

```sh
kubectl get svc -n ingress-nginx
```

## Source

[k8s service](https://kubernetes.io/docs/concepts/services-networking/service/)
[k8s deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/)
[ingress nginx install guide](https://kubernetes.github.io/ingress-nginx/deploy/)
[scaling a deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#scaling-a-deployment)


## Problem&Solve

### ingress is on pending state

```md
NAME                                 TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
ingress-nginx-controller             LoadBalancer   10.107.102.222   <pending>     80:31363/TCP,443:30797/TCP   2d
```

1. Delete existing ingress-nginx deployment:

```sh
kubectl delete -A ValidatingWebhookConfiguration ingress-nginx-admission
kubectl delete namespace ingress-nginx
```

2. Install NGINX Ingress Controller specifically for Docker Desktop:

```sh
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.12.2/deploy/static/provider/cloud/deploy.yaml
```

3. Wait for the deployment to complete:

```sh
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=120s
```

4. Verify the installation:

```sh
kubectl get svc -n ingress-nginx
```

## BUG
