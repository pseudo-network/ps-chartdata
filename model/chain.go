package model

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chain struct {
	Id   primitive.ObjectID `json:"id,omitempty" bson:"id"`
	Name string             `json:"name" bson:"name"`
	// todo...
}

// placeholder

func GetChainByID(chainName string) (string, error) {
	if chainName == "1" {
		return "ethereum", nil
	}
	if chainName == "2" {
		return "bsc", nil
	}
	if chainName == "3" {
		return "cardano", nil
	}

	// todo:
	err := errors.New("chain not found")
	return "", err
}
