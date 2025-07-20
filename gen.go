package gen

//go:generate sqlc -f db/sqlc.yaml generate
//go:generate templ generate
//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/user.proto ./proto/category.proto
