package main

import (
	"log"	
	"time"
	"sync"
)

type SyncMap struct {
	mx sync.Mutex
	m map[int]int
}

func NewSyncMap() *SyncMap {
	return &SyncMap{
		m: make(map[int]int),
	}
}

func (m *SyncMap) Store(key int, value int) {
	m.mx.Lock()
	defer m.mx.Unlock()
	m.m[key] = value
}

func (m *SyncMap) Load(key int) (int, bool) {
	val, ok := m.m[key]
	return val, ok
}

func main() {
	var wg sync.WaitGroup

	myMap := NewSyncMap()

	jobs := make(chan int)

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go workerW(&wg, myMap, jobs)
	}

	wg.Add(1)
	go producer(&wg, jobs)
	
	wg.Wait()

	// Чтобы не было гонки запись/чтение
	var wgR sync.WaitGroup	

	for i := 0; i < 2; i++ {
		wgR.Add(1)
		go workerR(&wgR, myMap)
	}

	wgR.Wait()
}

func workerW(wg *sync.WaitGroup, c *SyncMap, jobs <-chan int) {
	// Воркер, который пишит в map.
	defer wg.Done()
	
	for data := range jobs {
		c.Store(data, data)
	}
}

func workerR(wg *sync.WaitGroup, c *SyncMap) {
	// Воркер, который читает из map.
	defer wg.Done()

	// Проверка конкуретного чтения. Просто пример.
	time.Sleep(1*time.Second)

	for key, val := range c.m {
		log.Println("Чтение из map: ", key, val)
	}
}

func producer(wg *sync.WaitGroup, jobs chan<- int) {	
	defer wg.Done()	
	for i := 0; i < 5; i++ {
		jobs<-i	
	}
	close(jobs)
}
