package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/mgo.v2/bson"
)

const (
	allCh      = 0
	maxMsgLen  = 420
	maxUsrLen  = 17
	maxChanLen = 17
	maxMsgCnt  = 29
)

var mc = make(chan ChatMessage)  // message chan
var chc = make(chan ChatChannel) // channel chan
var auc = make(chan User)        // add user chan
var euc = make(chan User)        // edit user chan
var ruc = make(chan User)        // remove user chan

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin:     func(r *http.Request) bool { return true },
}

var stopMap = make(map[*websocket.Conn]chan bool)
var chConnMap = make(map[int]map[*websocket.Conn]bool)
var channelMap = make(map[*websocket.Conn]int)
var cidMap = make(map[int]ChatChannel)
var userMap = make(map[*websocket.Conn]User)
var muMap = make(map[*websocket.Conn]*sync.Mutex)

var cid = 1
var uid = 1
var mid = 1

func handleConn(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer conn.Close()
	mu := &sync.Mutex{}
	muMap[conn] = mu
	for {
		var im Message
		if err := conn.ReadJSON(&im); err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %#v\n", im)

		switch im.Name {
		case "message add":
			addMessage(im.Data, conn)
		case "channel add":
			addChannel(im.Data, conn)
		case "user edit":
			editUser(im.Data, conn)
		case "message subscribe":
			go subscribeMessage(im.Data, conn, mu)
		case "message unsubscribe":
			unsubscribeMessage(conn)
		case "channel subscribe":
			go subscribeAddChannel(conn, mu)
		case "user subscribe":
			subscribeUser(conn, mu)
		}
	}
	ru := userMap[conn] // removed user
	delete(userMap, conn)
	delete(muMap, conn)
	removeUser(ru)
	unsubscribeMessage(conn)
	fmt.Println("USER EXITS!", len(userMap))
}

func addMessage(d interface{}, conn *websocket.Conn) {
	var m ChatMessage
	err := mapstructure.Decode(d, &m)
	if err != nil {
		log.Println("mapstructure err:", err)
		return
	}
	m.Author = userMap[conn].Name
	m.CreatedAt = time.Now()
	m.ID = mid
	mid++
	err = dbm.Insert(&MessageModel{
		m.ChannelID,
		m.Content,
		m.Author,
		m.CreatedAt,
		m.ID,
	})
	if err != nil {
		log.Println(err)
	}
	if len(m.Content) > maxMsgLen {
		m.Content = m.Content[:maxMsgLen]
	}
	fmt.Printf("%#v\n", m)
	mc <- m
}

func addChannel(d interface{}, conn *websocket.Conn) {
	var c ChatChannel
	err := mapstructure.Decode(d, &c)
	if err != nil {
		log.Println("mapstructure err:", err)
		return
	}
	c.ID = cid
	cid++
	if len(c.Name) > maxChanLen {
		c.Name = c.Name[:maxChanLen]
	}

	fmt.Printf("%#v\n", c)
	cidMap[c.ID] = c
	chConnMap[c.ID] = make(map[*websocket.Conn]bool)
	chc <- c
}

func addUser(conn *websocket.Conn) {
	u := User{"anon" + string(uid), uid}
	uid++
	userMap[conn] = u
	auc <- u
}

func editUser(d interface{}, conn *websocket.Conn) {
	// fmt.Printf("data passed to editUser: %#v\n", d)
	var data struct {
		CurrentUserName string `json:"currentUserName"`
	}
	err := mapstructure.Decode(d, &data)
	if err != nil {
		log.Println("mapstructure err:", err)
		return
	}
	fmt.Printf("Edited to user: %#v\n", data)
	u := User{data.CurrentUserName, userMap[conn].ID}
	if len(u.Name) > maxUsrLen {
		u.Name = u.Name[:maxUsrLen]
	}
	userMap[conn] = u
	euc <- u
}

func removeUser(ru User) {
	ruc <- ru
}

