package connection

import (
	"bytes"
	"github.com/pkg/errors"
	"log"
	"net"
	"sync"
	"time"
)

type RoboMasterConn struct {
	option *Option
	IPConn net.PacketConn

	cmdLock *sync.Mutex

	roboIp      net.IP
	CtrlConn    net.Conn
	ctrlRecChan chan []byte

	VideoConn net.Conn
	AudioConn net.Conn
	PushConn  net.Conn
	EventConn net.Conn
}

func NewRoboMasterConn(option *Option) (*RoboMasterConn, error) {
	option = getDefaultOption(option)
	roboMaserConn := &RoboMasterConn{
		option:      option,
		ctrlRecChan: make(chan []byte, 1024),
		cmdLock:     new(sync.Mutex),
	}

	if option.IP == "" {
		err := roboMaserConn.scanRoboIp()
		if err != nil {
			return nil, err
		}
	} else {
		roboMaserConn.roboIp = net.ParseIP(option.IP)
	}

	if err := roboMaserConn.initCtrlConn(); err != nil {
		return nil, err
	}

	if err := roboMaserConn.dialConns(); err != nil {
		return nil, err
	}

	return roboMaserConn, nil
}

func (r *RoboMasterConn) scanRoboIp() (err error) {
	ipAddr := net.UDPAddr{Port: int(IP_PORT)}
	r.IPConn, err = net.ListenPacket(UDP, ipAddr.String())
	if err != nil {
		return errors.New("listen ip broadcast port error")
	}
	err = r.reciveBroadcastPack()
	return
}
func (r *RoboMasterConn) reciveBroadcastPack() error {
	r.IPConn.SetDeadline(time.Now().Add(r.option.ScanTimeout))
	buff := make([]byte, 64)
	for {
		n, addr, err := r.IPConn.ReadFrom(buff)
		if err != nil {
			return errors.Wrap(err, "recive ip broadcast error")
		}
		reciveIP := string(buff[9:n])

		udpIP, _ := net.ResolveUDPAddr(UDP, addr.String())
		log.Print("recive robomaster ip ", udpIP.IP.String())
		if udpIP.IP.String() == reciveIP {
			r.roboIp = udpIP.IP
			return nil
		}
	}
	return errors.New("recive ip broadcast error")
}
func (r *RoboMasterConn) initCtrlConn() (err error) {
	r.CtrlConn, err = r.dialConn(CTRL_NETWORK, CTRL_PORT)
	if err != nil {
		return errors.Wrap(err, "dial ctrl conn fail")
	}
	r.CtrlConn.SetReadDeadline(time.Now().Add(time.Hour))
	go func() {
		buff := make([]byte, 1024)
		for {
			n, err := r.CtrlConn.Read(buff)
			if err != nil {
				log.Fatal(err)
			}
			log.Print("recive command result:", string(buff[0:n]))
			r.ctrlRecChan <- buff[0:n]
		}
	}()
	return nil
}

func (r *RoboMasterConn) dialConns() (err error) {

	cmd, err := r.runCtrlCommnd("command")
	if err != nil || cmd != "ok" {
		return errors.New("join sdk command error")
	}
	if r.option.EnableVideo {
		err = r.EnableVideo()
		if err != nil {
			return
		}
	}
	if r.option.EnableAudio {
		err = r.EnableAudio()
		if err != nil {
			return
		}
	}
	r.PushConn, err = r.dialConn(PUSH_NETWORK, PUSH_PORT)
	if err != nil {
		return errors.Wrap(err, "dial push conn fail")
	}
	r.EventConn, err = r.dialConn(EVENT_NETWORK, EVENT_PORT)
	if err != nil {
		return errors.Wrap(err, "dial event conn fail")
	}
	if err != nil {
		return errors.Wrap(err, "dial ip conn fail")
	}

	return nil
}

func (r *RoboMasterConn) EnableVideo() (err error) {
	rec, err := r.runCtrlCommnd(EnableVideo)
	if err != nil || rec != "ok" {
		return errors.Wrap(err, "enable video stream fail")
	}
	r.VideoConn, err = r.dialConn(VIDEO_NETWORK, VIDEO_PORT)
	return
}
func (r *RoboMasterConn) DisableVideo() (err error) {
	if r.CtrlConn == nil {
		return
	}
	_, err = r.CtrlConn.Write([]byte(DisableVideo))
	if err != nil {
		return errors.Wrap(err, "disable video stream fail")
	}
	r.VideoConn = nil
	return
}
func (r *RoboMasterConn) EnableAudio() (err error) {
	rec, err := r.runCtrlCommnd(EnableAudio)
	if err != nil || rec != "ok" {
		return errors.Wrap(err, "enable audio fail")
	}
	r.AudioConn, err = r.dialConn(AUDIO_NETWORK, AUDIO_PORT)
	return
}
func (r *RoboMasterConn) DisableAudio() (err error) {
	if r.CtrlConn == nil {
		return
	}
	_, err = r.CtrlConn.Write([]byte(DisableAudio))
	if err != nil {
		return errors.Wrap(err, "disable audio fail")
	}
	r.AudioConn = nil
	return
}
func (r *RoboMasterConn) RunCmd(cmd string) (rec string, err error) {
	command := Command(cmd)
	return r.runCtrlCommnd(command)
}

func (r *RoboMasterConn) runCtrlCommnd(command Command) (rec string, err error) {
	r.cmdLock.Lock()
	defer r.cmdLock.Unlock()

	if r.CtrlConn == nil {
		return "", errors.New("must dial ctrl conn")
	}
	cmd := bytes.NewBuffer([]byte(command))

	cmd.Write([]byte(CommandSeparator))

	_, err = r.CtrlConn.Write(cmd.Bytes())
	log.Print("send command:", cmd.String())
	if err != nil {
		return "", err
	}
	//cmd响应
	select {
	case msg := <-r.ctrlRecChan:
		return string(msg), nil
	case <-time.After(r.option.CtrlTimeOut):
		return "", errors.New("rec ctrl command timeout")
	}
	return
}

func (r *RoboMasterConn) dialConn(network NetWork, port Port) (net.Conn, error) {
	var addr string
	if network == TCP {
		tcpAddr := net.TCPAddr{
			IP:   r.roboIp,
			Port: int(port),
		}
		addr = tcpAddr.String()
	} else {
		udpAddr := net.UDPAddr{
			IP:   r.roboIp,
			Port: int(port),
		}
		addr = udpAddr.String()
	}
	return net.Dial(string(network), addr)
}
