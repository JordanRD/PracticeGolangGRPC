Generate go stubs
```powershell
protoc --go-grpc_out=. --go_out=.  bookshop/bookshop.proto
```

Run
```powershell
go run main.go
```

