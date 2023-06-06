package main

import (
	"flag"
	"log"
	"time"

	"github.com/nsqio/go-nsq"
)

var (
	nsqdTCPAddress            = flag.String("nsqd-tcp-address", "", "nsqd TCP address")
	destinationNSQDTCPAddress = flag.String("destination-nsqd-tcp-address", "", "destination nsqd TCP address")
	topic                     = flag.String("topic", "", "nsq topic")
	channel                   = flag.String("channel", "nsq_to_nsq", "nsq channel")
	destinationTopic          = flag.String("destination-topic", "", "use this destination topic for all consumed topics")
	maxInFlight               = flag.Int("max-in-flight", 200, "max number of messages to allow in flight")
)

func main() {
	flag.Parse()

	cfg := nsq.NewConfig()
	cfg.MaxInFlight = *maxInFlight

	consumer, err := nsq.NewConsumer(*topic, *channel, cfg)
	if err != nil {
		log.Fatal(err)
	}

	producer, err := nsq.NewProducer(*destinationNSQDTCPAddress, cfg)
	if err != nil {
		log.Fatal(err)
	}

	consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		dstTopic := *destinationTopic
		if dstTopic == "" {
			dstTopic = *topic
		}
		if err := producer.Publish(dstTopic, message.Body); err != nil {
			log.Printf("Could not publish message: %v", err)
			return err
		}
		return nil
	}))

	err = consumer.ConnectToNSQD(*nsqdTCPAddress)
	if err != nil {
		log.Fatal(err)
	}

	for {
		time.Sleep(5 * time.Second)
		if consumer.Stats().MessagesReceived == consumer.Stats().MessagesFinished {
			producer.Stop()
			consumer.Stop()
			break
		}
	}

	<-consumer.StopChan
}
