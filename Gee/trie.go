package Gee

import (
	"strings"
)

// 路由前缀树
type node struct {
	pattern  string  // 待匹配路由
	part     string  // 当前路由的部分
	childern []*node // 子节点
	isWild   bool    // 是否精确匹配
}

// 第一个匹配的子节点
func (n *node) matchChild(part string) *node {
	for _, child := range n.childern {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配的节点，用于查找
func (n *node) matchChildern(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.childern {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 注册路由(添加新的路由规则)
func (n *node) insert(pattern string, parts []string, height int) {
	// 路径完全匹配
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]       // 每次匹配的节点
	child := n.matchChild(part) // 是否有对应节点
	if child == nil {           // 添加新节点
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*', // 动态路由
		}
		n.childern = append(n.childern, child) // 加入子节点列表
	}

	// 继续向下匹配
	child.insert(pattern, parts, height+1)
}

// 查找路由(寻找匹配的路径)
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	childern := n.matchChildern(part)
	for _, child := range childern {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
