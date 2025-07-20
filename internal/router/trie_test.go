package router_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/buglloc/sly64/v2/internal/router"
)

func TestRouteTrie(t *testing.T) {
	routes := []*router.Route{
		router.NewRoute(router.WithRouteName("1")),
		router.NewRoute(router.WithRouteName("2")),
		router.NewRoute(router.WithRouteName("3")),
		router.NewRoute(router.WithRouteName("4")),
		router.NewRoute(router.WithRouteName("5")),
	}

	trie := router.NewRouteTrie()
	trie.Insert("foo.example.com", routes[0])
	trie.Insert("*.example.com", routes[1])
	trie.Insert("bar.foo.example.com", routes[2])
	trie.Insert("example.com", routes[3])
	trie.Insert("*.wild.com", routes[4])

	tests := []struct {
		domain string
		expect *router.Route
	}{
		{"foo.example.com.", routes[0]},         // exact
		{"bar.example.com.", routes[1]},         // wildcard
		{"baz.foo.example.com.", routes[1]},     // wildcard matches any subdomain
		{"example.com.", routes[3]},             // exact
		{"bar.foo.example.com.", routes[2]},     // exact
		{"baz.bar.foo.example.com.", routes[1]}, // wildcard
		{"com.", nil},                           // no match
		{"some.com.", nil},                      // no match
		{"x.y.wild.com.", routes[4]},            // multi-wildcard
		{"foo.wild.com.", routes[4]},            // wildcard
		{"wild.com.", routes[4]},                // no match
	}

	for _, test := range tests {
		t.Run(test.domain, func(t *testing.T) {
			got := trie.Find(test.domain)
			if test.expect == nil {
				require.Nil(t, got)
				return
			}

			require.Equal(t, test.expect, got, "expected=%s got=%s", test.expect.Name(), got.Name())
		})
	}
}
