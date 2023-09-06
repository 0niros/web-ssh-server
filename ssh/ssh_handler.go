package ssh

import (
	"bufio"
	"errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"net"
	"unicode/utf8"
	"web-ssh-server/config"
)

type WebSshHandler interface {
	Conn(sessionId string, config config.SshConfig) error
	Write(cmd string) error
	Read() (*chan []byte, error)
	Close() error
}

type WebSshHandleService struct {
	HostName    string
	Password    string
	Address     string
	Port        string
	EmptyPasswd bool
	client      *ssh.Client
	channel     *ssh.Channel
}

type ptyRequestMsg struct {
	Term     string
	Columns  uint32
	Rows     uint32
	Width    uint32
	Height   uint32
	ModeList string
}

func (s *WebSshHandleService) Write(cmd string) error {
	if *s.channel == nil {
		logrus.Info("[Write] error case of nil channel: ", cmd)
		return errors.New("error case of nil channel")
	}
	_, err := (*s.channel).Write([]byte(cmd))
	logrus.Info("[Write] ", cmd)

	return err
}

func (s *WebSshHandleService) Read() (*chan []byte, error) {
	if *s.channel == nil {
		logrus.Info("[Write] error case of nil channel")
		return nil, errors.New("error case of nil channel")
	}
	reader := bufio.NewReader(*s.channel)
	msgChan := make(chan []byte)
	go func() {
		for {
			runeRead, _, err := reader.ReadRune()
			if err != nil {
				logrus.Info("[Reader] read error: ", err)
				return
			}
			p := make([]byte, utf8.RuneLen(runeRead))
			utf8.EncodeRune(p, runeRead)
			msgChan <- p
		}
	}()

	return &msgChan, nil
}

func (s *WebSshHandleService) Close() error {
	err := (*s.channel).Close()
	if err != nil {
		return err
	}
	err = s.client.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *WebSshHandleService) Conn(sessionId string, config config.SshConfig) error {
	// 1. Create ssh client.
	auth := make([]ssh.AuthMethod, 0)
	if !s.EmptyPasswd {
		logrus.Info("[Login] non empty passwd.")
		auth = append(auth, ssh.Password(s.Password))
	}
	client, err := ssh.Dial("tcp", s.Address+":"+s.Port, &ssh.ClientConfig{
		User: s.HostName,
		Auth: auth,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	})
	if err != nil {
		logrus.Error("[Conn] conn error for ", sessionId, err)
		return err
	}
	s.client = client

	// 2. Create ssh channel.
	// 2.1 Open channel.
	channel, signal, err := client.OpenChannel("session", nil)
	if err != nil {
		logrus.Errorf("[Channel] open channel error: %s", err)
		return err
	}
	// 2.2 Reply.
	go func() {
		for req := range signal {
			if req.WantReply {
				err = req.Reply(false, nil)
				if err != nil {
					logrus.Info("[Reply] reply error ", err)
				}
			}
		}
	}()

	// 3. Create ssh client mode.
	mode := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 115200,
		ssh.TTY_OP_OSPEED: 115200,
		ssh.ECHOE:         1,
		ssh.IEXTEN:        1,
		ssh.ISIG:          1,
	}
	var modeList []byte
	for k, v := range mode {
		kv := struct {
			Key byte
			Val uint32
		}{k, v}
		modeList = append(modeList, ssh.Marshal(&kv)...)
	}
	modeList = append(modeList, 0)
	req := ptyRequestMsg{
		Term:     "xterm",
		Columns:  config.Cols,
		Rows:     config.Rows,
		Width:    uint32(32 * 8),
		Height:   uint32(160 * 8),
		ModeList: string(modeList),
	}

	// 4. Create terminal.
	// 4.1 Send pty req request.
	ok, err := channel.SendRequest("pty-req", true, ssh.Marshal(&req))
	if !ok || err != nil {
		logrus.Info("[SSH] send pty-req request error ", err)
		return err
	}
	// 4.2 Send shell request.
	ok, err = channel.SendRequest("shell", true, nil)
	if !ok || err != nil {
		logrus.Info("[SSH] send shell request error ", err)
		return err
	}
	s.channel = &channel
	logrus.Info("[Channel] channel created.")

	return err
}

func New(addr string, port string, hostName string, password string) WebSshHandler {
	return &WebSshHandleService{
		HostName: hostName,
		Password: password,
		Address:  addr,
		Port:     port,
	}
}

func NewWithCfg(cfg config.SshConfig) WebSshHandler {
	return &WebSshHandleService{
		HostName:    cfg.Hostname,
		Password:    cfg.Password,
		Address:     cfg.Address,
		Port:        cfg.Port,
		EmptyPasswd: len(cfg.Password) == 0,
	}
}
