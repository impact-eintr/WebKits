package erouter

import (
	"strings"
)

// Priority   Path             Handle
// 9          \                *<1>
// 3          ├s               nil
// 2          |├earch\         *<2>
// 1          |└upport\        *<3>
// 2          ├blog\           *<4>
// 1          |    └:post      nil
// 1          |         └\     *<5>
// 2          ├about-us\       *<6>
// 1          |        └team\  *<7>
// 1          └contact\        *<8>
//
// 这个图相当于注册了下面这几个路由
// GET("/search/", func1)
// GET("/support/", func2)
// GET("/blog/:post/", func3)
// GET("/about-us/", func4)
// GET("/about-us/team/", func5)
// GET("/contact/", func6)

// 通过上面的示例可以看出：
// *<数字> 代表一个 handler 函数的内存地址（指针）
// search 和 support 拥有共同的父节点 s ，并且 s 是没有对应的 handle 的， 只有叶子节点（就是最后一个节点，下面没有子节点的节点）才会注册 handler 。
// 从根开始，一直到叶子节点，才是路由的实际路径。
// 路由搜索的顺序是从上向下，从左到右的顺序，为了快速找到尽可能多的路由，包含子节点越多的节点，优先级越高。

type nodeType uint8

const (
	static nodeType = iota
	root
	param
	catchAll
)

type node struct {
	// 当前节点的 URL 路径
	// 如上面图中的例子的首先这里是一个 /
	// 然后 children 中会有 path 为 [s, blog ...] 等的节点
	// 然后 s 还有 children node [earch,upport] 等，就不再说明了
	path string

	// 和下面的 children 对应，保留的子节点的第一个字符
	// 如上图中的 s 节点，这里保存的就是 eu （earch 和 upport）的首字母
	indices string

	// 判断当前节点路径是不是含有参数的节点, 上图中的 :post 的上级 blog 就是wildChild节点
	wildChild bool

	// 节点类型: static, root, param, catchAll
	// static: 静态节点, 如上图中的父节点 s （不包含 handler 的)
	// root: 如果插入的节点是第一个, 那么是root节点
	// catchAll: 有*匹配的节点
	// param: 参数节点，比如上图中的 :post 节点
	nType nodeType

	// 优先级，查找的时候会用到,表示当前节点加上所有子节点的数目
	priority uint32
	// 当前节点的所有直接子节点
	children []*node
	// 当前节点对应的 handler
	handle Handle
}

