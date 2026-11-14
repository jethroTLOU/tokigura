package server

import (
	"net/http"

	"go.uber.org/fx"
)

type Route interface {
	Register(mux *http.ServeMux)
}

type RegisterRoutesParams struct {
	fx.In
	Mux    *http.ServeMux
	Routes []Route `group:"routes"`
}

func RegisterRoutes(p RegisterRoutesParams) {
	apiMux := http.NewServeMux()

	for _, r := range p.Routes {
		r.Register(apiMux)
	}

	p.Mux.Handle("/api/", http.StripPrefix("/api", apiMux))
}
