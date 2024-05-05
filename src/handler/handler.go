package handler

import (
	"fmt"
	"net"
	"strings"

	"github.com/WilhelmWeber/toy_web_server/src/final"
	"github.com/WilhelmWeber/toy_web_server/src/router"
	"github.com/WilhelmWeber/toy_web_server/src/toy"
)

type Handler struct {
	Router            *router.Router
	static_asset_path string
}

func NewHadndler() *Handler {
	return &Handler{
		Router:            router.NewRouter(),
		static_asset_path: "",
	}
}

func (h *Handler) AddRoute(
	method final.Method,
	route string,
	controller func(c *toy.ToyContext),
) {
	route_paths := strings.Split(route, "/")
	if route_paths[0] != "" {
		route_paths = append([]string{""}, route_paths...)
	}
	var router *router.Router
	if _, obj := h.Router.SearchRouter(0, route_paths); obj != nil {
		router = obj
	} else {
		router = h.Router.CreateRouting(0, route_paths)
	}
	router.AddController(method, controller)
}

func (h *Handler) AddStaticPath(static_path string) {
	h.static_asset_path = static_path
}

// リクエストのパース
func (h *Handler) HandleConnection(conn net.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		panic(err)
	}
	requests := string(buf[:n])
	var http_method string
	var route string
	var http_ver string
	req_header_parserd := false
	var req_body string
	context := toy.NewContext()

	for i, line := range strings.Split(requests, "\n") {
		if len(line) == 0 {
			req_header_parserd = true
			continue
		}
		if i == 0 {
			req_line := strings.SplitN(line, " ", 3)
			http_method = req_line[0]
			route_with_query := strings.Split(req_line[1], "?")
			route = route_with_query[0]
			//query parse
			if len(route_with_query) > 1 {
				queries := strings.Split("&", route_with_query[2])
				for _, query := range queries {
					key_and_value := strings.SplitN(query, "=", 2)
					if len(key_and_value) == 2 {
						context.SetQueryParams(key_and_value[0], key_and_value[1])
					}
				}
			}
			http_ver = req_line[2]
		} else if !req_header_parserd {
			header := strings.Split(line, ": ")
			key := header[0]
			value := header[1]
			context.SetReqHeader(key, value)
		} else {
			req_body += line
		}
	}

	req_body = strings.TrimSpace(req_body)
	content_type := context.GetReqHeader("Content-Type")

	if strings.Contains(content_type, final.APPLICATION_WWW_FORM_URLENCODED) {
		form_params := strings.Split(req_body, "&")
		for _, param := range form_params {
			key_and_value := strings.SplitN(param, "=", 2)
			if len(key_and_value) == 2 {
				context.SetFormParams(key_and_value[0], key_and_value[1])
			}
		}
	} else if strings.Contains(content_type, final.MULTIPART_FORMDATA) {

	}
	//TODO: JSONのパース
	context.SetReqBody(req_body)

	var response, res_status string
	if _, isExist := final.Methods[http_method]; isExist {
		route_paths := strings.Split(route, "/")
		if route_paths[0] != "" {
			route_paths = append([]string{""}, route_paths...)
		}
		if info_paths, controller := h.Router.GetControllerAndInfo(http_method, route_paths); controller != nil {
			if !context.HasReponse() {
				for i, info_path := range info_paths {
					if strings.HasPrefix(info_path, ":") {
						key := strings.Trim(info_path, ":")
						value := route_paths[i]
						context.SetUrlParams(key, value)
					}
				}
				controller(context)
			}
			response, res_status = context.GetResponse(http_ver)
		} else {
			response, res_status = context.Default404(http_ver)
		}
	} else {
		response, res_status = context.Default400(http_ver)
	}

	conn.Write([]byte(response))

	defer conn.Close()

	fmt.Printf("|%s| %s %s", http_method, route, res_status)
}
