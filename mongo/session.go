package mongo

import (
	"gopkg.in/mgo.v2"
)

// Session holds a Mongo connection session.
type Session struct {
	session *mgo.Session
}

// NewSession creates a new Session towards the requested endpoint.
func NewSession(url string) (*Session, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	return &Session{session}, err
}

// Copy creates a new Session instance towards the same endpoint.
func (s *Session) Copy() *Session {
	return &Session{s.session.Copy()}
}

// GetCollection retrieves the requested collection from the database.
func (s *Session) GetCollection(db string, col string) *mgo.Collection {
	return s.session.DB(db).C(col)
}

// Close finishes a Session connection.
func (s *Session) Close() {
	if s.session != nil {
		s.session.Close()
	}
}

// DropDatabase drops the requestes database.
func (s *Session) DropDatabase(db string) error {
	if s.session != nil {
		return s.session.DB(db).DropDatabase()
	}
	return nil
}
