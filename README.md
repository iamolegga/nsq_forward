# nsq_forward

Like nsq_to_nsq but exit after all messages are processed in source channel.

## Usage

```shell
docker run --rm iamolegga/nsq_forward \
  -nsqd-tcp-address source-nsqd-host:4150 \
  -destination-nsqd-tcp-address destination-nsqd-host:4150 \
  -topic source_topic \
  -channel source_channel \
  -destination-topic destination_topic
```

## License

This project is licensed under the MIT License. See the LICENSE file for details.

