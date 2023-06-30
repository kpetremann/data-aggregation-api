package auth

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	unauthorizedResponse = `{"auth": "unauthorized"}`
	wwwAuthenticate      = "WWW-Authenticate"
	realm                = `Basic realm="restricted"`
)

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
