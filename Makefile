GATEWAY_PROTO_OUT_DIR = pkg/gateway-api
GATEWAY_API_PATH = api/gateway

gen-gateway:
		protoc \
				-I ${GATEWAY_API_PATH} \
				-I third_party/googleapis \
				-I third_party/envoyproxy/protoc-gen-validate \
				-I C:/ProgramData/protoc-24.4-win64/include \
				-I${GOPATH}/bin \
				--go_out=./$(GATEWAY_PROTO_OUT_DIR) --go_opt=paths=source_relative \
                --go-grpc_out=./$(GATEWAY_PROTO_OUT_DIR)  --go-grpc_opt=paths=source_relative \
                --validate_out="lang=go:./$(GATEWAY_PROTO_OUT_DIR)" --validate_opt=paths=source_relative \
                --grpc-gateway_out=./$(GATEWAY_PROTO_OUT_DIR) --grpc-gateway_opt=paths=source_relative \
                --openapiv2_out=use_go_templates=true,json_names_for_fields=false,allow_merge=true,merge_file_name=api:./$(GATEWAY_PROTO_OUT_DIR) \
                ./${GATEWAY_API_PATH}/*.proto


