package gorestframework

import (
	"log"
	"net/http"
	"strings"
)

// LoggingMiddleware is used to log every request with Method and RequestURI
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.Method, r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware is used to enable CORS capability for all hosts, PR are welcome!
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			return
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// CORS is a middleware which enables Access-Control-Allow-Methods
// for specified methods. It also handles request method, denying it
// with a 405 (method not allowed) if it is not among the provided
// methods.
// example:
// router.Handle("/path", CORS(
// 		http.MethodGet,
// 		http.MethodPatch,
// 		http.MethodPut,
// 		http.MethodDelete,
// 	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		switch r.Method {
// 		case http.MethodGet:
// 			srv.getDashboard(w, r)
// 		case http.MethodPatch, http.MethodPut:
// 			srv.editDashboard(w, r)
// 		case http.MethodDelete:
// 			srv.deleteDashboard(w, r)
// 		default:
// 			srv.RespondStatusError(w, r, http.StatusMethodNotAllowed)
// 		}
// 	}))
// )
func CORS(methods ...string) func(http.Handler) http.Handler {
	if !isInMethods(http.MethodOptions, methods) {
		methods = append(methods, http.MethodOptions)
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
			if r.Method == http.MethodOptions {
				return
			}
			if !isInMethods(r.Method, methods) {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// isInMethods is a utility function which checks if a method.
// For the sake of clarity, it is implemented in a linear search, but
// it can be improved probably.
func isInMethods(method string, methods []string) bool {
	for _, tmp := range methods {
		if tmp == method {
			return true
		}
	}
	return false
}

// AllowedOrigin is the middleware used to send Access-Control-Allow-Origin header
func AllowedOrigin(origin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			next.ServeHTTP(w, req)
		})
	}
}

// AllowedHeaders is the middleware used to send Access-Control-Allow-Headers header
func AllowedHeaders(headers ...string) func(http.Handler) http.Handler {
	if len(headers) == 0 {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				next.ServeHTTP(w, req)
			})
		}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
			next.ServeHTTP(w, req)
		})
	}
}
