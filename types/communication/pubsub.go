package communication

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

const (
	PUBLISH     = "publish"
	SUBSCRIBE   = "subscribe"
	UNSUBSCRIBE = "unsubscribe"
)

type PubSub struct {
	Clients       []PubSubClient
	Subscriptions []Subscription
}

type PubSubResponse struct {
	Message string
	Topic   string
	Sender  string
}

type PubSubClient struct {
	Id         string
	Connection *websocket.Conn
}

type PubSubMessage struct {
	Action  string          `json:"action"`
	Topic   string          `json:"topic"`
	Message json.RawMessage `json:"message"`
}

type Subscription struct {
	Topic        string
	PubSubClient *PubSubClient
}

func (ps *PubSub) AddClient(client PubSubClient) *PubSub {
	ps.Clients = append(ps.Clients, client)
	//fmt.Println("adding new client to the list", client.Id, len(ps.Clients))
	payload := []byte("Hello Client ID:" + client.Id)
	client.Connection.WriteMessage(1, payload)

	return ps
}

func (ps *PubSub) RemoveClient(client PubSubClient) *PubSub {
	// first remove all subscriptions by this client
	for index, sub := range ps.Subscriptions {
		if client.Id == sub.PubSubClient.Id {
			ps.Subscriptions = append(ps.Subscriptions[:index], ps.Subscriptions[index+1:]...)
		}
	}

	// remove client from the list
	for index, c := range ps.Clients {
		if c.Id == client.Id {
			ps.Clients = append(ps.Clients[:index], ps.Clients[index+1:]...)
		}
	}

	return ps
}

func (ps *PubSub) GetSubscriptions(topic string, client *PubSubClient) []Subscription {
	var subscriptionList []Subscription
	for _, subscription := range ps.Subscriptions {
		if client != nil {
			if subscription.PubSubClient.Id == client.Id && subscription.Topic == topic {
				subscriptionList = append(subscriptionList, subscription)

			}
		} else {
			if subscription.Topic == topic {
				subscriptionList = append(subscriptionList, subscription)
			}
		}
	}

	return subscriptionList
}

func (ps *PubSub) Subscribe(client *PubSubClient, topic string) *PubSub {
	clientSubs := ps.GetSubscriptions(topic, client)
	if len(clientSubs) > 0 {
		// client is subscribed this topic before
		return ps
	}
	newSubscription := Subscription{
		Topic:        topic,
		PubSubClient: client,
	}
	fmt.Println(client.Id, topic, client)
	ps.Subscriptions = append(ps.Subscriptions, newSubscription)

	return ps
}

func (ps *PubSub) Publish(topic string, message []byte, excludeClient *PubSubClient) {
	subscriptions := ps.GetSubscriptions(topic, nil)
	for _, sub := range subscriptions {
		fmt.Printf("Sending to client id %s message is %s \n", sub.PubSubClient.Id, message)
		//sub.Client.Connection.WriteMessage(1, message)
		re := PubSubResponse{
			Message: string(message[:]),
			Topic:   topic,
			Sender:  sub.PubSubClient.Id,
		}
		b, _ := json.Marshal(re)
		sub.PubSubClient.Send(b)
	}
}

func (client *PubSubClient) Send(message []byte) error {
	return client.Connection.WriteMessage(1, message)
}

func (ps *PubSub) Unsubscribe(client *PubSubClient, topic string) *PubSub {
	//clientSubscriptions := ps.GetSubscriptions(topic, client)
	for index, sub := range ps.Subscriptions {
		if sub.PubSubClient.Id == client.Id && sub.Topic == topic {
			// found this subscription from client and we do need remove it
			ps.Subscriptions = append(ps.Subscriptions[:index], ps.Subscriptions[index+1:]...)
		}
	}

	return ps
}

func (ps *PubSub) HandleReceiveMessage(client PubSubClient, messageType int, payload []byte) *PubSub {
	m := PubSubMessage{}
	err := json.Unmarshal(payload, &m)
	if err != nil {
		fmt.Println("This is not correct message payload", err)
		return ps
	}

	switch m.Action {
	case PUBLISH:
		fmt.Println("This is publish new message")
		ps.Publish(m.Topic, m.Message, nil)
		break
	case SUBSCRIBE:
		ps.Subscribe(&client, m.Topic)
		fmt.Println("new subscriber to topic", m.Topic, len(ps.Subscriptions), client.Id)
		break
	case UNSUBSCRIBE:
		fmt.Println("Client want to unsubscribe the topic", m.Topic, client.Id)
		ps.Unsubscribe(&client, m.Topic)
		break

	default:
		break
	}

	return ps
}
