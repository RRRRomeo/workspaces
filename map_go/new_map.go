package map_go

import (
	"log"
	"sync"
	"time"

	"golang.org/x/net/context"
)

var wg sync.WaitGroup

/**
* OOP
 */

type MyChan struct {
	sync.Once
	ch chan struct{}
}

type NewMap struct {
	sync.Mutex
	mp      map[int]string
	keyToch map[int]*MyChan
}

func NewMyChan() *MyChan {
	return &MyChan{
		ch: make(chan struct{}),
	}
}

func MallocNewMap(siz int) *NewMap {
	return &NewMap{
		mp:      make(map[int]string, siz),
		keyToch: make(map[int]*MyChan),
	}
}

func (c *MyChan) Close() {
	c.Do(func() {
		close(c.ch)
	})
}

func (m *NewMap) Put(k int, v string) {
	m.Lock()
	defer m.Unlock()
	m.mp[k] = v

	ch, ok := m.keyToch[k]
	if !ok {
		return
	}

	/**
	 * ch <- struct{}{} ====>为什么要close:!
	 * 多次close 会造成panic
	 * 	select {
	 *	case <-ch:
	 *		return
	 *	default:
	 *		close(ch)
	 *	}
	 */

	ch.Close()

}

func (m *NewMap) Get(k int, maxWaittingDuration time.Duration) (string, error) {
	m.Lock()
	v, ok := m.mp[k]
	if ok {
		m.Unlock()
		return v, nil
	}

	// ...
	ch, ok := m.keyToch[k]
	if !ok {
		// ch = make(chan struct{})
		ch = NewMyChan()
		m.keyToch[k] = ch
	}

	tCtx, cancle := context.WithTimeout(context.Background(), maxWaittingDuration)
	defer cancle()

	m.Unlock()
	select {
	case <-tCtx.Done():
		return "-1", tCtx.Err()
	case <-ch.ch:

	}

	m.Lock()
	v = m.mp[k]
	m.Unlock()

	return v, nil

}
func Test_Write(m *NewMap, i int) {
	defer wg.Done()
	log.Printf("%s...\n", "Test_GoRunTime_Write")
	// for i := 0; i < 50; i++ {
	// s := fmt.Sprintf("%d", i)
	m.Put(i, Test_RandomString(12))
	log.Printf("Test_Write:%d value:%s\n", i, m.mp[i])
	// }
	wg.Done()
}

func Test_Read(m *NewMap, i int) {
	defer wg.Done()
	log.Printf("%s...\n", "Test_GoRunTime_Read")
	// for i := 0; i < 50; i++ {
	r, _ := m.Get(i, 100*time.Millisecond)
	log.Printf("Test_Read :%d value:%s\n", i, r)
	// }
	wg.Done()
}

func GetWg() *sync.WaitGroup {
	return &wg
}
