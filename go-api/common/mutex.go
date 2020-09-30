package common

import (
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

const mutexLocked = 1 << iota

// Mutex 自定义实现TryLock锁
type Mutex struct {
	sync.Mutex
}

// TryLock 尝试获取锁
func (m *Mutex) TryLock() bool {
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), 0, mutexLocked)
}

// ChanMutex 利用chan构建的锁
type ChanMutex chan struct{}

// Lock 锁定
func (m *ChanMutex) Lock() {
	ch := (chan struct{})(*m)
	ch <- struct{}{}
}

// Unlock 解锁
func (m *ChanMutex) Unlock() {
	ch := (chan struct{})(*m)
	select {
	case <-ch:
	default:
		panic("unlock of unlocked mutex")
	}
}

// TryLock 尝试获取锁
func (m *ChanMutex) TryLock() bool {
	ch := (chan struct{})(*m)
	select {
	case ch <- struct{}{}:
		return true
	default:
	}
	return false
}

type spinLock uint32

// Lock 锁
func (sl *spinLock) Lock() {
	for !atomic.CompareAndSwapUint32((*uint32)(sl), 0, 1) {
		runtime.Gosched() //without this it locks up on GOMAXPROCS > 1
	}
}

// Unlock 解锁
func (sl *spinLock) Unlock() {
	atomic.StoreUint32((*uint32)(sl), 0)
}

// TryLock 尝试获取锁
func (sl *spinLock) TryLock() bool {
	return atomic.CompareAndSwapUint32((*uint32)(sl), 0, 1)
}

// SpinLock 自旋锁
func SpinLock() sync.Locker {
	var lock spinLock
	return &lock
}
