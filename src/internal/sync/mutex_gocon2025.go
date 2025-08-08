package sync

import "sync/atomic"

// TODO: Define constants separately for each version?
const (
	myMutexLocked                = 1 << iota // true if mutex is locked
	myMutexWoken                             // true if there is at least one awoken G
	myMutexStarving                          // true if mutex is in starving mode
	myMutexWaiterShift           = iota      // number of G's waiting on mutex
	myMutexStarvationThresholdNs = 1e6       // threshold time to enter starvation mode
)

/******************************************************************************/
/*                                  MyMutex1                                  */
/******************************************************************************/

type MyMutex1 struct {
	sema uint32
}

func NewMyMutex1() MyMutex1 {
	return MyMutex1{sema: 1} // need to init sema to 1
}

// NOTE: No TryLock() possible

func (m *MyMutex1) Lock() {
	println("Locking MyMutex1...") // using builtin println() to prevent cyclic deps
	defer println("Locking MyMutex1 complete!")

	queueLifo := false
	skipframes := 1 // skip 1 caller from stack trace (sync.(*MyMutex).Lock())
	runtime_SemacquireMutex(&m.sema, queueLifo, skipframes)
}

func (m *MyMutex1) Unlock() {
	println("Unlocking MyMutex1...")
	defer println("Unlocking MyMutex1 complete!")

	handoff := false
	skipframes := 1
	runtime_Semrelease(&m.sema, handoff, skipframes)
}

/******************************************************************************/
/*                                  MyMutex2                                  */
/******************************************************************************/

type MyMutex2 struct {
	state int32 // could use uint23 instead
	sema  uint32
}

func (m *MyMutex2) TryLock() bool {
	println("Trying to lock MyMutex2...")
	defer println("Trying to lock MyMutex2 complete!")

	if atomic.SwapInt32(&m.state, myMutexLocked) != 0 {
		return false
	}
	return true
}

func (m *MyMutex2) Lock() {
	println("Locking MyMutex2...")
	defer println("Locking MyMutex2 complete!")

	for atomic.SwapInt32(&m.state, myMutexLocked) != 0 {
		queueLifo := false
		skipframes := 1
		runtime_SemacquireMutex(&m.sema, queueLifo, skipframes)
	}
}

func (m *MyMutex2) Unlock() {
	println("Unlocking MyMutex2...")
	defer println("Unlocking MyMutex2 complete!")

	atomic.StoreInt32(&m.state, 0)
	handoff := false
	skipframes := 1
	// not part of runtime package
	runtime_SemreleaseNoDup(&m.sema, handoff, skipframes)
}

/******************************************************************************/
/*                                  MyMutex3                                  */
/******************************************************************************/

type MyMutex3 struct {
	state int32
	sema  uint32
}

func (m *MyMutex3) TryLock() bool {
	println("Trying to lock MyMutex3...")
	defer println("Trying to lock MyMutex3 complete!")

	if atomic.SwapInt32(&m.state, myMutexLocked) != 0 {
		return false
	}
	return true
}

func (m *MyMutex3) Lock() {
	println("Locking MyMutex3...")
	defer println("Locking MyMutex3 complete!")

	iter := 0
	for atomic.SwapInt32(&m.state, myMutexLocked) != 0 {
		if runtime_canSpin(iter) {
			runtime_doSpin() // spin by yielding CPU if possible
			iter++
			continue
		}
		queueLifo := false
		skipframes := 1
		runtime_SemacquireMutex(&m.sema, queueLifo, skipframes)
	}
}

func (m *MyMutex3) Unlock() {
	println("Unlocking MyMutex3...")
	defer println("Unlocking MyMutex3 complete!")

	atomic.StoreInt32(&m.state, 0)
	handoff := false
	skipframes := 1
	runtime_SemreleaseNoDup(&m.sema, handoff, skipframes)
}

/******************************************************************************/
/*                                  MyMutex4                                  */
/******************************************************************************/

type MyMutex4 struct {
	state int32
	sema  uint32
}

func (m *MyMutex4) TryLock() bool {
	println("Trying to lock MyMutex4...")
	defer println("Trying to lock MyMutex4 complete!")

	old := m.state
	if old&myMutexLocked != 0 {
		return false
	}
	if !atomic.CompareAndSwapInt32(&m.state, old, old|myMutexLocked) {
		// allows current G to barge in before waiting G's
		// it could be also be that mutex is still unlocked but waiter count decremented by Unlock()
		return false
	}
	return true
}

