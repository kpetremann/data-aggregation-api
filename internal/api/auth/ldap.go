package auth

import (
	"crypto/tls"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/go-ldap/ldap/v3"
)

type LDAPAuth struct {
	ldapClient *ldap.Conn
	bindDN     string
	password   string
	baseDN     string
}

func NewLDAPAuth(ldapURL string, bindDN string, password string, baseDN string, tlsConfig *tls.Config) *LDAPAuth {
	conn, err := ldap.DialURL(ldapURL, ldap.DialWithTLSConfig(tlsConfig))
	if err != nil {
		log.Error().Err(err).Str("ldapURL", ldapURL).Msg("failed to connect to the LDAP server")
		return nil
	}
	return &LDAPAuth{
		ldapClient: conn,
		bindDN:     bindDN,
		password:   password,
		baseDN:     baseDN,
	}
}

func (l *LDAPAuth) AuthenticateUser(username string, password string) bool {
	if err := l.ldapClient.Bind(l.bindDN, l.password); err != nil {
		log.Error().Err(err).Str("bindDN", l.bindDN).Msg("failed to bind to LDAP")
		return false
	}
	search, err := l.ldapClient.Search(ldap.NewSearchRequest(
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
		return false
	}
	if len(search.Entries) != 1 {
		log.Error().Str(userKey, username).Msg("no result or more than 1 result found for user")
		return false
	}
	if err := l.ldapClient.Bind(search.Entries[0].DN, password); err != nil {
		log.Error().Err(err).Str(userKey, username).Msg("failed to bind with user")
		return false
	}
	log.Debug().Str(userKey, username).Msg("succesfully authenticated user")
	return true
}
