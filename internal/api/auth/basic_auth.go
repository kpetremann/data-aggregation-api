package auth

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

const (
	unauthorizedResponse = `{"auth": "unauthorized"}`
	wwwAuthenticate      = "WWW-Authenticate"
	realm                = `Basic realm="restricted"`
)

type authMode int

const (
	NoAuth authMode = iota
	LDAPMode
)

type BasicAuth struct {
	ldapAuth *LDAPAuth
	mode     authMode
}

func NewBasicAuth(mode authMode) BasicAuth { return BasicAuth{mode: mode} }

func (b *BasicAuth) ConfigureLdap(ldap *LDAPAuth) error {
	if ldap == nil {
		return errors.New("LDAP configuration is missing")
	}
	b.ldapAuth = ldap

	return nil
}

func (b *BasicAuth) Wrap(next httprouter.Handle) httprouter.Handle {
	switch b.mode {
	case NoAuth:
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) { next(w, r, ps) }
	case LDAPMode:
		return BasicAuthLDAP(b.ldapAuth, next)
	default:
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			log.Error().Str("authentication issue", "bad server configuration").Send()
			http.Error(w, "authentication issue: bad server configuration", http.StatusInternalServerError)
		}
	}
}

// BasicAuthLDAP is a middleware wrapping the target HTTP HandlerFunc.
// It retrieves BasicAuth credentials and authenticate against LDAP.
func BasicAuthLDAP(ldapAuth *LDAPAuth, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set(wwwAuthenticate, realm)
			http.Error(w, unauthorizedResponse, http.StatusUnauthorized)
			return
		}
		if ok := ldapAuth.AuthenticateUser(username, password); !ok {
			w.Header().Set(wwwAuthenticate, realm)
			http.Error(w, unauthorizedResponse, http.StatusUnauthorized)
			return
		}
		next(w, r, ps)
	}
}
