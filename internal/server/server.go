package server

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"

	"subs-server/internal/config"
	"subs-server/internal/provider"
	"subs-server/internal/utils"
)

type Server struct {
	provider provider.Provider
	debug    bool
	config   *config.Config
}

func NewServer(provider provider.Provider, debug bool, cfg *config.Config) *Server {
	return &Server{
		provider: provider,
		debug:    debug,
		config:   cfg,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	endpoint := strings.TrimPrefix(r.URL.Path, "/")

	if endpoint == "" && s.debug {
		s.listEndpoints(w)
		return
	}

	if endpoint == "" {
		http.NotFound(w, r)
		return
	}

	content, exists := s.provider.GetFile(endpoint)
	if !exists {
		http.NotFound(w, r)
		return
	}

	processed, err := utils.ProcessContent(content)
	if err != nil {
		log.Printf("Error processing content for endpoint %s: %v", endpoint, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	w.Header().Set("profile-title", "base64:"+base64.StdEncoding.EncodeToString([]byte(s.config.ProfileTitle)))
	w.Header().Set("profile-update-interval", s.config.ProfileUpdateInterval)
	w.Header().Set("profile-web-page-url", s.config.ProfileWebPageURL)
	w.Header().Set("support-url", s.config.SupportURL)

	_, _ = w.Write(processed)
}

func (s *Server) listEndpoints(w http.ResponseWriter) {
	endpoints := s.provider.ListEndpoints()
	sort.Strings(endpoints)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, "<h1>Available endpoints:</h1>")
	fmt.Fprintln(w, "<ul>")
	for _, endpoint := range endpoints {
		fmt.Fprintf(w, `<li><a href="/%s">/%s</a></li>`, endpoint, endpoint)
	}
	fmt.Fprintln(w, "</ul>")

	fmt.Fprintln(w, "<h2>Configured Response Headers:</h2>")
	fmt.Fprintln(w, "<ul>")
	fmt.Fprintf(w, "<li>profile-title: base64:%s (%s)</li>", base64.StdEncoding.EncodeToString([]byte(s.config.ProfileTitle)), s.config.ProfileTitle)
	fmt.Fprintf(w, "<li>profile-update-interval: %s</li>", s.config.ProfileUpdateInterval)
	fmt.Fprintf(w, "<li>profile-web-page-url: <a href=\"%s\">%s</a></li>", s.config.ProfileWebPageURL, s.config.ProfileWebPageURL)
	fmt.Fprintf(w, "<li>support-url: <a href=\"%s\">%s</a></li>", s.config.SupportURL, s.config.SupportURL)
	fmt.Fprintln(w, "</ul>")
}
