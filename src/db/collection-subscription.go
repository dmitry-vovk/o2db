package db

import (
	"errors"
	. "types"
)

func (c *Collection) AddSubscription(p AddSubscription) (string, error) {
	if _, ok := c.Subscriptions[p.Key]; ok {
		return "", errors.New("Subscription already exists")
	}
	c.Subscriptions[p.Key] = p.Mask
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
	return "Subscribed using key " + p.Key, nil
}

func (c *Collection) subscriptionDispatcher(object *ObjectFields) {
	// TODO send copy fo the object to all relevant subscribers
}