func (m *MyMutex4) Lock() {
	println("Locking MyMutex4...")
	defer println("Locking MyMutex4 complete!")

	if atomic.CompareAndSwapInt32(&m.state, 0, myMutexLocked) {
		return // allows current G to barge in before waiting G's
	}
	m.lockSlow()
}

func (m *MyMutex4) lockSlow() {
	iter := 0
	old := m.state // not atomic but okay due to memory barriers
	for {
		if old&myMutexLocked != 0 && runtime_canSpin(iter) {
			runtime_doSpin()
			iter++
			old = m.state
			continue
		}
		new := old | myMutexLocked
		if old&myMutexLocked != 0 {
			new += 1 << myMutexWaiterShift
		}
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			if old&myMutexLocked == 0 {
				break // acquired mutex successfully with CAS
			}
			queueLifo := false
			skipframes := 2 // skip 2 callers from stack trace (isync.(*MyMutex4sync).Lock() and sync.(*MyMutex).Lock())
			runtime_SemacquireMutex(&m.sema, queueLifo, skipframes)
			iter = 0
		}
		old = m.state
	}
}

func (m *MyMutex4) Unlock() {
	println("Unlocking MyMutex4...")
	defer println("Unlocking MyMutex4 complete!")

	new := atomic.AddInt32(&m.state, -myMutexLocked)
	if new == 0 {
		return
	}
	m.unlockSlow(new)
}

func (m *MyMutex4) unlockSlow(new int32) {
	if (new+myMutexLocked)&myMutexLocked == 0 { // add back myMutexLocked in case it was not set initially
		fatal("Unlocked unlocked MyMutex4!")
	}

	old := new
	for {
		if old>>myMutexWaiterShift == 0 || // no need to wake up if there are no waiting G's
			old&myMutexLocked != 0 { // no need to wake up if some G barged in and acquired mutex already
			return
		}
		new = old - 1<<myMutexWaiterShift
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			handoff := false
			skipframes := 2
			runtime_Semrelease(&m.sema, handoff, skipframes)
		}
		old = m.state
	}
}

/******************************************************************************/
/*                                  MyMutex5                                  */
/******************************************************************************/

type MyMutex5 struct {
	state int32
	sema  uint32
}

func (m *MyMutex5) TryLock() bool {
	println("Trying to lock MyMutex5...")
	defer println("Trying to lock MyMutex5 complete!")

	old := m.state
	if old&myMutexLocked != 0 {
		return false
	}
	if !atomic.CompareAndSwapInt32(&m.state, old, old|myMutexLocked) {
		return false
	}
	return true
}

func (m *MyMutex5) Lock() {
	println("Locking MyMutex5...")
	defer println("Locking MyMutex5 complete!")

	if atomic.CompareAndSwapInt32(&m.state, 0, myMutexLocked) {
		return
	}
	m.lockSlow()
}

func (m *MyMutex5) lockSlow() {
	awoke := false // true if current G is spinning awake
	iter := 0
	old := m.state
	for {
		if old&myMutexLocked != 0 && runtime_canSpin(iter) {
			if !awoke && // awoke is set to true if myMutexWoken is successfully set by current G or waking up from sema acquire
				// no need to set myMutexWoken again if already set successfully current G
				// no need to set myMutexWoken when waking up from sema acquire because Unlock() handles it
				old&myMutexWoken == 0 && // no need to set myMutexWoken again if it is already set
				// this keeps awoke flag to false for G's barging in to acquire mutex
				// while another G is about to sema release in Unlock() so that sema value never becomes more than 1
				old>>myMutexWaiterShift != 0 && // if no waiters then Unlock() will not attempt to wake up anyways
				atomic.CompareAndSwapInt32(&m.state, old, old|myMutexWoken) {
				awoke = true
			}
			runtime_doSpin()
			iter++
			old = m.state
			continue
		}
		new := old | myMutexLocked
		if old&myMutexLocked != 0 {
			new += 1 << myMutexWaiterShift
		}
		if awoke {
			new &^= myMutexWoken // clear myMutexWoken bit if successfully acquired mutex or going to sleep
		}
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			if old&myMutexLocked == 0 {
				break
			}
			queueLifo := false
			skipframes := 2
			runtime_SemacquireMutex(&m.sema, queueLifo, skipframes)
			awoke = true
			iter = 0
		}
		old = m.state
	}
}

