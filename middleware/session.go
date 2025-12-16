package middleware

import (
	"github.com/gin-gonic/gin"
	"sync"
)

// Simple in-memory session store
var sessionStore = &SessionStore{
	sessions: make(map[string]map[string]interface{}),
}

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]map[string]interface{}
}

func (s *SessionStore) Get(sessionID, key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if session, ok := s.sessions[sessionID]; ok {
		val, exists := session[key]
		return val, exists
	}
	return nil, false
}

func (s *SessionStore) Set(sessionID, key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, ok := s.sessions[sessionID]; !ok {
		s.sessions[sessionID] = make(map[string]interface{})
	}
	s.sessions[sessionID][key] = value
}

func Session() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get or create session ID from cookie
		sessionID, err := c.Cookie("session_id")
		if err != nil || sessionID == "" {
			// Generate simple session ID
			sessionID = generateSessionID()
			c.SetCookie("session_id", sessionID, 3600*24*7, "/", "", false, true)
		}
		
		c.Set("session_id", sessionID)
		c.Next()
	}
}

func GetSession(c *gin.Context, key string) (interface{}, bool) {
	sessionID, _ := c.Get("session_id")
	if sid, ok := sessionID.(string); ok {
		return sessionStore.Get(sid, key)
	}
	return nil, false
}

func SetSession(c *gin.Context, key string, value interface{}) {
	sessionID, _ := c.Get("session_id")
	if sid, ok := sessionID.(string); ok {
		sessionStore.Set(sid, key, value)
	}
}

func generateSessionID() string {
	// Simple session ID generation (in production, use UUID or crypto/rand)
	return "sess_" + randomString(32)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[i%len(letters)]
	}
	return string(b)
}
