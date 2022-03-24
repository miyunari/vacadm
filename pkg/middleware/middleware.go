package middleware

import (
	"net/http"

	jwt "github.com/MninaTB/vacadm/pkg/jwt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func Auth(t *jwt.Tokenizer) mux.MiddlewareFunc {
	logger := logrus.WithField("component", "auth-middleware")
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := jwt.ExtractToken(r)
			if err != nil {
				logger.Error(err)
				return
			}
			err = t.Valid(token)
			if err != nil {
				logger.Error(err)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

func Logging() mux.MiddlewareFunc {
	logger := logrus.WithField("component", "log-middleware")
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger = logger.WithField("path", r.URL.String())
			logger.Info("new request")
			h.ServeHTTP(w, r)
			logger.Info("end request")
		})
	}
}
