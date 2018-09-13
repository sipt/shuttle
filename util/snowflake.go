package util

import (
	"errors"
	"sync"
	"time"
)

var IW *IdWorker

func init() {
	var err error
	IW, err = NewIdWorker(1)
	if err != nil {
		panic(err)
	}
}

func GetLongID() (i int64) {
	i, _ = IW.NextId()
	return
}

const (
	CEpoch         = 1474802888000
	CWorkerIdBits  = 10 // Num of WorkerId Bits
	CSenquenceBits = 12 // Num of Sequence Bits

	CWorkerIdShift  = 12
	CTimeStampShift = 22

	CSequenceMask = 0xfff // equal as getSequenceMask()
	CMaxWorker    = 0x3ff // equal as getMaxWorkerId()
)

// IdWorker Struct
type IdWorker struct {
	workerId      int64
	lastTimeStamp int64
	sequence      int64
	maxWorkerId   int64
	lock          *sync.Mutex
}

// NewIdWorker Func: Generate NewIdWorker with Given workerid
func NewIdWorker(workerid int64) (iw *IdWorker, err error) {
	iw = new(IdWorker)

	iw.maxWorkerId = getMaxWorkerId()

	if workerid > iw.maxWorkerId || workerid < 0 {
		return nil, errors.New("worker not fit")
	}
	iw.workerId = workerid
	iw.lastTimeStamp = -1
	iw.sequence = 0
	iw.lock = new(sync.Mutex)
	return iw, nil
}

func getMaxWorkerId() int64 {
	return -1 ^ -1<<CWorkerIdBits
}

func getSequenceMask() int64 {
	return -1 ^ -1<<CSenquenceBits
}

// return in ms
func (iw *IdWorker) timeGen() int64 {
	return time.Now().UnixNano() / 1000 / 1000
}

func (iw *IdWorker) timeReGen(last int64) int64 {
	ts := time.Now().UnixNano()
	for {
		if ts < last {
			ts = iw.timeGen()
		} else {
			break
		}
	}
	return ts
}

// NewId Func: Generate next id
func (iw *IdWorker) NextId() (ts int64, err error) {
	iw.lock.Lock()
	defer iw.lock.Unlock()
	ts = iw.timeGen()
	if ts == iw.lastTimeStamp {
		iw.sequence = (iw.sequence + 1) & CSequenceMask
		if iw.sequence == 0 {
			ts = iw.timeReGen(ts)
		}
	} else {
		iw.sequence = 0
	}

	if ts < iw.lastTimeStamp {
		err = errors.New("Clock moved backwards, Refuse gen id")
		return 0, err
	}
	iw.lastTimeStamp = ts
	ts = (ts-CEpoch)<<CTimeStampShift | iw.workerId<<CWorkerIdShift | iw.sequence
	return ts, nil
}

// ParseId Func: reverse uid to timestamp, workid, seq
func ParseId(id int64) (t time.Time, ts int64, workerId int64, seq int64) {
	seq = id & CSequenceMask
	workerId = (id >> CWorkerIdShift) & CMaxWorker
	ts = (id >> CTimeStampShift) + CEpoch
	t = time.Unix(ts/1000, (ts%1000)*1000000)
	return
}
