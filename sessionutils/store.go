package sessionutils

import (
	"encoding/base32"
	"errors"
	"github.com/diemus/tablecache"
	"net/http"
	"strings"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type TableCacheStore struct {
	Codecs  []securecookie.Codec
	Options *sessions.Options
	client  *tablecache.TableCache
}

func NewTableCacheStore(client *tablecache.TableCache, keyPairs ...[]byte) *TableCacheStore {
	store := TableCacheStore{
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		Options: &sessions.Options{
			Path:   "/",
			MaxAge: 86400 * 30,
		},
		client: client,
	}
	store.MaxAge(store.Options.MaxAge)
	return &store
}

// Get returns a session for the given name after adding it to the registry.
//
// It returns a new session if the sessions doesn't exist. Access IsNew on
// the session to check if it is an existing session or a new one.
//
// It returns a new session and an error if the session exists but could
// not be decoded.
func (s *TableCacheStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

// New returns a session for the given name without adding it to the registry.
//
// The difference between New() and Get() is that calling New() twice will
// decode the session data twice, while Get() registers and reuses the same
// decoded session after the first call.
func (s *TableCacheStore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(s, name)
	options := *s.Options
	session.Options = &options
	session.IsNew = true

	c, err := r.Cookie(name)
	if err != nil {
		// Cookie not found, this is a new session
		return session, nil
	}

	err = securecookie.DecodeMulti(name, c.Value, &session.ID, s.Codecs...)
	if err != nil {
		// Value could not be decrypted, consider this is a new session
		return session, nil
	}

	v, err := s.client.Get(session.ID)
	if errors.Is(err, tablecache.KeyNotFoundError) {
		// No value found in cache, don't set any values in session object,
		// consider a new session
		return session, nil
	} else if err != nil {
		//something not right
		return session, err
	}

	// Values found in session, this is not a new session
	err = s.load(session, v)
	if err != nil {
		return session, err
	}
	session.IsNew = false
	return session, nil
}

// Save adds a single session to the response.
// Set Options.MaxAge to -1 or call MaxAge(-1) before saving the session to delete all values in it.
func (s *TableCacheStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	// Delete if max-age is <= 0
	if session.Options.MaxAge < 0 {
		//delete from database
		err := s.client.Del(session.ID)
		if err != nil {
			return err
		}

		//delete local session
		for k := range session.Values {
			delete(session.Values, k)
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", s.Options))
		return nil
	}

	if session.ID == "" {
		session.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
	}

	encodedID, err := securecookie.EncodeMulti(session.Name(), session.ID, s.Codecs...)
	if err != nil {
		return err
	}

	encodedValues, err := securecookie.EncodeMulti(session.Name(), session.Values, s.Codecs...)
	if err != nil {
		return err
	}

	err = s.client.Set(session.ID, encodedValues)
	if err != nil {
		return err
	}

	http.SetCookie(w, sessions.NewCookie(session.Name(), encodedID, s.Options))
	return nil
}

// MaxAge sets the maximum age for the store and the underlying cookie
// implementation. Individual sessions can be deleted by setting Options.MaxAge
// = -1 for that session.
func (s *TableCacheStore) MaxAge(age int) {
	s.Options.MaxAge = age

	// Set the maxAge for each securecookie instance.
	for _, codec := range s.Codecs {
		if sc, ok := codec.(*securecookie.SecureCookie); ok {
			sc.MaxAge(age)
		}
	}
}

//load data into session
func (s *TableCacheStore) load(session *sessions.Session, value string) error {
	if err := securecookie.DecodeMulti(session.Name(), value, &session.Values, s.Codecs...); err != nil {
		return err
	}
	return nil
}
