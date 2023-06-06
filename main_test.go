package main_test

import (
	"os"
	"os/exec"
	"regexp"
	"testing"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	nsqdTCPAddress      string
	sourceTopic         string
	sourceChannel       string
	destinationTopic    string
	destinationChannel  string
	args                []string
	sourceProducer      *nsq.Producer
	destinationConsumer *nsq.Consumer
	forwardedMessages   chan *nsq.Message
}

func (s *Suite) SetupTest() {
	s.nsqdTCPAddress = "localhost:4150"
	prefix := regexp.MustCompile(`\D`).ReplaceAllString(time.Now().Format(time.RFC3339), "")
	s.sourceTopic = "source_topic_" + prefix
	s.sourceChannel = "source_channel_" + prefix
	s.destinationTopic = "destination_topic_" + prefix
	s.destinationChannel = "destination_channel_" + prefix

	s.args = append(s.args, "-nsqd-tcp-address", s.nsqdTCPAddress)
	s.args = append(s.args, "-destination-nsqd-tcp-address", s.nsqdTCPAddress)
	s.args = append(s.args, "-topic", s.sourceTopic)
	s.args = append(s.args, "-channel", s.sourceChannel)
	s.args = append(s.args, "-destination-topic", s.destinationTopic)

	// setup source producer
	config := nsq.NewConfig()
	var err error
	s.sourceProducer, err = nsq.NewProducer(s.nsqdTCPAddress, config)
	s.Nil(err, "failed to create nsq source producer")

	// setup destination consumer
	s.destinationConsumer, err = nsq.NewConsumer(s.destinationTopic, s.destinationChannel, config)
	s.Nil(err, "failed to create destination nsq consumer")

	s.forwardedMessages = make(chan *nsq.Message)
	s.destinationConsumer.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		s.forwardedMessages <- m
		return nil
	}))

	err = s.destinationConsumer.ConnectToNSQD(s.nsqdTCPAddress)
	s.Nil(err, "failed to connect to nsqd")
}

func (s *Suite) TearDownTest() {
	s.sourceProducer.Stop()
}

func (s *Suite) TestSingleMessage() {
	// Publish a single message
	testMessage := []byte("test")
	err := s.sourceProducer.Publish(s.sourceTopic, testMessage)
	s.Nil(err, "failed to publish message")

	s.runNsqForward()

	// Wait for the forwarded message
	select {
	case message := <-s.forwardedMessages:
		s.Equal(string(message.Body), "test", "forwarded message did not match expected")
	case <-time.After(10 * time.Second):
		s.Fail("timed out waiting for forwarded message")
	}
}

func (s *Suite) TestMultipleMessages() {
	// Publish multiple messages
	messages := []string{"test1", "test2", "test3"}
	for _, msg := range messages {
		err := s.sourceProducer.Publish(s.sourceTopic, []byte(msg))
		s.Nil(err, "failed to publish message")
	}

	s.runNsqForward()

	// Wait for the forwarded messages
	for _, _ = range messages {
		select {
		case message := <-s.forwardedMessages:
			s.Contains(messages, string(message.Body), "forwarded message was not among the expected")
		case <-time.After(10 * time.Second):
			s.Fail("timed out waiting for forwarded message")
		}
	}
}

func (s *Suite) runNsqForward() {
	cmd := exec.Command("./nsq_forward", s.args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	s.Nil(err, "failed to run nsq_forward")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
