package auth

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"

	"github.com/criteo/data-aggregation-api/internal/config"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

const (
	unauthorizedResponse = `{"auth": "unauthorized"}`
	wwwAuthenticate      = "WWW-Authenticate"
	realm                = `Basic realm="restricted"`
)

type authMode string

const (
	noAuth   authMode = "None"
	ldapMode authMode = "LDAP"
)

type BasicAuth struct {
	ldapAuth *LDAPAuth
	mode     authMode
}

func NewBasicAuth(cfg config.AuthConfig) (BasicAuth, error) {
	b := BasicAuth{mode: noAuth}

	if cfg.LDAP == nil {
		return b, nil
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: cfg.LDAP.InsecureSkipVerify, //nolint:gosec // configurable on purpose
	}
	ldap := NewLDAPAuth(cfg.LDAP.URL, cfg.LDAP.BindDN, cfg.LDAP.Password, cfg.LDAP.BaseDN, tlsConfig)
	if err := b.configureLdap(ldap); err != nil {
		return b, fmt.Errorf("failed to configure the request authenticator: %w", err)
	}
	b.mode = ldapMode

	return b, nil
}

func (b *BasicAuth) configureLdap(ldap *LDAPAuth) error {
	if ldap == nil {
		return errors.New("LDAP configuration is missing")
	}
	b.ldapAuth = ldap

	return nil
}

func (b *BasicAuth) Wrap(next httprouter.Handle) httprouter.Handle {
	switch b.mode {
	case noAuth:
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) { next(w, r, ps) }
	case ldapMode:
		return BasicAuthLDAP(b.ldapAuth, next)
	default:
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			log.Error().Str("auth-method", string(b.mode)).Str("authentication issue", "bad server configuration").Send()
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