// 因为路由是一个基数树，全部是从根节点开始，如果第一次调用注册方法的时候根是不存在的，
// 就注册一个根节点， 这里是每一种请求方法是一个根节点，会存在多个树。
// GET_/
//      \s
//        \earch
//        \upport
//      \blog
//           \:post
// POST_
// addRoute 将传入的 handle 添加到路径中
// 需要注意，这个操作不是并发安全的！！！！
func (n *node) addRoute(path string, handle Handle) {
	fullPath := path
	// 请求到达这个方法 就给当前节点的权重 + 1
	n.priority++

	// 如果树是空的
	if n.path == "" && n.indices == "" {
		// 如果 n 是一个空格的节点，就直接调用插入子节点方法
		n.insertChild(path, fullPath, handle)
		// 并且它只有第一次插入的时候才会是空的，所以将 nType 定义成 root
		n.nType = root
		return
	}

walk:
	for {
		// 先找到最长公共路径长度
		i := logestCommonPrefix(path, n.path)

		// 如果相同前缀的长度比当前节点保存的 path 短
		// 比如  n.path == search ， path == support
		// 它们相同的前缀就变成了 s ， s 比 search 要短，符合 if 的条件，要做处理
		if i < len(n.path) {
			// /_
			//   \search -> handler1
			//
			// /_
			//   \s
			//     \earch -> handler1
			child := node{
				path:      n.path[i:],
				wildChild: n.wildChild,
				// 将类型变更为static 默认没有处理函数的节点
				nType: static,
				// earch 继承 s 的indices
				indices: n.indices,
				// earch 继承 s 的子节点
				children: n.children,
				// earch 继承 s 的处理函数
				handle: n.handle,
				// 子节点(earch)优先级继承自父节点 并且-1
				priority: n.priority - 1,
			}
			// 更新节点信息
			n.children = []*node{&child}
			// 获取子节点的首字母,因为上面分割的时候是从 i 的位置开始分割
			// 所以 n.path[i] 可以去除子节点的首字母，理论上去 child.path[0] 也是可以的
			// 这里的 n.path[i] 取出来的是一个 uint8 类型的数字（代表字符），
			// 先用 []byte 包装一下数字再转换成字符串格式
			n.indices = string([]byte{n.path[i]})
			n.path = path[:i]
			// 变成一个没有处理函数的节点
			n.handle = nil
			// 肯定没有参数了，已经变成了一个没有 handle 的节点了
			n.wildChild = false
		}

		// 将新的节点添加到此节点的子节点， 这里是新添加节点的子节点
		// /_
		//   \abc -> handler1

		// abc/def

		// /_
		//   \abc -> handler1
		//   \(:abc) -> handler1
		//          \def -> handler2
		if i < len(path) {
			path = path[i:] // /def

			// 如果当前路径有参数
			// 就是定义路由时候是这种形式的： blog/:post/update
			// 如果进入了上面 if i < len(n.path) 这个条件，这里就不会成立了
			// 因为上一个 if 中将 n.wildChild 重新定义成了 false
			if n.wildChild {
				// 如果进入到了这里，证明这是一个参数节点，类似 :post 这种
				// 不会这个节点进行处理，直接将它的子节点赋值给当前节点
				// 比如： :post/ ，只要是参数节点，必有子节点，哪怕是
				// blog/:post 这种，也有一个 / 的子节点
				n = n.children[0]
				n.priority++ // 子节点 喜加一

				// 检查通配符是否匹配
				// 这里的 path 已经变成了去除了公共前缀的后面部分，比如
				// :abc/def ， 就是 /def
				// 这里的 n 也已经是 :abc 这种的下一级的节点，比如 / 或者 /d 等等
				// 如果添加的节点的 path >= 当前节点的 path &&
				// 当前节点的 path 长度和添加节点的前面相同数量的字符是相等的
				if len(path) >= len(n.path) && n.path == path[:len(n.path)] &&
					// 添加一个catchAll的子节点是不可能的
					n.nType != catchAll &&
					// 当前节点的 path >= 添加节点的 path ，其实有第一个条件限制，
					// 这里也只有 len(n.path) == len(path) 才会成立，
					// 就是当前节点的 path 和 添加节点的 path 相等 ||
					// 添加节点的 path 减去当前节点的 path 之后是 /
					// 例如： n.path = name, path = name 或
					// n.path = name, path = name/ 这两种情况
					(len(n.path) >= len(path) || path[len(n.path)] == '/') {
					// 跳出当前循环，进入下一次循环
					// 再次循环的时候
					// 1. if i < len(n.path) 这里就不会再进入了，现在 i == len(n.path)
					// 2. if n.wildChild 也不会进入了，
					// 当前节点已经在上次循环的时候改为 children[0]
					continue walk
				} else {
					// 当不是 n.path = name, path = name/ 这两种情况的时候，
					// 代表通配符冲突了，什么意思呢？
					// 简单的说就是通配符部分只允许定义相同的或者 / 结尾的
					// 例如：blog/:post/update，再定义一个路由 blog/:postabc/add，
					// 这个时候就会冲突了，是不被允许的，blog 后面只可以定义
					// :post 或 :post/ 这种，同一个位置不允许使用多种通配符
					// 这里的处理是直接 panic 了，如果想要支持，可以尝试重写下面部分代码
					// 下面做的事情就是组合 panic 用到的提示信息
					var pathSeg string

					// 如果当前节点的类型是有*匹配的节点
					if n.nType == catchAll {
						pathSeg = path
					} else {
						// 如果不是，将 path 做字符串分割
						// 这个是通过 / 分割，最多分成两个部分,然后取第一部分的值
						// 例如： path = "name/hello/world"
						// 分割两部分就是 name 和 hello/world , pathSeg = name
						pathSeg = strings.SplitN(path, "/", 2)[0]
					}

					// 通过传入的原始路径来处理前缀, 可以到上面看下，方法进入就定义了这个变量
					// 在原始路径中提取出 pathSeg 前面的部分在拼接上 n.path
					// 例如： n.path = ":post" , fullPath="/blog/:postnew/add"
					// 这时的 prefix = "/blog/:post"
					prefix := fullPath[:strings.Index(fullPath, pathSeg)] + n.path

					// 最终的提示信息就会生成类似这种：
					// panic: ':postnew' in new path '/blog/:postnew/update/' \
					// conflicts with existing wildcard ':post' in existing \
					// prefix '/blog/:post'
					// 就是说已经定义了 /blog/:post 这种规则的路由，
					// 再定义 /blog/:postnew 这种就不被允许了
					panic("'" + pathSeg +
						"' in new path '" + fullPath +
						"' conflicts with existing wildcard '" + n.path +
						"' in existing prefix '" + prefix +
						"'")
				}
			}

			// 如果没有进入到上面的参数节点，当前节点不是一个参数节点 :post 这种
			idxc := path[0] // indexchar

			if n.nType == param && idxc == '/' && len(n.children) == 1 {
				// /:post 这种节点不做处理，直接拿这个节点的子节点去匹配
				n = n.children[0]
				// 权重 + 1 ， 因为新的节点会变成这个节点的子节点
				n.priority++
				// 结束当前循环 再次进行匹配
				continue walk
			}

			// 检查添加的 path 的首字母是否保存在在当前节点的 indices 中
			for i, c := range []byte(n.indices) {
				if c == idxc {
					// 这里处理优先级和排序的问题，把这个方法看完再去查看这个方法干了什么
					i = n.increamentChildPrio(i)
					// 将当前的节点替换成它对应的子节点
					n = n.children[i]
					continue walk
				}
			}

			// 如果上面 for 中也没有匹配上，就将新添加的节点插入
			if idxc != ':' && idxc != '*' {
				n.indices += string([]byte{idxc})
				child := &node{}
				n.children = append(n.children, child)
				n.increamentChildPrio(len(n.indices) - 1)
				n = child
			}

			// 用当前节点发起插入子节点的动作
			// 注意这个 n 已经替换成了上面新初始化的 child 了，相当于是一个空的节点。
			n.insertChild(path, fullPath, handle)
			return
		}

		if n.handle != nil {
			panic("a handle is already registered for path '" + fullPath + "'")
		}
		n.handle = handle
		// 这个新的节点被添加了， 出现了 return ， 只有出现这个才会正常退出循环，一次添加完成。
		return
	}
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

// 返回最长的公共前缀长度
func logestCommonPrefix(a, b string) int {
	i := 0
	max := min(len(a), len(b))
	for i < max && a[i] == b[i] {
		i++
	}
	return i
}

// 通过之前两次的调用，我们知道，这个 pos 都是 n.indices 中指定字符的索引，也就是位置
func (n *node) increamentChildPrio(pos int) int {
	// 因为 children 和 indices 是同时添加的，所以索引是相同的 🤔
	// 可以通过 pos 代表的位置找到， 将对应的子节点的优先级 + 1
	cs := n.children
	cs[pos].priority++
	prio := cs[pos].priority

	// 重新排序
	newPos := pos
	for ; newPos > 0 && cs[newPos-1].priority < prio; newPos-- {
		cs[newPos-1], cs[newPos] = cs[newPos], cs[newPos-1]
	}
	// 重构idxc
	if newPos != pos {
		n.indices = n.indices[:newPos] + n.indices[pos:pos+1] + n.indices[newPos:pos] + n.indices[pos+1:]
	}
	return newPos
}

// 寻找通配符
func findWildcard(path string) (wildcard string, i int, valid bool) {
	// Find start
	for start, c := range []byte(path) {
		if c != ':' && c != '*' {
			continue
		}
		valid = true
		for end, c := range []byte(path[start+1:]) {
			switch c {
			case '/':
				return path[start : start+1+end], start, valid
			case ':', '*':
				valid = false
			}
		}
		return path[start:], start, valid
	}
	return "", -1, false
}

// path 插入的子节点的路径
// fullPath 完整路径，就是注册路由时候的路径，没有被处理过的
// 注册路由对应的 handle 函数
func (n *node) insertChild(path, fullPath string, handle Handle) {
	for {
		wildcard, i, valid := findWildcard(path)
		if i < 0 { // 没有通配符
			break
		}

		// 通配符非法
		if !valid {
			panic("每个路径只允许有一个通配符 " + wildcard + "in path '" + fullPath + "'")
		}

		if len(wildcard) < 2 {
			panic("通配符路由必须有明确的名字 '" + fullPath + "'")
		}

		// 检查通配符所在的位置，是否已经有子节点，如果有，就不能再插入
		// 例如： 已经定义了 /hello/name ， 就不能再定义 /hello/:param
		if len(n.children) > 0 {
			panic("该节点已经有子路由了 不支持继续添加通配符路由 " + wildcard + " " + fullPath)
		}

		// 正式开始匹配
		if wildcard[0] == ':' {
			if i > 0 {
				n.path = path[:i]
				path = path[i:]
			}

			// 标记上当前这个节点是一个包含参数的节点的节点
			n.wildChild = true
			// 将参数部分定义成一个子节点
			child := &node{
				nType: param, // 指定为通配符类型
			}
			// 用新定义的子节点初始化一个children属性
			n.children = []*node{child}

			// 将新创建的节点定义为当前节点，这个要想一下，到这里这种操作已经有不少了
			// 因为一直都是指针操作，修改都是指针的引用，所以定义好的层级关系不会被改变
			//type node struct {
			//        Count int
			//        Child *node
			//}
			//func (n *node) Test() {
			//        n.Count++
			//        child := &node{
			//                Count: n.Count,
			//        }
			//        n.Child = child
			//
			//        fmt.Printf("%p %v\n", n, n)
			//        n = child 方法接收者是n的地址复制 可以通过这个地址复制修改n的属性 但引用不是n
			//        fmt.Printf("%p %v\n", n, n)
			//}
			//func main() {
			//        n := new(node)
			//        n.Count = 1
			//        fmt.Printf("%p %v\n", n, n)
			//        n.Test()
			//        fmt.Printf("%p %v\n", n, n)
			//}
			n = child
			n.priority++

			// 如果小于路径的最大长度，代表还包含子路径（也就是说后面还有子节点）
			if len(wildcard) < len(path) {
				path = path[len(wildcard):]
				// 定义一个子节点，无论后面还有没有子节点 :name 这种格式的路由后面至少还有一个 /
				child := &node{
					priority: 1,
				}
				n.children = []*node{child}
				// 继续向后循环
				n = child
				continue
			}

			// 否则就结束循环 把处理函数嵌入新的叶子节点
			n.handle = handle
			return
		}

		// catchAll 注意 我们这里所说的路径指的是 GET("path", handler) 不是req.URL.Path
		// 这里的意思是， * 匹配的路径只允许定义在路由的最后一部分
		// 比如 : /hello/*world 是允许的， /hello/*world/more 这种就会 painc
		// 这种路径就是会将 hello/ 后面的所有内容变成 world 的变量
		// 比如地址栏输入： /hello/one/two/more ，获取到的参数 world = one/twq/more
		// 不会再将后面的 / 作为路径处理了
		if i+len(wildcard) != len(path) {
			panic("* 匹配的路径只允许定义在路由的最后一部分 " + wildcard + " " + fullPath)
		}

		// 这种情况是，新定义的 * 通配符路由和其他已经定义的路由冲突了 len(n.path)
		// 例如已经定义了一个 /hello/bro ， 又定义了一个 /hello/*world ，此时就会 panic 了
		if len(n.path) > 0 && n.path[len(n.path)-1] == '/' {
			panic("新定义的 * 通配符路由和其他已经定义的路由冲突了 " + wildcard + " " + fullPath)
		}

		// 这里是查询通配符前面是否有 / 没有 / 是不行的，panic
		i-- // 通配符前一个位置
		if path[i] != '/' {
			panic("no / before catch-all in path " + fullPath)
		}

		// 后面的套路基本和之前看到的类似，就是定义一个子节点，保存通配符前面的路径，
		// 有变化的就是将 nType 定义为 catchAll，就是说代表这是一个  * 号匹配的路由
		n.path = path[i:]
		child := &node{
			wildChild: true,
			nType:     catchAll,
		}
		n.children = []*node{child}
		n.indices = string('/')
		n = child
		n.priority++

		// 将下面的节点再添加到上面，不过 * 号路由不会再有下一级的节点了，因为它会将后面的
		// 的所有内容当做变量，即使它是个 / 符号
		child = &node{
			path:     path[i:],
			nType:    catchAll,
			handle:   handle,
			priority: 1,
		}
		n.children = []*node{child}
		return

	}
	n.path = path
	n.handle = handle

}

func countParams(path string) uint16 {
	var n uint
	for i := range []byte(path) {
		switch path[i] {
		case ':', '*':
			n++
		}
	}
	return uint16(n)
}
