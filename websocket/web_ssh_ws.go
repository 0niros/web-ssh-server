package websocket

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
	"web-ssh-server/config"
	"web-ssh-server/response"
	"web-ssh-server/ssh"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(*http.Request) bool {
			// No need to check origin here, return true directly.
			return true
		},
	}
	connPool = make(map[string]*websocket.Conn)
	mux      sync.Mutex
)

func WebSshWsHandler(c *gin.Context) {
	// 1. Generate session id.
	sessionId := c.Param("sessionId")
	if len(sessionId) == 0 {
		logrus.Error("[AUTH] sessionId is null.")
		response.ErrorStatusHandler(c, http.StatusForbidden, 403, "session id is null")
		return
	}

	// 2. Update http conn.
	updateConn(c.Writer, c.Request, sessionId)
}

func updateConn(w http.ResponseWriter, r *http.Request, id string) {
	// 1. Do upgrade.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Error("[Upgrader] upgrade error\n", err)
		return
	}
	messageChannel := addClient(id, conn)

	// 2. Ticker.
	pingTicker := time.NewTicker(time.Second * 20)

	// 3. Parse config file.
	_, msg, err := conn.ReadMessage()
	if err != nil {
		logrus.Error("[Parse] read ssh config error\n", err)
		return
	}
	sshConfig := config.SshConfig{}
	err = json.Unmarshal(msg, &sshConfig)
	if err != nil {
		logrus.Error("[Parse] parse ssh config error\n", err)
		return
	}
	logrus.Infof("[Config]: {%s}:{%s} {%s} {%s}", sshConfig.Address, sshConfig.Port, sshConfig.Hostname, sshConfig.Password)

	// 4. Connect ssh.
	sshHandler := ssh.NewWithCfg(sshConfig)
	if err = sshHandler.Conn(id, sshConfig); err != nil {
		logrus.Errorf("[SSH] login ssh error: \n %s", err)
		return
	}

	// 5. Create reader thread and writer thread.
	go readFromWs(conn, messageChannel)
	go startReader(messageChannel, &sshHandler, conn)
	go startWriter(&sshHandler, conn)
	go startPing(pingTicker, conn)
}

func startReader(m *chan string, h *ssh.WebSshHandler, c *websocket.Conn) {
	for {
		select {
		case msg, _ := <-*m:
			err := (*h).Write(msg)
			if err != nil {
				logrus.Errorf("[WsReader] write cmd [%s] error: %s\n", msg, err)
				stop(c, h)
				return
			}
		}
	}
}

func startWriter(h *ssh.WebSshHandler, c *websocket.Conn) {
	read, err := (*h).Read()
	buffer := make([]byte, 0)
	t := time.NewTimer(time.Microsecond * 50)
	if err != nil {
		logrus.Errorf("[WsWriter] read tty {%s} error.", err)
		stop(c, h)
		return
	}
	for {
		select {
		case msg := <-*read:
			buffer = append(buffer, msg...)
		case <-t.C:
			if len(buffer) != 0 {
				err = c.WriteMessage(websocket.TextMessage, buffer)
				if err != nil {
					logrus.Error("[WS] send msg error")
					stop(c, h)
				}
				buffer = []byte{}
			}
			t.Reset(time.Microsecond * 50)
		}
	}
}

func startPing(t *time.Ticker, c *websocket.Conn) {
	for {
		select {
		case <-t.C:
			err := c.SetWriteDeadline(time.Now().Add(time.Second * 25))
			if err != nil {
				logrus.Error("[PING] ws ping error\n", err)
				return
			}
			err = c.WriteMessage(websocket.PingMessage, []byte{})
		}
	}
}

func stop(c *websocket.Conn, s *ssh.WebSshHandler) {
	err := c.Close()
	if err != nil {
		logrus.Info("[WS] stop ws error: ", err)
		return
	}
	err = (*s).Close()
	if err != nil {
		logrus.Info("[SSH] stop ssh error: ", err)
		return
	}
}

func addClient(id string, conn *websocket.Conn) *chan string {
	mux.Lock()
	defer mux.Unlock()
	connPool[id] = conn
	m := make(chan string)
	return &m
}

func readFromWs(c *websocket.Conn, m *chan string) {
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			logrus.Error("[Read] read from ws error", err)
			return
		}
		ms := string(msg)
		*m <- ms
	}
}
