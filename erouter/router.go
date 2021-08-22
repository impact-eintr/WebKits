package erouter

import (
	"net/http"
	"sync"
)

type Param struct {
	key   string
	value string
}

// Params 就是一个 Param 的切片，这样就可以看出来， URL 参数可以设置多个了。
// 它是在 tree 的 GetValue() 方法调用的时候设置的
// 这个切片是有顺序的，第一个设置的参数就是切片的第一个值，所以通过索引获取值是安全的
type Params []Param

func (ps Params) ByName(name string) string {
	for _, p := range ps {
		if p.key == name {
			return p.value
		}
	}
	return ""
}

type paramsKey struct{}

var ParamsKey = paramsKey{}

type Handle func(http.ResponseWriter, *http.Request, Params)

type Router struct {
	trees map[string]*node

	paramsPool sync.Pool
	maxParams  uint16

	SaveMatchedRoutePath bool

	// 这个参数是否自动处理当访问路径最后带的 /，一般为 true 就行。
	// 例如： 当访问 /foo/ 时， 此时没有定义 /foo/ 这个路由，但是定义了
	// /foo 这个路由，就对自动将 /foo/ 重定向到 /foo (GET 请求
	// 是 http 301 重定向，其他方式的请求是 http 307 重定向）。
	RedirectTrailingSlash bool

	// 是否自动修正路径， 如果路由没有找到时，Router 会自动尝试修复。
	// 首先删除多余的路径，像 ../ 或者 // 会被删除。
	// 然后将清理过的路径再不区分大小写查找，如果能够找到对应的路由， 将请求重定向到
	// 这个路由上 ( GET 是 301， 其他是 307 ) 。
	RedirectFixedPath bool

	HandleMethodNotAllowed bool

	// 如果为 true ，会自动回复 OPTIONS 方式的请求。
	// 如果自定义了 OPTIONS 路由，会使用自定义的路由，优先级高于这个自动回复。
	HandleOPTHONS bool

	GlobalOPTHONS http.Handler
	GlobalAllowed string

	NotFound         http.Handler
	MethodNotAllowed http.Handler
	PanicHandler     func(http.ResponseWriter, *http.Request, interface{})
}

func New() *Router {
	return &Router{}
}

// 验证Router是否实现了http.Handler
var _ http.Handler = New()

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r.PanicHandler != nil {
		defer r.recv(w, req)
	}
	// 请求路径
	path := req.URL.Path

	// 到基数树中去查找匹配的路由
	if root := r.trees[req.Method]; root != nil {
		// 如果路由成功匹配 从路由从基数树中取出
		if handle, ps, tsr := root.getValue(path, r.getParams); handle != nil {
			if ps != nil {
				handle(w, req, *ps) // 向httprouter注册的函数
				r.putParams(ps)
			} else {
				handle(w, req, nil)
			}
			return // 此次生命周期结束
		} else if req.Method != "CONNECT" && path != "/" {
			// 在 HTTP 协议中，CONNECT 方法可以开启一个客户端与所请求资源之间的双向沟通的通道。
			// 它可以用来创建隧道（tunnel）。
			// 例如，CONNECT 可以用来访问采用了 SSL (en-US) (HTTPS)  协议的站点。
			// 客户端要求代理服务器将 TCP 连接作为通往目的主机隧道。
			// 之后该服务器会代替客户端与目的主机建立连接。
			// 连接建立好之后，代理服务器会面向客户端发送或接收 TCP 消息流。

			// 这里就要做重定向处理， 默认是 301
			code := http.StatusMovedPermanently
			// 如果请求的方式不是 GET 就将 http 的响应码设置成 308
			if req.Method != http.MethodGet {
				code = http.StatusPermanentRedirect
			}

			// tsr 返回值是一个 bool 值，用来判断是否需要重定向, getValue 返回来的
			// RedirectTrailingSlash 这个就是初始化时候定义的，只有为 true 才会处理
			if tsr && r.RedirectTrailingSlash {
				if len(path) > 1 && path[len(path)-1] == '/' {
					req.URL.Path = path[:len(path)-1]
				} else {
					req.URL.Path = path + "/"
				}
				// 执行重定向
				http.Redirect(w, req, req.URL.String(), code)
				return
			}
			// 路由没有找到，重定向规则也不符合，这里会尝试修复路径
			// 需要在初始化的时候定义 RedirectFixedPath 为 true，允许修复
			if r.RedirectFixedPath {
				// 这里就是在处理 Router 里面说的，将路径通过 CleanPath 方法去除多余的部分
				// 并且 RedirectTrailingSlash 为 ture 的时候，去匹配路由
				// 比如： 定义了一个路由 /foo , 但实际访问的是 ////FOO ，就会被重定向到 /foo
				fixedPath, found := root.findCaseInsensitivePath(
					CleanPath(path),
					r.RedirectTrailingSlash,
				)

				// 修复好的路径有处理路由的话 执行重定向
				if found {
					req.URL.Path = fixedPath
					http.Redirect(w, req, req.URL.String(), code)
				}
			}
		}
	}

	if req.Method == http.MethodOptions && r.HandleOPTHONS {
		// 处理 OPTHIONS 请求

	}

	// 处理 404
	if r.NotFound != nil {
		r.NotFound.ServeHTTP(w, req)
	} else {
		http.NotFound(w, req)
	}

}

func (r *Router) recv(w http.ResponseWriter, req *http.Request) {
	if rcv := recover(); rcv != nil {
		r.PanicHandler(w, req, rcv)
	}
}

// Params 是一个 []Param struct{key,value}
func (r *Router) getParams() *Params {
	ps, _ := r.paramsPool.Get().(*Params)
	*ps = (*ps)[0:0]
	return ps

}

func (r *Router) putParams(ps *Params) {
	if ps != nil {
		r.paramsPool.Put(ps)
	}
}
