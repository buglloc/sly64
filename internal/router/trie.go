package router

import (
	"strings"
)

type TrieNode struct {
	children map[string]*TrieNode
	route    *Route
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
	domain = strings.TrimRight(domain, ".") + "."
	labels := splitDomain(domain)
	node := t.root
	for _, label := range labels {
		if node.children[label] == nil {
			node.children[label] = &TrieNode{children: make(map[string]*TrieNode)}
		}
		node = node.children[label]
	}

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

		if child, ok := node.children["*"]; ok && child.route != nil {
			return child.route
		}

		return nil
	}

	label := labels[idx]
	if child, ok := node.children[label]; ok {
		if route := findTrieLabel(child, labels, idx+1); route != nil {
			return route
		}
	}

	if child, ok := node.children["*"]; ok {
		return child.route
	}

	return nil
}

func splitDomain(domain string) []string {
	labels := strings.Split(domain, ".")
	for i, j := 0, len(labels)-1; i < j; i, j = i+1, j-1 {
		labels[i], labels[j] = labels[j], labels[i]
	}
	return labels
}
