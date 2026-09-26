package main

import (
	"sync"
	"time"
)

type SessionManager struct {
	mu       sync.RWMutex
	sessions map[int64]*UserSession
	ttl      time.Duration
}

func NewSessionManager(ttl time.Duration) *SessionManager {
	sm := &SessionManager{
		sessions: make(map[int64]*UserSession),
		ttl:      ttl,
	}
	go sm.startCleaner(10 * time.Minute)
	return sm
}

func (sm *SessionManager) GetOrCreate(userID int64) UserSession {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, exists := sm.sessions[userID]
	if !exists {
		session = &UserSession{
			UserID:      userID,
			CurrentStep: StepBudget,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		sm.sessions[userID] = session
	}
	return *session
}

func (sm *SessionManager) Mutate(userID int64, fn func(s *UserSession)) UserSession {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, exists := sm.sessions[userID]
	if !exists {
		session = &UserSession{
			UserID:      userID,
			CurrentStep: StepBudget,
			CreatedAt:   time.Now(),
		}
		sm.sessions[userID] = session
	}

	fn(session)
	session.UpdatedAt = time.Now()
	return *session
}

func (sm *SessionManager) Reset(userID int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, userID)
}

func (sm *SessionManager) startCleaner(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		sm.mu.Lock()
		now := time.Now()
		for id, session := range sm.sessions {
			if now.Sub(session.UpdatedAt) > sm.ttl {
				delete(sm.sessions, id)
			}
		}
		sm.mu.Unlock()
	}
}
