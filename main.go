package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"flag"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"os"
	"time"
)

var (
	create       bool
	projectID    string
	topic        string
	subscription string
	publisher    bool
	pbTopic      *pubsub.Topic
	pbSubs       *pubsub.Subscription
)

// a simple pubsub client that creates topic and subscription and publishes or receive messages to/from it
func main() {
	flag.BoolVar(&create, "create", false, "If set the topic will be created")
	flag.StringVar(&projectID, "projectId", "", "google cloud projectId")
	flag.StringVar(&topic, "topic", "", "pubsub topic")
	flag.StringVar(&subscription, "subscription", "", "pubsub subscription")
	flag.BoolVar(&publisher, "publisher", false, "If set to True program publishes console inputs on the topic, otherwise listens to the topic")
	flag.Parse()

	if projectID == "" || topic == "" || subscription == "" {
		fmt.Println("Usage:")
		fmt.Println("gopubsub [-create] -projectId=<projectID> -topic=<topic> -subscription=<subscription> [-publisher]")
		os.Exit(1)
	}
	ctx := context.Background()
	pbClient, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile("credentials.json"))
	if err != nil {
		log.Fatalf("Error creating a new pubsub client - %s", err)
	}
	if create {
		pbTopic, err = pbClient.CreateTopic(ctx, topic)
		if err != nil {
			log.Fatalf("Error creatng topic - %s", err)
		}
		cfg := pubsub.SubscriptionConfig{
			Topic:       pbTopic,
			AckDeadline: 10 * time.Second,
		}
		pbSubs, err = pbClient.CreateSubscription(ctx, subscription, cfg)
		if err != nil {
			log.Fatalf("Error creating subscription %s", err)
		}
		log.Println("Topic and subscription are created")
		os.Exit(0)
	}

	pbTopic = pbClient.Topic(topic)
	pbSubs = pbClient.Subscription(subscription)

	if publisher {
		msg := ""
		defer pbTopic.Stop()
		for {
			fmt.Print("message: ")
			l, err := fmt.Scanf("%s", &msg)
			if l == 0 {
				log.Println("Finished publishing messages")
				os.Exit(0)
			}
			if err != nil {
				log.Fatalf("Error reading input - %s", err)
			}
			pubRes := pbTopic.Publish(ctx, &pubsub.Message{Data: []byte(msg)})
			_, err = pubRes.Get(ctx)
			if err != nil {
				log.Fatalf("Error publishing message - %s", err)
			}
			log.Println("message sent!")
		}
	} else {
		err := pbSubs.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
			log.Printf("Got message: %s", m.Data)
			m.Ack()
		})
		if err != nil {
			log.Fatalf("Error receiving messages - %s", err)
		}
	}
}
