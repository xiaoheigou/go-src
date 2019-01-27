package timewheel

import (
	"sync"
	"time"
	"yuudidi.com/pkg/utils"
)

type slot struct {
	id       int
	elements map[interface{}]interface{}
}

func newSlot(id int) *slot {
	s := &slot{id: id}
	s.elements = make(map[interface{}]interface{})
	return s
}

func (s *slot) add(c interface{}) {
	s.elements[c] = c
}

func (s *slot) remove(c interface{}) {
	delete(s.elements, c)
}

type handler func(interface{})

type TimeWheel struct {
	tickDuration     time.Duration
	ticksPerWheel    int
	currentTickIndex int
	ticker           *time.Ticker
	redisKey         string
	onTick           handler
	wheel            []*slot
	indicator        map[interface{}]*slot
	sync.RWMutex

	taskChan chan interface{}
	quitChan chan interface{}
}

func New(tickDuration time.Duration, ticksPerWheel int, redisKey string, f handler) *TimeWheel {
	if tickDuration < 1 || ticksPerWheel < 1 || nil == f {
		return nil
	}

	ticksPerWheel++
	t := &TimeWheel{
		tickDuration:     tickDuration,
		ticksPerWheel:    ticksPerWheel,
		redisKey:         redisKey,
		onTick:           f,
		currentTickIndex: 0,
		taskChan:         make(chan interface{}),
		quitChan:         make(chan interface{}),
	}
	t.indicator = make(map[interface{}]*slot, 0)

	t.wheel = make([]*slot, ticksPerWheel)
	for i := 0; i < ticksPerWheel; i++ {
		t.wheel[i] = newSlot(i)
	}

	return t
}

func (t *TimeWheel) Start() {
	t.ticker = time.NewTicker(t.tickDuration)
	go t.run()
	if values, err := utils.GetCacheSetMembers(t.redisKey); err != nil {
		for _, v := range values {
			t.taskChan <- v
		}
	}
}

func (t *TimeWheel) Add(c interface{}) {
	t.taskChan <- c
	utils.SetCacheSetMember(t.redisKey, 0, c)
}

func (t *TimeWheel) Remove(c interface{}) {
	if v, ok := t.indicator[c]; ok {
		v.remove(c)
		utils.DelCacheSetMember(t.redisKey, c)
	}
}

func (t *TimeWheel) getPreviousTickIndex() int {
	t.RLock()
	defer t.RUnlock()

	cti := t.currentTickIndex
	if 0 == cti {
		return t.ticksPerWheel - 1
	}
	return cti - 1
}

func (t *TimeWheel) Stop() {
	close(t.quitChan)
	utils.RedisClient.Del(t.redisKey)
}

func (t *TimeWheel) run() {
	for {
		select {
		case <-t.quitChan:
			t.ticker.Stop()
			break
		case <-t.ticker.C:
			if t.ticksPerWheel == t.currentTickIndex {
				t.currentTickIndex = 0
			}

			slot := t.wheel[t.currentTickIndex]
			for _, v := range slot.elements {
				slot.remove(v)
				delete(t.indicator, v)
				t.onTick(v)
				utils.DelCacheSetMember(t.redisKey, v)
			}

			t.currentTickIndex++
		case v := <-t.taskChan:
			t.Remove(v)
			slot := t.wheel[t.getPreviousTickIndex()]
			slot.add(v)
			t.indicator[v] = slot
		}
	}
}
