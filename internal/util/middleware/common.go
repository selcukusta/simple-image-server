package middleware

import "net/http"

//CommonMiddleware is using to add common headers to the response
func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET")
		next.ServeHTTP(w, r)
	})
}
