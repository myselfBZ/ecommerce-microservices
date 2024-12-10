package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

type MiddlewareError struct{
    err string 
    status int
}

func NewMiddlewareError(err string, s int) *MiddlewareError{
    return &MiddlewareError{
        err: err,
        status: s,
    }
}

func (m *MiddlewareError) setError(msg string) {
    m.err = msg
}

func (m *MiddlewareError) setStatus(s int){
    m.status = s
}

type Middleware func(w http.ResponseWriter, r *http.Request) *MiddlewareError


func middlewareFunc(next http.Handler, middlwareFuncs...Middleware)  http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        defer func(){
            log.Printf("Took: %f", time.Since(start).Seconds())
        }()
        for _, m := range middlwareFuncs{
            // stop the chain if one of middleware functions return an error
            if resp := m(w, r); resp.err != ""{
                http.Error(w, resp.err, resp.status)
                return 
            }
        }

		next.ServeHTTP(w, r)
	})
}

func JWTValidate(w http.ResponseWriter, r *http.Request) *MiddlewareError {

    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        return NewMiddlewareError("unauthorized", http.StatusUnauthorized) 
    }
    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        return NewMiddlewareError("unauthorized", http.StatusUnauthorized) 
    }

    token := parts[1]
    claims, err := ValidateToken(token)
    if err != nil {
        return NewMiddlewareError("unauthorized", http.StatusUnauthorized) 
    }
    ctx := context.WithValue(r.Context(), "user_id", claims.UserID) 
    r.WithContext(ctx)
    return nil 
}
