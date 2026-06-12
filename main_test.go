package bfstree_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/bcicen/bfstree"
)

var tree *bfstree.BFSTree

type TestEdge struct {
	from string
	to   string
}

// TestEdge implements bfstree.Edge interface
func (t TestEdge) To() string   { return t.to }
func (t TestEdge) From() string { return t.from }

func TestCreateTree(t *testing.T) {
	tree = bfstree.New(
		TestEdge{"a", "b"},
		TestEdge{"b", "d"},
		TestEdge{"d", "e"},
		TestEdge{"e", "f"},
		TestEdge{"f", "g"},
		TestEdge{"c", "d"},
		TestEdge{"c", "a"},
		TestEdge{"a", "c"},
		TestEdge{"b", "c"},
	)
	t.Logf("created tree with %d edges, %d nodes", tree.Len(), len(tree.Nodes()))
}

func TestFindLongPath(t *testing.T) {
	path, err := tree.FindPath("a", "g")
	if err != nil {
		t.Error(err)
	}
	t.Logf("found path: %s", path)
}

func TestFindShortPath(t *testing.T) {
	path, err := tree.FindPath("a", "b")
	if err != nil {
		t.Error(err)
	}
	t.Logf("found path: %s", path)
}

func TestFindNoPath(t *testing.T) {
	_, err := tree.FindPath("a", "z")
	if err == nil {
		t.Errorf("no error returned on missing path")
	}
	t.Logf("got expected error: %s", err)
}

// TestFindNoPathInDenseGraph guards against combinatorial path blowup: a densely
// connected graph has super-linearly many simple paths between nodes. Without a
// global visited set, FindPath enumerates all of them for an unreachable target
// and never returns. The search must terminate with an error.
func TestFindNoPathInDenseGraph(t *testing.T) {
	const n = 300
	var edges []bfstree.Edge
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i != j {
				edges = append(edges, TestEdge{strconv.Itoa(i), strconv.Itoa(j)})
			}
		}
	}
	dense := bfstree.New(edges...)

	done := make(chan error, 1)
	go func() {
		_, err := dense.FindPath("0", "unreachable")
		done <- err
	}()

	select {
	case err := <-done:
		if err == nil {
			t.Errorf("expected error for unreachable target")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("FindPath did not terminate: combinatorial path blowup on a dense graph")
	}
}
