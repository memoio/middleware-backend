package auth

import (
	"errors"
	"sync"
	"time"

	"golang.org/x/xerrors"
)

type Session struct {
	Nonce     string
	LastLogin int64
	RequestID int64
}

type SessionStore struct {
	sessions map[string]Session
	mutex    sync.Mutex
}

var sessionStore = SessionStore{
	sessions: make(map[string]Session),
}

func (s *SessionStore) AddSession(did, nonce string, timestamp int64) error {
	if timestamp <= time.Now().Add(-1*time.Minute).Unix() {
		return xerrors.Errorf("the request has timed out, please log in within one minute")
	}

	session, ok := s.sessions[did]
	if ok {
		if timestamp <= session.LastLogin {
			return xerrors.Errorf("the current request is later than the latest request")
		}
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.sessions[did] = Session{
		Nonce:     nonce,
		LastLogin: timestamp,
		RequestID: 1,
	}
	return nil
}

func (s *SessionStore) GetSession(did string) (Session, error) {
	session, ok := s.sessions[did]
	if !ok {
		return Session{}, errors.New("cannot find session, please log in first")
	}

	return session, nil
}

func (s *SessionStore) VerifySession(did, nonce string, requestID int64) error {
	session, ok := s.sessions[did]
	if !ok {
		return xerrors.Errorf("cannot find session, please log in first")
	}

	if nonce != session.Nonce {
		return xerrors.Errorf("can't match the nonce, please check your input or log in again")
	}

	if requestID != session.RequestID {
		return xerrors.Errorf("not a sequential request")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.sessions[did] = Session{
		Nonce:     session.Nonce,
		LastLogin: session.LastLogin,
		RequestID: session.RequestID + 1,
	}
	return nil
}
