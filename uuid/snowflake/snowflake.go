package snowflake

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strconv"
	"sync"
	"time"
)

var (
	instance *Node
	once     = sync.Once{}
)

func Generate() ID {
	once.Do(func() {
		var err error
		if instance, err = NewNode(0); err != nil {
			panic(err)
		}
	})
	return instance.Generate()
}

var (
	Epoch    int64 = 1288834974657
	NodeBits uint8 = 10
	StepBits uint8 = 12

	nodeMax   int64 = -1 ^ (-1 << NodeBits)
	nodeMask  int64 = nodeMax << StepBits
	stepMask  int64 = -1 ^ (-1 << StepBits)
	timeShift uint8 = NodeBits + StepBits
	nodeShift uint8 = StepBits
)

type Node struct {
	mu   sync.Mutex
	time int64
	node int64
	step int64
}

type ID int64

func NewNode(node int64) (*Node, error) {
	nodeMax = -1 ^ (-1 << NodeBits)
	nodeMask = nodeMax << StepBits
	stepMask = -1 ^ (-1 << StepBits)
	timeShift = NodeBits + StepBits
	nodeShift = StepBits

	if node < 0 || node > nodeMax {
		return nil, errors.New("node number must be between 0 and " + strconv.FormatInt(nodeMax, 10))
	}

	return &Node{
		time: 0,
		node: node,
		step: 0,
	}, nil
}

func (n *Node) Generate() ID {
	n.mu.Lock()
	now := time.Now().UnixNano() / 1000000
	if n.time == now {
		n.step = (n.step + 1) & stepMask

		if n.step == 0 {
			for now <= n.time {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		n.step = 0
	}
	n.time = now
	r := ID((now-Epoch)<<timeShift |
		(n.node << nodeShift) |
		(n.step),
	)
	n.mu.Unlock()
	return r
}

func (f ID) Int64() int64 {
	return int64(f)
}

func (f ID) String() string {
	return strconv.FormatInt(int64(f), 10)
}

func (f ID) Base2() string {
	return strconv.FormatInt(int64(f), 2)
}

func (f ID) Base36() string {
	return strconv.FormatInt(int64(f), 36)
}

func (f ID) Base64() string {
	return base64.StdEncoding.EncodeToString(f.Bytes())
}

func (f ID) MD5() string {
	m := md5.Sum(f.Bytes())
	return hex.EncodeToString(m[:])
}

func (f ID) Bytes() []byte {
	return []byte(f.String())
}

func (f ID) Time() int64 {
	return (int64(f) >> timeShift) + Epoch
}

func (f ID) Node() int64 {
	return int64(f) & nodeMask >> nodeShift
}

func (f ID) Step() int64 {
	return int64(f) & stepMask
}
