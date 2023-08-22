package auth

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/go-ldap/ldap/v3"
)

type connectionStatus bool

const (
	connectionClosed connectionStatus = false
	connectionUp     connectionStatus = true
)

type authRequest struct {
	authResp chan bool
	username string
	password string
}

type result struct {
	auth bool
	conn connectionStatus
}

type LDAPAuth struct {
	tlsConfig *tls.Config
	reqCh     chan authRequest
	ldapURL   string
	bindDN    string
	password  string
	baseDN    string
}

func NewLDAPAuth(ldapURL string, bindDN string, password string, baseDN string, tlsConfig *tls.Config) *LDAPAuth {
	return &LDAPAuth{
		tlsConfig: tlsConfig,
		ldapURL:   ldapURL,
		bindDN:    bindDN,
		password:  password,
		baseDN:    baseDN,
		reqCh:     make(chan authRequest),
	}
}

func SetDefaultTimeout(timeout time.Duration) {
	ldap.DefaultTimeout = timeout //nolint:reassign  // we want to customize the default timeout
}

func (l *LDAPAuth) StartWorkers(ctx context.Context, maxWorker int) error {
	if maxWorker <= 0 {
		return fmt.Errorf("maxWorker must be greater than 0")
	}
	for i := 0; i < maxWorker; i++ {
		go l.spawnWorker(ctx)
	}
	return nil
}
func (l *LDAPAuth) AuthenticateUser(username string, password string) bool {
	req := authRequest{
		username: username,
		password: password,
		authResp: make(chan bool),
	}
	l.reqCh <- req
	return <-req.authResp
}

func (l *LDAPAuth) spawnWorker(ctx context.Context) {
	const maxRetry = 1
	var conn *ldap.Conn
	var err error
	tick := time.NewTicker(time.Minute)

	for {
		select {
		case req := <-l.reqCh:
			auth := false
			retry := 0
			for retry <= maxRetry+1 {
				retry++
				log.Debug().Msgf("worker LDAP authentication attempt number %d", retry)
				// (re)connect if needed
				if conn == nil || conn.IsClosing() {
					log.Debug().Msg("LDAP connection is closed, reconnecting")
					conn, err = l.connect()
					if err != nil {
						log.Error().Err(err).Msg("worker LDAP reconnection failed")
						req.authResp <- false
						break
					}
				}

				// bind with the user credentials
				var connState connectionStatus
				auth, connState = l.authenticateWithTimeout(ctx, conn, req.username, req.password, ldap.DefaultTimeout)

				if connState == connectionClosed {
					log.Debug().Msg("LDAP connection was closed by the server, closing on client side")
					if err := conn.Close(); err != nil {
						log.Error().Err(err).Msg("connection was closed by the server but failed to close on client side")
					}
				} else {
					// LDAP connection is still up, we accept the authentication result
					log.Debug().Msg("auth response valid")
					break
				}
			}

			log.Debug().Msgf("worker LDAP authentication attempt number %d, result: %t", retry, auth)

			req.authResp <- auth
			tick.Reset(time.Minute)

		case <-tick.C:
			// close connection if no request has been made for a minute
			log.Debug().Msg("timer reached, closing connection")
			if conn != nil {
				if err := conn.Close(); err != nil {
					log.Error().Err(err).Msg("unable to close the LDAP connection")
				}
			}
			return

		case <-ctx.Done():
			// gracefully close connection if context is done
			log.Debug().Msg("context is closed, closing connection")
			closeLDAPConnection(conn)
			return
		}
	}
}

func closeLDAPConnection(conn *ldap.Conn) {
	if conn != nil && !conn.IsClosing() {
		_ = conn.Close()
	}
}

// authenticateWithTimeout performs the authentication against LDAP with a timeout.
func (l *LDAPAuth) authenticateWithTimeout(ctx context.Context, conn *ldap.Conn, username, password string, timeout time.Duration) (bool, connectionStatus) {
	// request the authentication
	res := make(chan result)
	go func() {
		a, c := l.authenticate(conn, username, password)
		res <- result{auth: a, conn: c}
	}()

	var connState connectionStatus
	var auth bool

	// handle timeout and context closing
	select {
	case r := <-res:
		auth = r.auth
		connState = r.conn
	case <-time.After(timeout):
		log.Error().Msg("LDAP authentication timeout")
		closeLDAPConnection(conn)
		auth = false
		connState = connectionClosed
	case <-ctx.Done():
		auth = false
	}
	return auth, connState
}

func (l *LDAPAuth) connect() (*ldap.Conn, error) {
	conn, err := ldap.DialURL(l.ldapURL, ldap.DialWithTLSConfig(l.tlsConfig))
	if err != nil {
		return nil, fmt.Errorf("failed to dial to the LDAP server: %w", err)
	}

	return conn, nil
}

// authenticate performs the authentication against LDAP.
// The first returned boolean is true if the authentication is successful, false otherwise.
// The second returned boolean is true if the connection is closed, false otherwise.
func (l *LDAPAuth) authenticate(conn *ldap.Conn, username string, password string) (bool, connectionStatus) {
	if err := conn.Bind(l.bindDN, l.password); err != nil {
		log.Error().Err(err).Str("bindDN", l.bindDN).Msg("failed to bind to LDAP")

		// detect TCP connection closed or any network errors
		if ldap.IsErrorWithCode(err, ldap.ErrorNetwork) {
			return false, connectionClosed
		}
		return false, connectionUp
	}
	search, err := conn.Search(ldap.NewSearchRequest(
		l.baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=user)(sAMAccountName=%s))", username),
		[]string{"dn"},
		nil,
	))

	const userKey = "user"
	if err != nil {
		log.Error().Err(err).Str(userKey, username).Msg("failed to perform LDAP search to find user")
		return false, connectionUp
	}
	if len(search.Entries) != 1 {
		log.Error().Str(userKey, username).Msg("no result or more than 1 result found for user")
		return false, connectionUp
	}
	if err := conn.Bind(search.Entries[0].DN, password); err != nil {
		log.Error().Err(err).Str(userKey, username).Msg("failed to bind with user")
		return false, connectionUp
	}
	return true, connectionUp
}
