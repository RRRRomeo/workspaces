package map_go

import (
	"errors"
	"log"
	"math/rand"
	"sync"
)

// type SampleMap map[int]int
type SampleMap struct {
	sync.Mutex
	Dat map[int]string
}

// var wg sync.WaitGroup
var baseStrArr string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func SampleMap_Init(elementSize int) (*SampleMap, error) {
	log.Printf("%s ...\n", "SampleMap_Init")

	if elementSize < 0 {
		log.Printf("elementSize < 0\n")
		return nil, errors.New("elementSize < 0")
	}

	sampleMap := &SampleMap{
		Dat: make(map[int]string, elementSize),
	}
	// sampleMap := make(SampleMap, elementSize)

	return sampleMap, nil

}

func SampleMap_Write(m *SampleMap, element string, key int) (int, error) {
	log.Printf("%s ...\n", "SampleMap_Write")
	if m == nil {
		return -1, errors.New("SampleMap ptr is nil")
	}

	m.Lock()
	defer m.Unlock()

	if m.Dat[key] == element {
		return 0x8001, errors.New("element already exist and same")
	}

	if m.Dat[key] != element && m.Dat[key] != "" {
		return 0x8002, errors.New("element already exist difference")

	}

	m.Dat[key] = element
	return 0, nil
}

func SampleMap_Read(m *SampleMap, key int) (string, error) {
	log.Printf("%s...\n", "SampleMap_Read")
	if m == nil {
		log.Fatalf("inner map is nil\n")
		return "-1", errors.New("inner map is nil")
	}

	m.Lock()
	defer m.Unlock()
	val, ok := m.Dat[key]
	if ok {
		// m.Unlock()
		return val, nil
	}

	return "-1", errors.New("the locker dont unlock")
}

func Test_GoRunTime_Write(m *SampleMap, i int) {
	// log.Printf("%s...\n", "Test_GoRunTime_Write")
	// for i := 0; i < 50; i++ {
	// s := fmt.Sprintf("%d", i)
	SampleMap_Write(m, Test_RandomString(12), i)
	log.Printf("Sample_Write:%d\n", i)
	// }
	wg.Done()
}

func Test_GoRunTime_Read(m *SampleMap, i int) {
	// log.Printf("%s...\n", "Test_GoRunTime_Read")
	// for i := 0; i < 50; i++ {
	r, _ := SampleMap_Read(m, i)
	log.Printf("read :%d value:%s\n", i, r)
	// }
	wg.Done()
}

func Test_RandomString(lengh int) string {
	str := make([]byte, lengh)
	for i := 0; i < lengh; i++ {
		str[i] = baseStrArr[rand.Intn(len(baseStrArr))]
	}
	return string(str)
}
