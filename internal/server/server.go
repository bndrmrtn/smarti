package server

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bndrmrtn/smarti/internal/ast"
	"github.com/bndrmrtn/smarti/internal/lexer"
	"github.com/bndrmrtn/smarti/internal/packages"
	"github.com/bndrmrtn/smarti/internal/runtime"
	"github.com/fatih/color"
)

type Server struct {
	dir string

	booster map[string][]ast.Node
}

func New(directory string) (*Server, error) {
	stat, err := os.Stat(directory)
	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		return nil, os.ErrNotExist
	}

	return &Server{
		dir:     directory,
		booster: make(map[string][]ast.Node),
	}, nil
}

func (s *Server) Start(listenAddr string) error {
	color.NoColor = true
	return http.ListenAndServe(listenAddr, s)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(s.dir, r.URL.Path)
	if r.URL.Path == "" || strings.HasSuffix(r.URL.Path, "/") {
		path = filepath.Join(path, "index.smt")
	}

	if !strings.HasSuffix(path, ".smt") {
		http.ServeFile(w, r, path)
		return
	}

	lx := lexer.New(path)
	if err := lx.Parse(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if nodes, ok := s.booster[lx.Sum()]; ok {
		s.execute(path, nodes, w, r)
		return
	}

	parser := ast.NewParser(lx.Tokens)
	if err := parser.Parse(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.booster[lx.Sum()] = parser.Nodes

	s.execute(path, parser.Nodes, w, r)
}

func (s *Server) execute(file string, nodes []ast.Node, w http.ResponseWriter, r *http.Request) {
	runt := runtime.New()

	runt.With("response", packages.NewResponse(w))
	runt.With("request", packages.NewRequest(r))

	if err := runt.Run(file, nodes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}
