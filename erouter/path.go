package erouter

// 处理以下的各种情况
//	// Already clean
//	{"/", "/"},
//	{"/abc", "/abc"},
//	{"/a/b/c", "/a/b/c"},
//	{"/abc/", "/abc/"},
//	{"/a/b/c/", "/a/b/c/"},
//
//	// missing root
//	{"", "/"},
//	{"a/", "/a/"},
//	{"abc", "/abc"},
//	{"abc/def", "/abc/def"},
//	{"a/b/c", "/a/b/c"},
//
//	// Remove doubled slash
//	{"//", "/"},
//	{"/abc//", "/abc/"},
//	{"/abc/def//", "/abc/def/"},
//	{"/a/b/c//", "/a/b/c/"},
//	{"/abc//def//ghi", "/abc/def/ghi"},
//	{"//abc", "/abc"},
//	{"///abc", "/abc"},
//	{"//abc//", "/abc/"},
//
//	// Remove . elements
//	{".", "/"},
//	{"./", "/"},
//	{"/abc/./def", "/abc/def"},
//	{"/./abc/def", "/abc/def"},
//	{"/abc/.", "/abc/"},
//
//	// Remove .. elements
//	{"..", "/"},
//	{"../", "/"},
//	{"../../", "/"},
//	{"../..", "/"},
//	{"../../abc", "/abc"},
//	{"/abc/def/ghi/../jkl", "/abc/def/jkl"},
//	{"/abc/def/../ghi/../jkl", "/abc/jkl"},
//	{"/abc/def/..", "/abc"},
//	{"/abc/def/../..", "/"},
//	{"/abc/def/../../..", "/"},
//	{"/abc/def/../../..", "/"},
//	{"/abc/def/../../../ghi/jkl/../../../mno", "/mno"},
//
//	// Combinations
//	{"abc/./../def", "/def"},
//	{"abc//./../def", "/def"},
//	{"abc/../../././../def", "/def"},

func CleanPath(p string) string {
	const stackBufSize = 128

	//log.Println("======", p, "======")
	if p == "" {
		return "/"
	}

	buf := make([]byte, 0, stackBufSize)
	n := len(p)
	r := 1           // read指针
	w := 1           // write指针
	if p[0] != '/' { // 处理不是以 / 开头的情况
		r = 0
		if n+1 > stackBufSize {
			buf = make([]byte, n+1)
		} else {
			buf = buf[:n+1]
		}
		buf[0] = '/' // '/'00000000...
	}

	trailing := n > 1 && p[n-1] == '/'

	for r < n {
		switch {
		case p[r] == '/':
			r++
		case p[r] == '.' && r+1 == n: // /a/. end
			trailing = true
			r++
		case p[r] == '.' && p[r+1] == '/': // /a/./? || /a/./ end
			r += 2 // 跳过.
		case p[r] == '.' && p[r+1] == '.' && (r+2 == n || p[r+2] == '/'): // /a/.. end || /a/../
			r += 3 // 跳过..
			// 这里w--是为了一会儿复制非 .. 元素的时候跳过已经准备写入的路径 /a/b/.. 这时候会越过b不写入
			if w > 1 {
				w--
				if len(buf) == 0 {
					// 缓存为空 调整源数据的写入位置
					for w > 1 && p[w] != '/' {
						//log.Println("len(buf) == 0", p[:w+1])
						w--
					}
				} else {
					// 缓存不为空 调整缓存的写入位置
					for w > 1 && buf[w] != '/' {
						//log.Println("len(buf) != 0", string(buf[:w+1]))
						w--
					}
				}
			}
		default:
			// 实际的路径元素
			if w > 1 {
				bufApp(&buf, p, w, '/') // 去除多个连续的 /
				w++
			}

			// 复制元素
			for r < n && p[r] != '/' { // 直到碰到下一个 / 为止进行元素复制
				bufApp(&buf, p, w, p[r])
				w++
				r++
			}
		}
	}

	if trailing && w > 1 {
		bufApp(&buf, p, w, '/')
		w++
	}

	if len(buf) == 0 {
		return p[:w]
	}

	//log.Println("buf[:w]", string(buf[:w]))
	return string(buf[:w])
}

func bufApp(buf *[]byte, s string, w int, c byte) {
	b := *buf
	//log.Println("len(buf)", len(b), string(b))
	if len(b) == 0 {
		if s[w] == c {
			//log.Println("检测分支", s[:w+1])
			return
		}
		// 处理 r 与 w 不同步的情况
		if l := len(s); l > cap(b) {
			*buf = make([]byte, len(s))
		} else {
			*buf = (*buf)[:l]
		}
		b = *buf
		copy(b, s[:w])
		//log.Println("复制分支", string(b), w)
	}
	b[w] = c
}
