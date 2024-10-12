package semaphore

type Semaphore struct {
	sem chan bool
}

func New(limit int) *Semaphore {
	return &Semaphore{
		sem: make(chan bool, limit),
	}
}

func (s *Semaphore) Acquire() bool {
	select {
	case s.sem <- true:
		return true
	default:
		return false
	}
}

func (s *Semaphore) Release() {
	<-s.sem
}
