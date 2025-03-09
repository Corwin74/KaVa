package concurrency

// Semaphore -- ограничивает количество конкуретных задач в момент времени
type Semaphore struct {
	tickets chan struct{}
}

// NewSemaphore -- возвращает новый объект Semaphore
func NewSemaphore(n int) Semaphore {
	return Semaphore{tickets: make(chan struct{}, n)}
}

// Acquire -- попытка пройти за семафор
func (s *Semaphore) Acquire() {
	if s == nil || s.tickets == nil {
		return
	}
	s.tickets <- struct{}{}
}

// Release -- освобождение семафора для других задач
func (s *Semaphore) Release() {
	if s == nil || s.tickets == nil {
		return
	}
	
	<-s.tickets
}

// WithSemaphore -- helper
func (s *Semaphore) WithSemaphore(action func()) {
	if action == nil || s == nil {
		return
	}

	s.Acquire()
	action()
	s.Release()
}
