package plex

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTokenCache(t *testing.T) {
	oldtokenfile := tokenFile
	tokenFile = path.Join(os.TempDir(), "mockTokenFile")
	defer func() { tokenFile = oldtokenfile }()

	Convey("with a missing token file", t, func() {
		os.Remove(tokenFile) // Start with no tokenfile

		Convey("when asked to retreve the token", func() {
			token := Token()

			Convey("the token should be an empty string", func() {
				So(token, ShouldEqual, "")
			})
		})

		Convey("given a token to cache", func() {
			err := cacheToken("ValidToken")

			Convey("the token file should be generated", func() {
				So(err, ShouldBeNil)
				So(exists(tokenFile), ShouldBeTrue)
			})

			Convey("the test token should be able to be read", func() {
				So(err, ShouldBeNil)
				So(Token(), ShouldEqual, "ValidToken")
			})

		})

	})

	Convey("given an empty token file", t, func() {
		f, err := os.Create(tokenFile)
		if err != nil {
			t.Fatal("Unable to create empty token file.")
		}
		f.Close()

		Convey("an empty string should be returned", func() {
			So(Token(), ShouldEqual, "")
		})

	})

	Convey("when instructed to remove an existing token file", t, func() {
		// populate the token file for this test
		if err := cacheToken("ValidToken"); err != nil {
			t.Fatalf("unable to create token file: %v", err)
		}

		err := RemoveCachedToken()

		Convey("the file should be removed", func() {
			So(err, ShouldBeNil)
			So(exists(tokenFile), ShouldBeFalse)
		})

	})

}

func TestMyPlexToken(t *testing.T) {
	oldtokenfile := tokenFile
	tokenFile = path.Join(os.TempDir(), "mockTokenFile")
	defer func() {
		tokenFile = oldtokenfile
		os.Remove(tokenFile)
	}()

	ts := startMyPlexTestServer()
	defer ts.Close()

	Convey("when attempting to authenticate to MyPlex with valid credentials", t, func() {
		err := tokenRequest(Credentials{}, ts.URL+"/users/goodSign_in.xml")

		Convey("the returned token should be cached", func() {
			So(err, ShouldBeNil)
			So(Token(), ShouldEqual, "GoodToken")
		})

	})

	Convey("when attempting to authenticate to MyPlex with invalid credentials", t, func() {
		err := tokenRequest(Credentials{}, ts.URL+"/users/badSign_in.xml")

		Convey("an error should be returned", func() {
			So(err, ShouldResemble, fmt.Errorf("401 Unauthorized"))
		})

	})

}

func startMyPlexTestServer() *httptest.Server {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.EscapedPath() {

		case "/users/goodSign_in.xml":
			w.Header().Set("Content-Type", "text/xml;charset=utf-8")
			fmt.Fprintln(w, `<?xml version="1.0" encoding="UTF-8"?>
                <user>
                    <username>testuser</username>
                    <email>testuser@domain.com</email>
                    <joined-at type="datetime">2010-12-31 12:59:59 UTC</joined-at>
                    <authentication-token>GoodToken</authentication-token>
                </user>`)

		case "/users/badSign_in.xml":
			w.Header().Set("Content-Type", "text/xml;charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, `<?xml version="1.0" encoding="UTF-8"?>
                <errors>
                    <error>Invalid email, username, or password.</error>
                </errors>`)

		}

	}))

	return s
}

func exists(fn string) bool {
	if _, err := os.Stat(fn); err == nil {
		return true
	}
	return false
}
