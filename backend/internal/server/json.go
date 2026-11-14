package server

import (
	"encoding/json"
	"net"
	"net/http"
	"net/netip"

	"github.com/jethroTLOU/tokigura/internal/core/apperr"
)

const maxBodyBytes = 1 << 20 // 1 MB

type Map map[string]any

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(dst); err != nil {
		return apperr.NewBindError(err)
	}

	return nil
}

func ExtractIP(r *http.Request) (netip.Addr, error) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}

	ip, err := netip.ParseAddr(host)
	if err != nil {
		return netip.Addr{}, apperr.New(apperr.CodeBadRequest, "invalid IP address")
	}

	return ip.WithZone(""), nil
}
