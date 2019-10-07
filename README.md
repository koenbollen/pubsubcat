# pubsubcat ![CI](https://github.com/koenbollen/pubsubcat/workflows/CI/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/koenbollen/pubsubcat)](https://goreportcard.com/report/github.com/koenbollen/pubsubcat)
The Google Pub/Sub Swiss Army Knife

---

_pubsubcat_ is a developer's utility for using [Google Cloud Pub/Sub]. Use it to
quickly peek into a topic or when you want to manually publish a message.

One of the main features is that this tool automatically creates a new temporary
subscription on topic when you subscribe and cleans it when it's no longer used.

Inspired by [kafkacat] and initially created to get acquainted with Pub/Sub.

## Installation

```
go get -u github.com/koenbollen/pubsubcat
```

## Usage

```
Usage:
  pubsubcat [flags]
  pubsubcat [command]

Available Commands:
  create      Create topics and subscriptions
  help        Help about any command
  ls          List topics and subscriptions
  pipe        Pipe a topic directly to another using a temporary subscription
  pop         alias for: `pubsubcat subscribe --count 1 TOPIC`
  publish     Publish input lines as messages
  subscribe   Subscribe to a topic using a temporary subscription

Flags:
      --config string    config file (default is ~/.pubsubcat)
  -h, --help             help for pubsubcat
  -p, --project string   Google Cloud Project to work under
  -q, --quiet            only output messages
  -v, --verbose count    increase verbosity

Use "pubsubcat [command] --help" for more information about a command.
```

## Examples

#### Create and list topics:

```
$ pubsubcat create my-topic
projects/bnl-blendle/topics/my-topic
$ pubsubcat ls
projects/bnl-blendle/topics/my-topic
```

_pubsubcat_ can create and list topics and subscriptions with the default configuration.

#### Publish a message:

```
$ date | pubsubcat publish my-topic
] publishing lines from stdin to topic my-topic
] eof, syncing
$
```

_pubsubcat_ read from stdin and publish each line as a message on the given topic.

#### Subscribing to a topic:

```
$ pubsubcat subscribe my-topic
] listening on topic "my-topic" using subscription "pubsubcat-my-topic-4d65822107fcfd52-1530794401"
Thu Jul  5 14:40:06 CEST 2018
^C
] stopped receiving, cleaning up temporary subscription
```

_pubsubcat_ creates a temporary subscription to read from the given topic without consuming messages.

[google cloud pub/sub]: https://cloud.google.com/pubsub/
[kafkacat]: https://github.com/edenhill/kafkacat
