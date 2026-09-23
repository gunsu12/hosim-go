package migrations

import "embed"

// FS menyematkan seluruh file .sql di folder migrations ke dalam binary Go
//
//go:embed *.sql
var FS embed.FS
