package upstream

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/jackc/puddle/v2"
)

type DialerFn func(ctx context.Context) (net.Conn, error)

type ConnPool interface {
	Acquire(ctx context.Context) (PoolConn, error)
	Close()
}

type PoolConn interface {
	net.Conn
	Destroy()
}

var _ ConnPool = (*NetPool)(nil)

type NetPool struct {
	dialer DialerFn
}

func NewNetPool(d DialerFn) *NetPool {
	return &NetPool{
		dialer: d,
	}
}

func (p *NetPool) Acquire(ctx context.Context) (PoolConn, error) {
	conn, err := p.dialer(ctx)
	if err != nil {
		return nil, err
	}

	return &netConn{
		Conn: conn,
	}, nil
}

func (p *NetPool) Close() {}

var _ PoolConn = (*netConn)(nil)

type netConn struct {
	net.Conn
}

func (c *netConn) Destroy() {
	_ = c.Close()
}

var _ ConnPool = (*PuddlePool)(nil)

type PuddlePool struct {
	pool       *puddle.Pool[net.Conn]
	maxRetries int
}

func NewPuddlePool(dialer DialerFn, maxSize int32) (*PuddlePool, error) {
	pool, err := puddle.NewPool(&puddle.Config[net.Conn]{
		Constructor: puddle.Constructor[net.Conn](dialer),
		Destructor: func(conn net.Conn) {
			conn.Close()
		},
		MaxSize: maxSize,
	})

	if err != nil {
		return nil, err
	}

	return &PuddlePool{
		pool:       pool,
		maxRetries: 16,
	}, nil
}

func (p *PuddlePool) Acquire(ctx context.Context) (PoolConn, error) {
	for range p.maxRetries {
		res, err := p.pool.Acquire(ctx)
		if err != nil {
			return nil, fmt.Errorf("acuire conn: %w", err)
		}

		conn := res.Value()
		if isConnAlive(conn) {
			return &puddleConn{
				Conn: conn,
				res:  res,
			}, nil
		}

		res.Destroy()
	}

	return nil, fmt.Errorf("failed to get valid connection after %d attempts", p.maxRetries)
}

func (p *PuddlePool) Close() {
	p.pool.Close()
}

var _ PoolConn = (*puddleConn)(nil)

type puddleConn struct {
	net.Conn
	res    *puddle.Resource[net.Conn]
	closed bool
}

func (p *puddleConn) Close() error {
	if p.closed {
		return nil
	}

	p.closed = true
	if p.res != nil {
		p.Conn.SetWriteDeadline(time.Time{})
		p.Conn.SetReadDeadline(time.Time{})
		p.res.Release()
		return nil
	}

	return p.Conn.Close()
}

func (p *puddleConn) Destroy() {
	if p.res != nil {
		p.res.Destroy()
	}
}

func isConnAlive(conn net.Conn) bool {
	conn.SetReadDeadline(time.Now().Add(1 * time.Millisecond))
	defer conn.SetReadDeadline(time.Time{}) // Reset deadline

	var buf [1]byte
	_, err := conn.Read(buf[:0])

	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			// Timeout is expected and means connection is likely alive
			return true
		}

		// Any other error means connection is dead
		return false
	}

	return true
}
