### Title
Grpc and api gateway 

### refer

https://github.com/grpc-ecosystem/grpc-gateway#usage  

https://www.liwenzhou.com/posts/Go/grpc-gateway/


### build
```bash
protoc -I . --go_out=gen --go_opt=paths=source_relative --grpc-gateway_out gen --grpc-gateway_opt paths=source_relative  --grpc-gateway_opt generate_unbound_methods=true  --go-grpc_out=gen --go-grpc_opt=paths=source_relative  types/*.proto 
```
