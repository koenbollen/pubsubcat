# pubsubcat

The Google Pub/Sub Swiss Army Knife

---

_pubsubcat_ is a developer's utility for using [Google Cloud Pub/Sub]. Use it to
quickly peek into a topic or when you want to manually publish a message.

One of the main features is that this tool automatically creates a new temporary
subscription on topic when you subscribe and cleans it when it's no longer used.

Inspired by [kafkacat] and initially created to get acquainted with Pub/Sub.

[google cloud pub/sub]: https://cloud.google.com/pubsub/
[kafkacat]: https://github.com/edenhill/kafkacat
