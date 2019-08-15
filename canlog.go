package canlog

import (
	"fmt"
	"sync"

	uuid "github.com/satori/go.uuid"
)

type Reference string

var refMap sync.Map

type entry struct {
	data map[string]interface{}
	ch   chan (pair)
	wg   *sync.WaitGroup
}

type pair struct {
	key   string
	value interface{}
	stop  bool
}

func initRef(ref Reference) {
	e := &entry{
		data: make(map[string]interface{}),
		ch:   make(chan (pair), 50),
		wg:   &sync.WaitGroup{},
	}

	refMap.Store(ref, *e)

	e.wg.Add(1)
	go func() {
		for data := range e.ch {
			if data.stop {
				e.wg.Done()
				return
			}

			e.data[data.key] = data.value
		}
	}()
}

func Ref() Reference {
	ref := Reference(uuid.NewV4().String())
	initRef(ref)
	return ref
}

func Push(ref interface{}, key string, value interface{}) error {
	rme, ok := refMap.Load(ref)

	if !ok {
		return fmt.Errorf("canlog reference '%s' is already popped", ref)
	}

	refMapEntry := rme.(entry)
	refMapEntry.ch <- pair{key: key, value: value}

	return nil
}

func Pop(ref Reference) (map[string]interface{}, error) {
	if rme, ok := refMap.Load(ref); ok {
		refMapEntry := rme.(entry)
		refMapEntry.ch <- pair{stop: true}
		refMapEntry.wg.Wait()

		refMap.Delete(ref)

		return refMapEntry.data, nil
	}
	return nil, fmt.Errorf("invalid canlog reference: %s", ref)
}
