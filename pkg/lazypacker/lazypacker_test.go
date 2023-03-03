package lazypacker

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	c "github.com/smartystreets/goconvey/convey"
	"github.comn/SnowPhoenix0105/smartpacker/pkg/optional"
	"github.comn/SnowPhoenix0105/smartpacker/pkg/result"
)

type testPackParam[TV any] struct {
	DoPackFn func(ctx context.Context, param *testPackParam[TV]) result.Of[TV]
}

type testPackManager[TV any] struct{}

func (tpm testPackManager[TV]) DoPack(ctx context.Context, param *testPackParam[TV]) result.Of[TV] {
	return param.DoPackFn(ctx, param)
}

func TestLazyPacker_Pack(t *testing.T) {
	c.Convey("Test DoPack can be called only once", t, func() {
		const (
			TestTimes      = 16
			GoroutineCount = 128
		)

		for testTimes := range [TestTimes]struct{}{} {
			expect := testTimes*0xabcdef + 0xdeaddead
			c.Convey(fmt.Sprintf("Given expect=%d\t[%d]", expect, testTimes), func() {
				c.Convey(fmt.Sprintf("When %d gorountines called packer.Pack()", GoroutineCount), func() {
					doPackIsCalled := atomic.Int32{}
					ctx := context.TODO()
					param := &testPackParam[int]{
						DoPackFn: func(ctx context.Context, param *testPackParam[int]) result.Of[int] {
							doPackIsCalled.Add(1)
							return result.OfValue(expect)
						},
					}

					packer := NewLazyPacker[int, *testPackParam[int], testPackManager[int]]()
					resList := make([]result.Of[int], GoroutineCount)

					wg := sync.WaitGroup{}
					wg.Add(GoroutineCount)
					for goroutineID := range [GoroutineCount]struct{}{} {
						go func(gid int) {
							defer wg.Done()
							time.Sleep(time.Duration(rand.Float64() * float64(time.Millisecond) * 2))
							resList[gid] = packer.Pack(ctx, param)
						}(goroutineID)
					}
					wg.Wait()

					c.Convey(fmt.Sprintf("The result of each goroutine should equals to %d", expect), func() {
						for _, res := range resList {
							c.So(res.Error(), c.ShouldBeNil)
							c.So(res.Value(), c.ShouldEqual, expect)
						}
					})

					c.Convey("The PackManager::DoPack() should be called only once", func() {
						c.So(doPackIsCalled.Load(), c.ShouldEqual, 1)
					})
				})
			})
		}
	})
}

func TestLazyPacker_Set(t *testing.T) {
	c.Convey("SetParallel", t, func() {
		const (
			TestTimes      = 16
			GoroutineCount = 128
		)

		for testTimes := range [TestTimes]struct{}{} {
			c.Convey(fmt.Sprintf("Given nothing [%d]", testTimes), func() {
				c.Convey(fmt.Sprintf("When %d getter-goroutines called packer.TryGetResult() "+
					"and %d setter-goroutines called packer.Set(goroutineID)", GoroutineCount, GoroutineCount), func() {
					packer := NewLazyPacker[int, *testPackParam[int], testPackManager[int]]()
					setterResList := make([]bool, GoroutineCount)
					getterResList := make([]optional.Of[result.Of[int]], GoroutineCount)

					wg := sync.WaitGroup{}
					wg.Add(GoroutineCount * 2)
					for goroutineID := range [GoroutineCount]struct{}{} {
						go func(getterID int) {
							defer wg.Done()
							time.Sleep(time.Duration(rand.Float64() * float64(time.Millisecond) * 2))
							getterResList[getterID] = packer.TryGetResult()
						}(goroutineID)
						go func(setterID int) {
							defer wg.Done()
							time.Sleep(time.Duration(rand.Float64() * float64(time.Millisecond) * 2))
							setterResList[setterID] = packer.Set(setterID)
						}(goroutineID)
					}
					wg.Wait()

					c.Convey("The count of success setter should be 1", func() {
						okCnt := 0
						var okSetterID int
						for setterID, ok := range setterResList {
							if ok {
								okCnt++
								okSetterID = setterID
							}
						}
						_, _ = c.Print("and the No.", okSetterID, " setter success set the value")
						c.So(okCnt, c.ShouldEqual, 1)

						c.Convey(fmt.Sprintf("The results of getters should be nil or %d", okSetterID), func() {
							nilCnt := 0
							valCnt := 0
							for _, res := range getterResList {
								if !res.Ok() {
									nilCnt++
									continue
								}

								valCnt++
								c.So(res.Value().Value(), c.ShouldEqual, okSetterID)
							}
							c.Convey(fmt.Sprintf("The result is nil for %d getters and %d for %d getters",
								nilCnt, okSetterID, valCnt), func() {})
						})
					})
				})
			})
		}
	})
}
