package lazypacker

import (
	"context"
	"sync"
	"sync/atomic"

	"github.comn/SnowPhoenix0105/smartpacker/pkg/optional"
	"github.comn/SnowPhoenix0105/smartpacker/pkg/result"
	"github.comn/SnowPhoenix0105/smartpacker/pkg/utils"
)

type PackManager[TV, TP any] interface {
	DoPack(ctx context.Context, param TP) result.Of[TV]
}

type LazyPacker[TV, TP any, TM PackManager[TV, TP]] struct {
	mutex sync.Mutex
	val   atomic.Pointer[result.Of[TV]]
}

func NewLazyPacker[TV, TP any, TM PackManager[TV, TP]]() *LazyPacker[TV, TP, TM] {
	res := MakeLazyPacker[TV, TP, TM]()
	return &res
}

func MakeLazyPacker[TV, TP any, TM PackManager[TV, TP]]() LazyPacker[TV, TP, TM] {
	return LazyPacker[TV, TP, TM]{}
}

func (lp *LazyPacker[TV, TP, TM]) Pack(ctx context.Context, param TP) result.Of[TV] {
	if res := lp.loadResult(); res != nil {
		return *res
	}
	return lp.packFromMgr(ctx, param)
}

func (lp *LazyPacker[TV, TP, TM]) TryGetResult() optional.Of[result.Of[TV]] {
	return optional.FromPtr(lp.loadResult())
}

func (lp *LazyPacker[TV, TP, TM]) Set(val TV) bool {
	if lp.loadResult() != nil {
		return false
	}
	return lp.setResult(result.PtrOfValue(val))
}

func (lp *LazyPacker[TV, TP, TM]) Touch(ctx context.Context, param TP) {
	if lp.loadResult() != nil {
		return
	}

	go func() {
		if lp.loadResult() != nil {
			return
		}
		lp.packFromMgr(ctx, param)
	}()
}

func (lp *LazyPacker[TV, TP, TM]) loadResult() *result.Of[TV] {
	return lp.val.Load()
}

func (lp *LazyPacker[TV, TP, TM]) packFromMgr(ctx context.Context, param TP) (finalRes result.Of[TV]) {
	lp.mutex.Lock()
	defer lp.mutex.Unlock()

	if ptr := lp.loadResult(); ptr != nil {
		return *ptr
	}

	defer lp.val.Store(&finalRes)
	defer utils.PanicRecoverWithFinalErrorPtr(&finalRes.Err)
	return utils.ZeroOf[TM]().DoPack(ctx, param)
}

func (lp *LazyPacker[TV, TP, TM]) setResult(res *result.Of[TV]) bool {
	lp.mutex.Lock()
	defer lp.mutex.Unlock()
	if lp.loadResult() != nil {
		return false
	}

	lp.val.Store(res)
	return true
}
