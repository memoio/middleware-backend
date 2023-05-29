package auth

import (
	"sync"
	"time"

	"golang.org/x/xerrors"
)

type Identity struct {
	address string
	chainID int
}

type Session struct {
	randomToken string
	lastLogin   int64
	requestID   int64
}

type SessionStore struct {
	sessions map[Identity]Session
	mutex    sync.Mutex
}

var sessionStore = SessionStore{
	sessions: make(map[Identity]Session),
}

func (s *SessionStore) AddSession(address, token string, chainID int, timestamp int64) error {
	identity := Identity{
		address: address,
		chainID: chainID,
	}

	if timestamp <= time.Now().Add(-1*time.Minute).Unix() {
		return xerrors.Errorf("the request has timed out, please log in within one minute")
	}

	session, ok := s.sessions[identity]
	if ok {
		if timestamp <= session.lastLogin {
			return xerrors.Errorf("the current request is later than the latest request")
		}
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.sessions[identity] = Session{
		randomToken: token,
		lastLogin:   timestamp,
		requestID:   0,
	}
	return nil
}

func (s *SessionStore) VerifySession(address, token string, chainID int, requestID int64) error {
	identity := Identity{
		address: address,
		chainID: chainID,
	}
	session, ok := s.sessions[identity]
	if !ok {
		return xerrors.Errorf("cannot find session, please log in first")
	}

	if token != session.randomToken {
		return xerrors.Errorf("can't match the token, please check your input or log in again")
	}
	if requestID != session.requestID+1 {
		return xerrors.Errorf("not a sequential request")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.sessions[identity] = Session{
		randomToken: session.randomToken,
		lastLogin:   session.lastLogin,
		requestID:   session.requestID + 1,
	}
	return nil
}
