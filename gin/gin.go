package main

import (
	"html/template"
	"sync"

	"github.com/gin-gonic/gin/render"
)

type (
	RoutesInfo []RouteInfo
	RouteInfo  struct {
		Method  string
		Path    string
		Handler string
	}

	Engine struct {
		RouteInfoGroup
		delims      render.Delims
		FuncMap     template.FuncMap
		allNoRoute  HandlersChain
		allNoMethod HandlersChain
		noRoute     HandlersChain
		noMethod    HandlersChain
		pool        sync.Pool
		trees       methodTrees

		RedirectTrailingSlash bool
		RedirectFixedPath     bool
	}
)
