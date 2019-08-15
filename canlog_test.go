package canlog_test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/mikesimons/canlog"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCanlog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Canlog")
}

var _ = Describe("Canlog", func() {
	Describe("Ref", func() {
		It("should return a unique id", func() {
			for i := 0; i < 200; i++ {
				r1 := canlog.Ref()
				r2 := canlog.Ref()

				Expect(r1).ShouldNot(BeEquivalentTo(r2))
			}
		})
	})

	Describe("Push", func() {
		It("should push data which can be popped", func() {
			ref := canlog.Ref()
			canlog.Push(ref, "test", "data")
			data, _ := canlog.Pop(ref)
			Expect(data["test"]).To(Equal("data"))
		})

		It("should append data to context", func() {
			ref := canlog.Ref()
			canlog.Push(ref, "test", "data")
			canlog.Push(ref, "test1", "data1")
			canlog.Push(ref, "test2", "data2")

			data, _ := canlog.Pop(ref)
			Expect(data["test"]).To(Equal("data"))
			Expect(data["test1"]).To(Equal("data1"))
			Expect(data["test2"]).To(Equal("data2"))
		})

		It("should return an error if ref is popped", func() {
			ref := canlog.Ref()
			canlog.Pop(ref)
			err := canlog.Push(ref, "test", "data")
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if ref is invalid", func() {
			err := canlog.Push(canlog.Reference(""), "", "")
			Expect(err).NotTo(BeNil())
		})

		It("should be goroutine safe", func() {
			ref := canlog.Ref()
			wg := &sync.WaitGroup{}

			for i := 0; i < 20; i++ {
				wg.Add(1)
				go func(i int) {
					canlog.Push(ref, fmt.Sprintf("test%d", i), "test")
					time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
					wg.Done()
				}(i)
			}

			wg.Wait()
			data, _ := canlog.Pop(ref)

			for i := 0; i < 20; i++ {
				key := fmt.Sprintf("test%d", i)
				Expect(data[key]).To(Equal("test"), key)
			}
		})
	})

	Describe("Pop", func() {
		It("should return pushed data", func() {
			ref := canlog.Ref()
			canlog.Push(ref, "string", "stringval")
			canlog.Push(ref, "int", 1)
			canlog.Push(ref, "map", map[string]string{"k1": "v1"})
			data, _ := canlog.Pop(ref)
			Expect(data["string"]).To(Equal("stringval"))
			Expect(data["int"]).To(Equal(1))
			Expect(data["map"].(map[string]string)["k1"]).To(Equal("v1"))
		})

		It("should return error for invalid ref", func() {
			_, err := canlog.Pop(canlog.Reference(""))
			Expect(err).NotTo(BeNil())
		})

		It("should remove ref once popped", func() {
			ref := canlog.Ref()
			_, err1 := canlog.Pop(ref)
			_, err2 := canlog.Pop(ref)

			Expect(err1).To(BeNil())
			Expect(err2).NotTo(BeNil())
		})
	})
})
