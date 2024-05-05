package router

import (
	"strings"

	"github.com/WilhelmWeber/toy_web_server/src/final"
	"github.com/WilhelmWeber/toy_web_server/src/toy"
)

type Router struct {
	Route              string
	children           map[string]*Router
	params_child_route string
	full_route         []string
	controllers        map[string]func(c *toy.ToyContext)
}

func NewRouter() *Router {
	return &Router{
		Route:              "",
		children:           make(map[string]*Router),
		params_child_route: "",
	}
}

// 再帰的にルーティング木を検索
func (r *Router) SearchRouter(
	curr_index int,
	routes []string,
) ([]string, *Router) {
	//現在のインデックスがルートパラメータ配列の最後を指すのであれば
	if curr_index == len(routes)-1 {
		return r.full_route, r
	} else {
		child, isExist := r.children[routes[curr_index+1]]
		if isExist {
			return child.SearchRouter(curr_index+1, routes)
		} else {
			return routes, nil
		}
	}
}

func (r *Router) GetControllerAndInfo(
	method_type string,
	routes []string,
) ([]string, func(c *toy.ToyContext)) {
	full_routes, router := r.SearchRouter(0, routes)
	if router != nil {
		controller := router.GetController(method_type)
		return full_routes, controller
	} else {
		return full_routes, nil
	}
}

func (r *Router) GetController(method_type string) func(c *toy.ToyContext) {
	return r.controllers[method_type]
}

// 再帰的にルーティング木を構築する
func (r *Router) CreateRouting(
	curr_index int,
	routes []string,
) *Router {
	if curr_index == len(routes)-1 {
		return r
	} else {
		obj_path := routes[curr_index+1]
		//パラメータの登録
		if strings.HasPrefix(obj_path, ":") {
			if r.params_child_route != "" && obj_path != r.params_child_route {
				panic("Error: params" + obj_path + " is conflicted with params" + r.params_child_route)
			}
			r.params_child_route = obj_path
		}
		if _, isExist := r.children[obj_path]; !isExist {
			r.children[obj_path] = &Router{
				Route:      obj_path,
				children:   make(map[string]*Router),
				full_route: routes[:curr_index+1],
			}
		}
		next := r.children[routes[curr_index+1]]
		return next.CreateRouting(curr_index+1, routes)
	}
}

func (r *Router) AddController(
	method_type final.Method,
	controller func(c *toy.ToyContext),
) {
	method_str := method_type.ToString()
	if _, isExist := r.controllers[method_str]; isExist {
		panic("controller of " + method_str + " already exists")
	}
	r.controllers[method_str] = controller
}
