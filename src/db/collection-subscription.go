package db

import (
	"errors"
	"github.com/kr/pretty"
	"logger"
	. "types"
)

type Subscription struct {
	Key     string
	Mask    ObjectFields
	Clients []*Client
}

func (c *Collection) AddSubscription(p AddSubscription) (string, error) {
	if _, ok := c.Subscriptions[p.Key]; ok {
		return "", errors.New("Subscription already exists")
	}
	c.Subscriptions[p.Key] = &Subscription{
		Key:  p.Key,
		Mask: p.Mask,
	}
	return "Subscription created", nil
}

func (c *Collection) CancelSubscription(p CancelSubscription) (string, error) {
	if _, ok := c.Subscriptions[p.Key]; !ok {
		return "", errors.New("Subscription does not exist")
	}
	delete(c.Subscriptions, p.Key)
	return "Subscription cancelled", nil
}

func (c *Collection) Subscribe(p Subscribe, client *Client) (string, error) {
	if _, ok := c.Subscriptions[p.Key]; !ok {
		return "", errors.New("Subscription does not exist")
	}
	c.Subscriptions[p.Key].Clients = append(c.Subscriptions[p.Key].Clients, client)
	return "Subscribed using key " + p.Key, nil
}

func (c *Collection) SubscriptionDispatcher(object *ObjectFields) {
	for _, v := range c.Subscriptions {
		if response := v.Match(object); response {
			logger.ErrorLog.Printf("Subscription dispatched: %# v", pretty.Formatter(response))
		}
	}
}

func (s *Subscription) Match(object *ObjectFields) bool {
	return true
}
