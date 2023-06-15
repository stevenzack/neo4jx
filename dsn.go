package neo4jx

import "net/url"

// neo4j+s://fa16ab6c.databases.neo4j.io?password=ASnA99CF7bojkXWgAXNCUr48XmcXGI1-uLaFOkdWlmM
func ParseDsn(dsn string) (addr string, username string, password string, e error) {
	vs, e := url.Parse(dsn)
	if e != nil {
		return "", "", "", e
	}
	addr = vs.Scheme + "://" + vs.Host
	username = vs.Query().Get("username")
	password = vs.Query().Get("password")
	return
}