func (m *MyMutex5) Unlock() {
	println("Unlocking MyMutex5...")
	defer println("Unlocking MyMutex5 complete!")

	new := atomic.AddInt32(&m.state, -myMutexLocked)
	if new == 0 {
		return
	}
	m.unlockSlow(new)
}

func (m *MyMutex5) unlockSlow(new int32) {
	if (new+myMutexLocked)&myMutexLocked == 0 {
		fatal("Unlocked unlocked MyMutex5!")
	}

	old := new
	for {
		if old>>myMutexWaiterShift == 0 ||
			old&(myMutexLocked|myMutexWoken) != 0 { // no need to wake up if there is some G spinning awake
			return
		}
		new = (old - 1<<myMutexWaiterShift) | myMutexWoken // set myMutexWoken bit if successfully woke up some waiting G
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			handoff := false
			skipframes := 2
			runtime_Semrelease(&m.sema, handoff, skipframes)
		}
		old = m.state
	}
}

/******************************************************************************/
/*                                  MyMutex6                                  */
/******************************************************************************/

type MyMutex6 struct {
	state int32
	sema  uint32
}

func (m *MyMutex6) TryLock() bool {
	println("Trying to lock MyMutex6...")
	defer println("Trying to lock MyMutex6 complete!")

	old := m.state
	if old&myMutexLocked != 0 {
		return false
	}
	if !atomic.CompareAndSwapInt32(&m.state, old, old|myMutexLocked) {
		return false
	}
	return true
}

func (m *MyMutex6) Lock() {
	println("Locking MyMutex6...")
	defer println("Locking MyMutex6 complete!")

	if atomic.CompareAndSwapInt32(&m.state, 0, myMutexLocked) {
		return
	}
	m.lockSlow()
}

func (m *MyMutex6) lockSlow() {
	waited := false // true if waited at least once already
	awoke := false
	iter := 0
	old := m.state
	for {
		if old&myMutexLocked != 0 && runtime_canSpin(iter) {
			if !awoke &&
				old&myMutexWoken == 0 &&
				old>>myMutexWaiterShift != 0 &&
				atomic.CompareAndSwapInt32(&m.state, old, old|myMutexWoken) {
				awoke = true
			}
			runtime_doSpin()
			iter++
			old = m.state
			continue
		}
		new := old | myMutexLocked
		if old&myMutexLocked != 0 {
			new += 1 << myMutexWaiterShift
		}
		if awoke {
			new &^= myMutexWoken
		}
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			if old&myMutexLocked == 0 {
				break
			}
			queueLifo := waited // insert at the front of waiter queue if waiting for more than once
			skipframes := 2
			waited = true
			runtime_SemacquireMutex(&m.sema, queueLifo, skipframes)
			awoke = true
			iter = 0
		}
		old = m.state
	}
}

func (m *MyMutex6) Unlock() {
	println("Unlocking MyMutex6...")
	defer println("Unlocking MyMutex6 complete!")

	new := atomic.AddInt32(&m.state, -myMutexLocked)
	if new == 0 {
		return
	}
	m.unlockSlow(new)
}

func (m *MyMutex6) unlockSlow(new int32) {
	if (new+myMutexLocked)&myMutexLocked == 0 {
		fatal("Unlocked unlocked MyMutex6!")
	}

	old := new
	for {
		if old>>myMutexWaiterShift == 0 ||
			old&(myMutexLocked|myMutexWoken) != 0 {
			return
		}
		new = (old - 1<<myMutexWaiterShift) | myMutexWoken
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			handoff := false
			skipframes := 2
			runtime_Semrelease(&m.sema, handoff, skipframes)
		}
		old = m.state
	}
}

/******************************************************************************/
/*                                  MyMutex7                                  */
/******************************************************************************/

type MyMutex7 struct {
	state int32
	sema  uint32
}

func (m *MyMutex7) TryLock() bool {
	println("Trying to lock MyMutex7...")
	defer println("Trying to lock MyMutex7 complete!")

	old := m.state
	if old&(myMutexLocked|myMutexStarving) != 0 {
		return false // do not allow current G to barging in when starvation mode is on
	}
	if !atomic.CompareAndSwapInt32(&m.state, old, old|myMutexLocked) {
		return false
	}
	return true
}

