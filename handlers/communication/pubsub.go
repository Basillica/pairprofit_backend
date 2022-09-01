package communication

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"pairprofit.com/x/types/communication"
)

func autoId() string {
	var err error
	return uuid.Must(uuid.New(), err).String()
}

var ps = &communication.PubSub{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// check origin will check the cross region source (note : please not using in production)
	CheckOrigin: func(r *http.Request) bool {
		//Here we just allow the chrome extension client accessable (you should check this verify accourding your client source)
		return true //origin == "chrome-extension://cbcbkhdmedgianpaifchdaddpnmgnknn"
	},
}

func WebsocketHandler(c *gin.Context) {
	id := c.Param("id")
	// verify id from application OR NOT
	fmt.Println("User id is: ", id)

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("An error occured: ", err)
		return
	}
	client := communication.PubSubClient{
		Id:         autoId(),
		Connection: conn,
	}
	// add this client into the list
	ps.AddClient(client)
	fmt.Println("New Client is connected, total: ", len(ps.Clients))

	defer conn.Close()
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Something went wrong", err)
			ps.RemoveClient(client)
			fmt.Println("total clients and subscriptions ", len(ps.Clients), len(ps.Subscriptions))
			return
		}
		ps.HandleReceiveMessage(client, messageType, message)
	}
}
