package db

import (
	"errors"
	"github.com/kr/pretty"
	"logger"
	. "types"
)

func (c *Collection) AddSubscription(p AddSubscription) (string, uint, error) {
	if _, ok := c.Subscriptions[p.Key]; ok {
		return "", RSubscriptionAlreadyExists, errors.New("Subscription already exists")
	}
	newSubscription := &Subscription{
		Key:   p.Key,
		Query: p.Query,
	}
	if err := newSubscription.Validate(); err != nil {
		return "Invalid subscription format", RSubscriptionInvalidFormat, err
	}
	c.Subscriptions[p.Key] = newSubscription
	return "Subscription created", RSubscriptionCreated, nil
}

func (c *Collection) CancelSubscription(p CancelSubscription) (string, uint, error) {
	if _, ok := c.Subscriptions[p.Key]; !ok {
		return "", RSubscriptionDoesNotExist, errors.New("Subscription does not exist")
	}
	delete(c.Subscriptions, p.Key)
	return "Subscription cancelled", RSubscriptionCancelled, nil
}

func (c *Collection) Subscribe(p Subscribe, client *Client) (string, uint, error) {
	if _, ok := c.Subscriptions[p.Key]; !ok {
		return "", RSubscriptionDoesNotExist, errors.New("Subscription does not exist")
	}
	c.Subscriptions[p.Key].Clients = append(c.Subscriptions[p.Key].Clients, client)
	return "Subscribed using key " + p.Key, RSubscribed, nil
}

func (c *Collection) ListSubscriptions() []SubscriptionItem {
	subscriptions := []SubscriptionItem{}
	for _, s := range c.Subscriptions {
		subscriptions = append(
			subscriptions,
			SubscriptionItem{
				Collection: c.Name,
				Key:        s.Key,
				Query:      s.Query,
			},
		)
	}
	return subscriptions
}

func (c *Collection) SubscriptionDispatcher(object *ObjectFields) {
	for _, v := range c.Subscriptions {
		if response := v.Match(*object); response {
			logger.ErrorLog.Printf("Iterating over clients (%d)", len(v.Clients))
			for _, client := range v.Clients {
				client.Respond(Response{
					Result:   true,
					Response: object,
				})
			}
			logger.ErrorLog.Printf("Subscription dispatched: %# v", pretty.Formatter(response))
		}
	}
}