func (m *MyMutex7) Lock() {
	println("Locking MyMutex7...")
	defer println("Locking MyMutex7 complete!")

	if atomic.CompareAndSwapInt32(&m.state, 0, myMutexLocked) {
		return
	}
	m.lockSlow()
}

func (m *MyMutex7) lockSlow() {
	var waitStartTime int64 // merged the waited flag into waitStartTime
	starving := false
	awoke := false
	iter := 0
	old := m.state
	for {
		if old&(myMutexLocked|myMutexStarving) == myMutexLocked && // only spin if not in starvation mode
			runtime_canSpin(iter) {
			if !awoke &&
				old&myMutexWoken == 0 &&
				old>>myMutexWaiterShift != 0 &&
				atomic.CompareAndSwapInt32(&m.state, old, old|myMutexWoken) {
				awoke = true
			}
			runtime_doSpin()
			iter++
			old = m.state
			continue
		}
		new := old
		if old&myMutexStarving == 0 {
			new |= myMutexLocked // newly arriving G should not barge in during starvation mode
		}
		if old&(myMutexLocked|myMutexStarving) != 0 {
			new += 1 << myMutexWaiterShift // newly arriving G must always sleep in starvation mode
		}
		if starving && old&mutexLocked != 0 {
			// enter starvation mode if current G is starving
			// but no need if mutex is already unlocked
			new |= myMutexStarving
		}
		if awoke {
			new &^= myMutexWoken
		}
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			if old&(myMutexLocked|myMutexStarving) == 0 {
				break // newly arriving G should not barge in during starvation mode
			}
			queueLifo := waitStartTime != 0
			if waitStartTime == 0 {
				waitStartTime = runtime_nanotime() // start timer since first sleep
			}
			skipframes := 2
			runtime_SemacquireMutex(&m.sema, queueLifo, skipframes)
			// flag starvation mode for next CAS after threshold is reached for current G
			starving = starving || runtime_nanotime()-waitStartTime > myMutexStarvationThresholdNs
			old = m.state // get latest state after wake up
			if old&myMutexStarving != 0 {
				delta := int32(myMutexLocked - 1<<myMutexWaiterShift)
				if !starving || // if current G is not starving then other waiters are not starving either because of LIFO order
					old>>myMutexWaiterShift == 1 { // if current G is last waiting G then clearly no waiters are starving
					delta -= myMutexStarving
				}
				// starvation mode guarantees no other G's will barge in
				// so must be safe to set myMutexLocked bit and decrement waiter count without CAS
				atomic.AddInt32(&m.state, delta)
				break // successfully acquired mutex via hand off
			}
			awoke = true
			iter = 0
		} else {
			old = m.state
		}
	}
}

func (m *MyMutex7) Unlock() {
	println("Unlocking MyMutex7...")
	defer println("Unlocking MyMutex7 complete!")

	// myMutexLocked bit is dropped during handoff in starvation mode
	// this is okay because Lock() and TryLock() checks myMutexStarving before barging i
	new := atomic.AddInt32(&m.state, -myMutexLocked)
	if new == 0 {
		return
	}
	m.unlockSlow(new)
}

func (m *MyMutex7) unlockSlow(new int32) {
	if (new+myMutexLocked)&myMutexLocked == 0 {
		fatal("Unlocked unlocked MyMutex7!")
	}

	if new&myMutexStarving == 0 {
		old := new
		for {
			if old>>myMutexWaiterShift == 0 ||
				old&(myMutexLocked|myMutexWoken|myMutexStarving) != 0 { // no need to wake up if some G put mutex into starvation mode
				return
			}
			new = (old - 1<<myMutexWaiterShift) | myMutexWoken
			if atomic.CompareAndSwapInt32(&m.state, old, new) {
				handoff := false
				skipframes := 2
				runtime_Semrelease(&m.sema, handoff, skipframes)
			}
			old = m.state
		}
	} else {
		handoff := true // directly hand off mutex to starving G at the front of waiter queue
		skipframes := 2
		// setting handoff to true in runtime semaphore makes releasing G to yield CPU immediately
		// so that starving G's can be rescheduled
		runtime_Semrelease(&m.sema, handoff, skipframes)
	}
}
