package main

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// 哈希表记录下标+递归
func buildTree(preorder []int, inorder []int) *TreeNode {
	var (
		n   = len(preorder)
		idx = make(map[int]int)
		dfs func(pL, pR, iL, iR int) *TreeNode // 左闭右开
	)
	for i, x := range inorder {
		idx[x] = i
	}

	dfs = func(pL, pR, iL, iR int) *TreeNode {
		if pL == pR {
			return nil
		}

		i := idx[preorder[pL]]
		lL := i - iL
		left := dfs(pL+1, pL+1+lL, iL, i)
		right := dfs(pL+1+lL, pR, i+1, iR)
		return &TreeNode{
			Val:   preorder[pL],
			Left:  left,
			Right: right,
		}
	}

	return dfs(0, n, 0, n)
}
