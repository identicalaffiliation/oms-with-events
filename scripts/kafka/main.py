from confluent_kafka.admin import AdminClient, NewTopic
import time

KAFKA_BOOTSTRAP_SERVERS = "kafka:9092"

topicInfo = {"topicName": "orders", "partitions": 3, "replicationFactor": 1}

TOPICS = [
  (topicInfo["topicName"], topicInfo["partitions"], topicInfo["replicationFactor"]),
]

retryCount = 20
sleepTime = 2

def wait_kafka():
  admin = AdminClient({"bootstrap.servers": KAFKA_BOOTSTRAP_SERVERS})

  for _ in range(retryCount):
    try:
      admin.list_topics(timeout=3)
      print("Kafka is ready")
      return
    except Exception:
      time.sleep(sleepTime)

  raise Exception("Kafka not ready")


def create_topics():
  admin = AdminClient({"bootstrap.servers": KAFKA_BOOTSTRAP_SERVERS})

  new_topics = [
      NewTopic(topic=name, num_partitions=p, replication_factor=r)
      for name, p, r in TOPICS
  ]

  fs = admin.create_topics(new_topics)

  for topic, f in fs.items():
    try:
      f.result()
      print(f"Created topic: {topic}")
    except Exception as e:
      print(f"Failed: {topic} -> {e}")


if __name__ == "__main__":
  wait_kafka()
  create_topics()