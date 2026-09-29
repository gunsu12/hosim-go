package web

import (
	"embed"
	"io/fs"
)

// FS menyematkan seluruh isi folder dist (hasil build Svelte/Vite) ke dalam binary Go
//
//go:embed all:dist
var distFS embed.FS

// Dist mengembalikan fs.FS yang langsung mengarah ke dalam folder dist
func Dist() fs.FS {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic(err)
	}
	return sub
}
