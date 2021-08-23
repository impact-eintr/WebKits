package erouter

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

type node struct {
	// 当前节点的URL路径
	//
	path      string
	wildChild bool
	nType     nodeType
	maxParams uint8
	priority  []*node
	handle    Handle
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
