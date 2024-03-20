module 7days

replace Gee => ./Gee

replace Orm => ./Orm

replace Cache => ./Cache

go 1.22.1

require (
	Cache v0.0.0-00010101000000-000000000000
	Gee v0.0.0-00010101000000-000000000000
	Orm v0.0.0-00010101000000-000000000000
	github.com/mattn/go-sqlite3 v1.14.22
)
