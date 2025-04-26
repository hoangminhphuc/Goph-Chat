package ws

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	serviceHub "github.com/hoangminhphuc/goph-chat/boot"
	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/common/models"
	ws "github.com/hoangminhphuc/goph-chat/internal/server/websocket"
)

var upgrader = websocket.Upgrader {
	CheckOrigin: func(r *http.Request) bool { return true },
}


func HandleWebSocketConnection(serviceCtx serviceHub.ServiceHub) func(*gin.Context) {
	return func(c *gin.Context) {
		wsServer := serviceCtx.MustGetRuntimeService(common.PluginWSMain).(*ws.WebSocketServer)
		
		roomID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			common.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		
		room, err := wsServer.GetRoom(roomID)
		if err != nil {
			common.ErrNotFound("room", roomID)
			return
		}
		
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return 
		}

		currentUser := c.MustGet(common.CurrentUser).(*models.Requester)

		
		client := ws.NewClient(currentUser.GetUserId(), conn, room)
		
		room.Register <- client

		requestBody := make(chan []byte) 

		go client.Read(c, requestBody)
	}

}