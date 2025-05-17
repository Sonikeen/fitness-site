package storage

import "fitness-site/internal/models"

var subscriptions []model.Subscription
var subscriptionID = 1

func GetAllSubscriptions() []model.Subscription {
	return subscriptions
}

func AddSubscription(sub model.Subscription) {
	sub.ID = subscriptionID
	subscriptionID++
	subscriptions = append(subscriptions, sub)
}
