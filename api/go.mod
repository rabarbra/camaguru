module api

go 1.22.5

replace jwt => ./jwt

replace orm => ./orm

require (
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/crypto v0.25.0 // indirect
	jwt v0.0.0-00010101000000-000000000000 // indirect
	orm v0.0.0-00010101000000-000000000000 // indirect
)
