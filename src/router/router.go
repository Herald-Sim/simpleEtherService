package router

import (
	"fmt"
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

type Router struct {
	Handlers map[string]map[string]HandlerFunc
}

type Context struct {
	Params   map[string]interface{}
	TempData interface{}

	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

func (r *Router) HandleFunc(method, pattern string, h HandlerFunc) {
	m, ok := r.Handlers[method]
	if !ok {
		m = make(map[string]HandlerFunc)
		r.Handlers[method] = m
	}

	m[pattern] = h
}

func match(pattern, path string) (bool, map[string]string) {
	fmt.Println("Call match")

	if pattern == path {
		return true, nil
	}

	patterns := strings.Split(pattern, "/")
	paths := strings.Split(path, "/")

	if len(patterns) != len(paths) {
		return false, nil
	}

	params := make(map[string]string)

	for i := 0; i < len(patterns); i++ {
		switch {
		case patterns[i] == paths[i]:
		case len(patterns[i]) > 0 && patterns[i][0] == ':':
			params[patterns[i][1:]] = paths[i]

		default:
			return false, nil
		}
	}

	return true, params
}

// ServeHTTP : 웹 요청의 http method와 URL path를 분석하여 그에 맞는 핸들러를 찾아 동작,
// 만약 해당 요청에 일치하는 핸들러가 존재하지 않을경우, NotFound 에러를 반환
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Call ServeHTTP")

	for pattern, handler := range r.Handlers[req.Method] {
		if ok, parameter := match(pattern, req.URL.Path); ok {
			fmt.Println(parameter)
			c := Context{
				Params:         make(map[string]interface{}),
				ResponseWriter: w,
				Request:        req,
			}

			for k, v := range parameter {
				c.Params[k] = v
			}

			handler(&c)
			return
		}
	}
	http.NotFound(w, req)
	return
}