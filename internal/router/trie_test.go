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
		router.NewRoute(router.WithRouteName("6")),
	}

	trie := router.NewRouteTrie()
	// Both wildcard and non-wildcard syntax work the same way
	// "*." prefix is dropped during insertion
	trie.Insert("foo.example.com", routes[0])     // stored as foo.example.com
	trie.Insert("*.example.com", routes[1])       // stored as example.com
	trie.Insert("bar.foo.example.com", routes[2]) // stored as bar.foo.example.com
	trie.Insert("example.com", routes[3])         // stored as example.com (overwrites routes[1])
	trie.Insert("wild.com", routes[4])            // stored as wild.com
	trie.Insert(".", routes[5])                   // stored as .

	tests := []struct {
		domain string
		expect *router.Route
	}{
		{"foo.example.com.", routes[0]},         // exact match
		{"bar.example.com.", routes[3]},         // matches example.com pattern
		{"baz.foo.example.com.", routes[0]},     // matches foo.example.com pattern
		{"example.com.", routes[3]},             // exact match
		{"bar.foo.example.com.", routes[2]},     // exact match
		{"baz.bar.foo.example.com.", routes[2]}, // matches bar.foo.example.com pattern
		{"com.", routes[5]},                     // match global
		{"some.com.", routes[5]},                // match global
		{"x.y.wild.com.", routes[4]},            // matches wild.com pattern
		{"foo.wild.com.", routes[4]},            // matches wild.com pattern
		{"wild.com.", routes[4]},                // exact match
		{"wild.com", routes[4]},                 // exact match
	}

	for _, test := range tests {
		t.Run(test.domain, func(t *testing.T) {
			got := trie.Find(test.domain)
			if test.expect == nil {
				require.Nil(t, got)
				return
			}

			require.NotNil(t, got)
			require.Equal(t, test.expect, got, "expected=%s got=%s", test.expect.Name(), got.Name())
		})
	}
}
