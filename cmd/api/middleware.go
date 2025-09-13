package main

import (
  "fmt"
  "net/http"
)

func (a *application)recoverPanic(next http.Handler)http.Handler  {
   return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
   // defer will be called when the stack unwinds
       defer func() {
           // recover() checks for panics
           err := recover();
           if err != nil {
               w.Header().Set("Connection", "close")
               a.serverErrorResponse(w, r, fmt.Errorf("%s", err))
           }
       }()
       next.ServeHTTP(w,r)
   })  
}


// SCENARIO #1: Let's assume our API is a public API. Anyone can use it once they
// signup. So for the clients' browsers to be able to read our API responses, we
// need to set the Access-Control-Allow-Origin header to everyone 
// (using the * operator). Notice we are back to returning 'http.Handler'
// Also notice that this middleware sets the header on the response object (w). We 
// set the response header early in the middleware chain to enable our
// response to be accepted by the client's browser

func (a *application) enableCORS (next http.Handler) http.Handler {                             
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
         w.Header().Add("Vary", "Origin")
        // Let's check the request origin to see if it's in the trusted list
        origin := r.Header.Get("Origin")
        // Once we have a origin from the request header we need need to check
 if origin != "" {
    for i := range a.config.cors.trustedOrigins {
        if origin == a.config.cors.trustedOrigins[i] {
           w.Header().Set("Access-Control-Allow-Origin", origin)                                          
           break
        }
   }
}
 
         next.ServeHTTP(w, r)
     })
 }
 
 //NEED GET THESE DONE!!
 func (a *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: check Authorization header / session here
		next.ServeHTTP(w, r)
	})
}

func (a *application) rateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: token bucket / IP-based limiter
		next.ServeHTTP(w, r)
	})
}