func subscribeMessage(d interface{}, conn *websocket.Conn, mu *sync.Mutex) {
	var data struct {
		ChannelID int `json:"channelId"`
	}
	if err := mapstructure.Decode(d, &data); err != nil {
		log.Println("mapstructure err:", err)
		return
	}
	fmt.Printf("Subscribed to channel id: %d\n", data.ChannelID)
	channelMap[conn] = data.ChannelID
	chConnMap[data.ChannelID][conn] = true
	stop := make(chan bool)
	stopMap[conn] = stop

	var result []MessageModel
	if err := dbm.Find(bson.M{"channelId": data.ChannelID}).Sort("-createdAt").Limit(maxMsgCnt).All(&result); err != nil {
		log.Println(err)
		return
	}
	sort.Slice(result, func(a int, b int) bool {
		return result[a].CreatedAt.Before(result[b].CreatedAt)
	})
	for i := 0; i < len(result); i++ {
		res := result[i]
		fmt.Printf("res: %+v\n", res)
		om := Message{"message add",
			ChatMessage{
				res.ChannelID,
				res.Content,
				res.Author,
				res.CreatedAt,
				res.ID,
			},
		}
		mu.Lock()
		if err := conn.WriteJSON(om); err != nil {
			log.Println("write error:", err)
			break
		}
		mu.Unlock()
	}

	for {
		select {
		case <-stop:
			fmt.Println("channel stopped")
			return
		case m := <-mc:
			// fmt.Printf("data channel: %d\n", data.ChannelID)
			fmt.Printf("conn channel: %d\n", m.ChannelID)
			om := Message{"message add", m}
			fmt.Printf("add message: %#v\n", om)
			for co := range chConnMap[m.ChannelID] {
				muMap[co].Lock()
				if err := co.WriteJSON(om); err != nil {
					log.Println("write error:", err)
					break
				}
				muMap[co].Unlock()
			}
		}
	}
}

func unsubscribeMessage(conn *websocket.Conn) {
	currCid := channelMap[conn]
	delete(chConnMap[currCid], conn)
	delete(channelMap, conn)
	if stop, ok := stopMap[conn]; ok {
		fmt.Printf("stop chan\n")
		delete(stopMap, conn)
		stop <- true
	}
}

func subscribeAddChannel(conn *websocket.Conn, mu *sync.Mutex) {
	for cid := range chConnMap {
		ch := cidMap[cid]
		om := Message{"channel add", ch}
		mu.Lock()
		if err := conn.WriteJSON(om); err != nil {
			log.Println("write error:", err)
			break
		}
		mu.Unlock()
	}
	for c := range chc {
		om := Message{"channel add", c}
		fmt.Printf("add channel: %#v\n", om)
		for co := range userMap {
			muMap[co].Lock()
			if err := co.WriteJSON(om); err != nil {
				log.Println("write error:", err)
				break
			}
			muMap[co].Unlock()
		}
	}
}

func subscribeUser(conn *websocket.Conn, mu *sync.Mutex) {
	for _, usr := range userMap {
		om := Message{"user add", usr}
		mu.Lock()
		if err := conn.WriteJSON(om); err != nil {
			log.Println("write error:", err)
			break
		}
		mu.Unlock()
	}

	go subscribeAddUser(conn, mu)
	go subscribeEditUser(conn, mu)
	go subscribeRemoveUser(conn, mu)

	addUser(conn)
}

func subscribeAddUser(conn *websocket.Conn, mu *sync.Mutex) {
	for au := range auc {
		om := Message{"user add", au}
		fmt.Printf("add user: %#v\n", om)
		for co := range userMap {
			muMap[co].Lock()
			if err := co.WriteJSON(om); err != nil {
				log.Println("write error:", err)
				break
			}
			muMap[co].Unlock()
		}
	}
}

func subscribeEditUser(conn *websocket.Conn, mu *sync.Mutex) {
	for eu := range euc {
		om := Message{"user edit", eu}
		fmt.Printf("edit user: %#v\n", om)
		for co := range userMap {
			muMap[co].Lock()
			if err := co.WriteJSON(om); err != nil {
				log.Println("write error:", err)
				break
			}
			muMap[co].Unlock()
		}
	}
}

func subscribeRemoveUser(conn *websocket.Conn, mu *sync.Mutex) {
	for ru := range ruc {
		om := Message{"user remove", ru}
		fmt.Printf("remove user: %#v\n", om)
		for co := range userMap {
			muMap[co].Lock()
			if err := co.WriteJSON(om); err != nil {
				log.Println("write error:", err)
				break
			}
			muMap[co].Unlock()
		}
	}
}
