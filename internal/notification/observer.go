package notification

import "sync"

type NotificationObserver interface {
	Update(notification *Notification)
}

type NotificationSubject struct {
	observers []NotificationObserver
	mu        sync.RWMutex
}

func NewNotificationSubject() *NotificationSubject {
	return &NotificationSubject{
		observers: []NotificationObserver{},
	}
}

func (s *NotificationSubject) Attach(observer NotificationObserver) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.observers = append(s.observers, observer)
}

func (s *NotificationSubject) Detach(observer NotificationObserver) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, obs := range s.observers {
		if obs == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

func (s *NotificationSubject) Notify(notification *Notification) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, observer := range s.observers {
		go observer.Update(notification)
	}
}

type WebSocketObserver struct {
	connections map[int64]chan *Notification
	mu          sync.RWMutex
}

func NewWebSocketObserver() *WebSocketObserver {
	return &WebSocketObserver{
		connections: make(map[int64]chan *Notification),
	}
}

func (w *WebSocketObserver) AddConnection(userID int64, ch chan *Notification) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.connections[userID] = ch
}

func (w *WebSocketObserver) RemoveConnection(userID int64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if ch, exists := w.connections[userID]; exists {
		close(ch)
		delete(w.connections, userID)
	}
}

func (w *WebSocketObserver) Update(notification *Notification) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if ch, exists := w.connections[notification.UserID]; exists {
		select {
		case ch <- notification:
		default:
		}
	}
}

type LogObserver struct{}

func NewLogObserver() *LogObserver {
	return &LogObserver{}
}

func (l *LogObserver) Update(notification *Notification) {
	println("[NOTIFICATION]", notification.Type, "to user", notification.UserID, ":", notification.Message)
}
