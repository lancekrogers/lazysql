package tree

import (
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/rivo/tview"
)

func Search(root *tview.TreeNode, searchText string) []*tview.TreeNode {
	if root == nil {
		return nil
	}
	lowerSearchText := strings.ToLower(searchText)
	if lowerSearchText == "" {
		return nil
	}

	parts := strings.SplitN(lowerSearchText, " ", 2)
	databaseNameFilter := ""
	tableNameFilter := ""

	if len(parts) == 1 {
		tableNameFilter = parts[0]
	} else {
		databaseNameFilter = parts[0]
		tableNameFilter = parts[1]
	}

	type rankedNode struct {
		node *tview.TreeNode
		rank int
	}
	var rankedNodes []rankedNode

	root.Walk(func(node, parent *tview.TreeNode) bool {
		nodeText := strings.ToLower(node.GetText())

		if databaseNameFilter == "" {
			rank := fuzzy.RankMatch(tableNameFilter, nodeText)
			if rank >= 0 {
				if parent != nil {
					parent.SetExpanded(true)
				}
				adjustedRank := prioritizeResult(tableNameFilter, nodeText, rank)
				rankedNodes = append(rankedNodes, rankedNode{node: node, rank: adjustedRank})
			}
		} else {
			rank := fuzzy.RankMatch(tableNameFilter, nodeText)
			if rank >= 0 && parent != nil {
				parentText := strings.ToLower(parent.GetText())
				parentRank := fuzzy.RankMatch(databaseNameFilter, parentText)
				if parentRank >= 0 {
					parent.SetExpanded(true)
					adjustedTableRank := prioritizeResult(tableNameFilter, nodeText, rank)
					adjustedParentRank := prioritizeResult(databaseNameFilter, parentText, parentRank)
					combinedRank := adjustedTableRank + (adjustedParentRank / 2)
					rankedNodes = append(rankedNodes, rankedNode{node: node, rank: combinedRank})
				}
			}
		}

		return true
	})

	sort.Slice(rankedNodes, func(i, j int) bool {
		return rankedNodes[i].rank < rankedNodes[j].rank
	})

	results := make([]*tview.TreeNode, 0, len(rankedNodes))
	for _, rn := range rankedNodes {
		results = append(results, rn.node)
	}
	return results
}

func prioritizeResult(pattern, target string, fuzzyRank int) int {
	if pattern == target {
		return 0
	}

	if strings.HasPrefix(target, pattern) {
		lengthDiff := len(target) - len(pattern)
		if lengthDiff > 98 {
			lengthDiff = 98
		}
		return 1 + lengthDiff
	}

	if strings.Contains(target, pattern) {
		index := strings.Index(target, pattern)
		lengthPenalty := len(target) - len(pattern)
		score := 100 + index + lengthPenalty
		if score > 9999 {
			score = 9999
		}
		return score
	}

	return 10000 + fuzzyRank
}
