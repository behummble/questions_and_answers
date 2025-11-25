package http

import (
	"net/http"
	"strings"
	"regexp"
	"context"
)

type Route struct {
    Method  string
    Pattern *regexp.Regexp
    Handler http.HandlerFunc
}

type Router struct {
    routes []Route
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    for _, route := range r.routes {
        if route.Method != req.Method {
            continue
        }
        
        matches := route.Pattern.FindStringSubmatch(req.URL.Path)
        if matches != nil {
            ctx := req.Context()
            for i, name := range route.Pattern.SubexpNames() {
                if i > 0 && name != "" {
                    ctx = context.WithValue(ctx, name, matches[i])
                }
            }
            route.Handler(w, req.WithContext(ctx))
            return
        }
    }
    
    http.NotFound(w, req)
}

func newRouter() *Router {
    return &Router{
		routes: make([]Route, 0, 7),
	}
}

func (r *Router) addRoute(method, pattern string, handler http.HandlerFunc) {
    regexPattern := "^" + pattern + "$"
    regexPattern = strings.ReplaceAll(regexPattern, "{id}", `([0-9]+)`)
    regexPattern = strings.ReplaceAll(regexPattern, "{slug}", `([^/]+)`)
    
    compiled := regexp.MustCompile(regexPattern)
    r.routes = append(r.routes, Route{
        Method:  method,
        Pattern: compiled,
        Handler: handler,
    })
}

func(r *Router) register(s *Server) {
    r.addRoute("POST", "/questions/{id}/answers/", s.createAnswer)
    r.addRoute("GET", "/answers/{id}", s.getAnswer)
    r.addRoute("DELETE", "/answers/{id}", s.deleteAnswer)
    
    r.addRoute("POST", "/questions", s.createQuestion)
    r.addRoute("GET", "/questions", s.getAllQuestions)
    r.addRoute("GET", "/questions/{id}", s.getQuestion)
    r.addRoute("DELETE", "/questions/{id}", s.deleteQuestion)
}

// Helper to get path parameter from context
func getPathParam(ctx context.Context, param string) string {
    if value, ok := ctx.Value(param).(string); ok {
        return value
    }
    return ""
}