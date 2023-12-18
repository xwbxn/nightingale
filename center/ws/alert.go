package ws

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/toolkits/pkg/ginx"
)

var (
	clients = make(map[uint]map[*websocket.Conn]bool)
	mux     sync.Mutex
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 20 * time.Second,
	// 取消 ws 跨域校验
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WsHandler 处理ws请求
func WsHandler(c *gin.Context) {
	var conn *websocket.Conn
	var err error

	idStr := ginx.UrlParamStr(c, "id")
	u64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		fmt.Println(err)
		return
	}
	id := uint(u64)

	conn, err = upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}

	// var json struct {
	// 	Jwt string `json:"jwt"`
	// }

	// err = conn.ReadJSON(&json)

	// if err != nil {
	// 	log.Println("read json err:", err)
	// 	_ = conn.Close()
	// 	return
	// }

	// // 验证逻辑
	// var userId uint
	// userId, err = model.ValidateJWT(json.Jwt)

	// if err != nil {
	// 	log.Println("validate jwt err:", err)
	// 	_ = conn.Close()
	// 	return
	// }

	// // 将 conn 保存到字典  以上注释部分暂时不用，以后根据分配的Id来完善消息推送
	addClient(id, conn)
}

func addClient(id uint, conn *websocket.Conn) {

	mux.Lock()
	if clients[id] == nil {
		clients[id] = make(map[*websocket.Conn]bool)
	}
	clients[id][conn] = true
	mux.Unlock()
}

func getClients(id uint) (conns []*websocket.Conn) {
	mux.Lock()
	_conns, ok := clients[id]
	if ok {
		for k := range _conns {
			conns = append(conns, k)
		}
	}
	mux.Unlock()
	return
}

func deleteClient(id uint, conn *websocket.Conn) {
	mux.Lock()
	_ = conn.Close()
	delete(clients[id], conn)
	mux.Unlock()
}

type Data struct {
	Err string      `json:"err"`
	Dat interface{} `json:"dat"`
}

func SetMessage(userId uint, content interface{}) {
	mux.Lock()
	conns := getClients(userId)
	for i := range conns {
		i := i

		data := Data{"", content}
		err := conns[i].WriteJSON(data)

		if err != nil {
			log.Println("write json err:", err)
			deleteClient(userId, conns[i])
		}

	}
	mux.Unlock()
}
