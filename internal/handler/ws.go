package handler

import (
	"fmt"
	"log"
	"net/http"
	"websocket-chat/internal/model"
	mHandler "websocket-chat/internal/model/handler"
	mRes "websocket-chat/internal/model/response"
	mService "websocket-chat/internal/model/service"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// We'll need to check the origin of our connection
	// this will allow us to make requests from our React
	// development server to here.
	// For now, we'll do no checking and just allow any connection
	CheckOrigin: func(r *http.Request) bool { return true },
}

type wsHandler struct {
	wsService mService.WsService
	logger    *logrus.Logger
	clients   map[*model.Client]bool
	mRes.ResponseDto
}

func InitWsHandler(service mService.Service, logger *logrus.Logger) mHandler.WsHandler {
	return wsHandler{
		wsService: service.WebSocket,
		logger:    logger,
	}
}

func (h wsHandler) WebSocketConnect(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host, r.URL.Query())

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	client := &model.Client{Conn: ws}
	// register client
	h.clients[client] = true
	fmt.Println("clients", len(h.clients), h.clients, ws.RemoteAddr())

	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	h.wsService.Receiver(client)

	fmt.Println("exiting", ws.RemoteAddr().String())
	delete(h.clients, client)
}
