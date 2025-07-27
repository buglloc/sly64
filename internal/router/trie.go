package router

import (
	"strings"
)

type TrieNode struct {
	children map[string]*TrieNode
	route    *Route
	wildcard bool
}

type RouteTrie struct {
	root *TrieNode
}

func NewRouteTrie() *RouteTrie {
	return &RouteTrie{
		root: &TrieNode{children: make(map[string]*TrieNode)},
	}
}

func (t *RouteTrie) Insert(domain string, route *Route) {
	// Drop "*. " prefix if present to normalize wildcard patterns
	domain = strings.TrimPrefix(domain, "*.")

	labels := splitDomain(domain)
	node := t.root

	for _, label := range labels {
		if node.children[label] == nil {
			node.children[label] = &TrieNode{children: make(map[string]*TrieNode)}
		}
		node = node.children[label]
	}

	node.wildcard = true
	node.route = route
}

func (t *RouteTrie) Find(domain string) *Route {
	labels := splitDomain(domain)
	return findTrieLabel(t.root, labels, 0)
}

func findTrieLabel(node *TrieNode, labels []string, idx int) *Route {
	if node == nil {
		return nil
	}

	if idx == len(labels) {
		if node.route != nil {
			return node.route
		}
		return nil
	}

	if child, ok := node.children[labels[idx]]; ok {
		if route := findTrieLabel(child, labels, idx+1); route != nil {
			return route
		}
	}

	if node.wildcard {
		return node.route
	}

	return nil
}

func splitDomain(domain string) []string {
	dLen := len(domain)
	if dLen == 0 {
		return nil
	}

	if domain[dLen-1] == '.' {
		// Drop trailing "."
		dLen--
	}

	if dLen == 0 {
		return nil
	}

	labels := strings.Split(domain[:dLen], ".")
	for i, j := 0, len(labels)-1; i < j; i, j = i+1, j-1 {
		labels[i], labels[j] = labels[j], labels[i]
	}

	return labels
}
