package main

import (
	"time"

	mgo "gopkg.in/mgo.v2"
)

// Message is a representation of a chat message
type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

// ChatMessage is a representation of a message in chat
type ChatMessage struct {
	ChannelID int       `json:"channelId"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"createdAt"`
	ID        int       `json:"id"`
}

// ChatChannel is a representation of a channel in chat
type ChatChannel struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

// User is a representation of a chat user
type User struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

// ChannelModel is a Channel schema in MongoDB
// type ChannelModel struct {
// 	Name string `bson:"name"`
// 	ID   int    `bson:"id"`
// }

// MessageModel is a Message schema in MongoDB
type MessageModel struct {
	ChannelID int       `bson:"channelId"`
	Content   string    `bson:"content"`
	Author    string    `bson:"author"`
	CreatedAt time.Time `bson:"createdAt"`
	ID        int       `bson:"id"`
}

var dbm *mgo.Collection

// var dbc *mgo.Collection
