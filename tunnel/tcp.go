package tunnel

import (
	"io"
	"net"
	"sync"
	"time"

	"github.com/xjasonlyu/tun2socks/v2/common/pool"
	"github.com/xjasonlyu/tun2socks/v2/core/adapter"
	"github.com/xjasonlyu/tun2socks/v2/log"
	M "github.com/xjasonlyu/tun2socks/v2/metadata"
	"github.com/xjasonlyu/tun2socks/v2/proxy"
	"github.com/xjasonlyu/tun2socks/v2/tunnel/statistic"
)

const (
	tcpWaitTimeout = 5 * time.Second
)

func newTCPTracker(conn net.Conn, metadata *M.Metadata) net.Conn {
	return statistic.NewTCPTracker(conn, metadata, statistic.DefaultManager)
}

func handleTCPConn(localConn adapter.TCPConn) {
	defer localConn.Close()

	id := localConn.ID()
	log.Debugf("[TCP] handleTCPConn id %v", id)
	metadata := &M.Metadata{
		Network: M.TCP,
		SrcIP:   net.IP(id.RemoteAddress),
		SrcPort: id.RemotePort,
		DstIP:   net.IP(id.LocalAddress),
		DstPort: id.LocalPort,
	}

	targetConn, err := proxy.Dial(metadata)
	if err != nil {
		log.Warnf("[TCP] dial %s error: %v", metadata.DestinationAddress(), err)
		return
	}
	metadata.MidIP, metadata.MidPort = parseAddr(targetConn.LocalAddr())

	// igniter: no statistic
	//targetConn = newTCPTracker(targetConn, metadata)
	defer targetConn.Close()

	log.Infof("[TCP] %s <-> %s", metadata.SourceAddress(), metadata.DestinationAddress())
	relay(localConn, targetConn) /* relay connections */
}

// relay copies between left and right bidirectionally.
func relay(left, right net.Conn) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		_ = copyBuffer(right, left) /* ignore error */
		right.SetReadDeadline(time.Now().Add(tcpWaitTimeout))
	}()

	go func() {
		defer wg.Done()
		_ = copyBuffer(left, right) /* ignore error */
		left.SetReadDeadline(time.Now().Add(tcpWaitTimeout))
	}()

	wg.Wait()
}

func copyBuffer(dst io.Writer, src io.Reader) error {
	buf := pool.Get(pool.RelayBufferSize)
	defer pool.Put(buf)

	_, err := io.CopyBuffer(dst, src, buf)
	return err
}
