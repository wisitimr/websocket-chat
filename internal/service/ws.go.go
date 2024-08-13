package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"websocket-chat/internal/model"
	mRepo "websocket-chat/internal/model/repository"
	mService "websocket-chat/internal/model/service"

	"github.com/sirupsen/logrus"
)

type wsService struct {
	redisRepo mRepo.RedisRepository
	logger    *logrus.Logger
	broadcast chan *model.Chat
}

func InitWsService(repo mRepo.Repository, logger *logrus.Logger, broadcast chan *model.Chat) mService.WsService {
	return &wsService{
		redisRepo: repo.Redis,
		logger:    logger,
		broadcast: broadcast,
	}
}

// define a receiver which will listen for
// new messages being sent to our WebSocket
// endpoint
func (s wsService) Receiver(client *model.Client) {
	for {
		// read in a message
		// readMessage returns messageType, message, err
		// messageType: 1-> Text Message, 2 -> Binary Message
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		m := &model.Message{}

		err = json.Unmarshal(p, m)
		if err != nil {
			log.Println("error while unmarshaling chat", err)
			continue
		}

		fmt.Println("host", client.Conn.RemoteAddr())
		if m.Type == "bootup" {
			// do mapping on bootup
			client.Username = m.User
			fmt.Println("client successfully mapped", &client, client, client.Username)
		} else {
			fmt.Println("received message", m.Type, m.Chat)
			c := m.Chat
			c.Timestamp = time.Now().Unix()

			// save in redis
			id, err := s.redisRepo.CreateChat(&c)
			if err != nil {
				log.Println("error while saving chat in redis", err)
				return
			}

			c.ID = id
			s.broadcast <- &c
		}
	}
}
