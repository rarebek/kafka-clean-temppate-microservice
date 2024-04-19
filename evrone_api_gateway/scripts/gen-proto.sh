#!/bin/bash

ROOT_DIR="."

PROTO_DIRS=($(find "$ROOT_DIR" -type f -name "*.proto" -exec dirname {} \; | sort -u))

for PROTO_DIR in "${PROTO_DIRS[@]}"; do
    GENPROTO_DIR="genproto/$(basename "$PROTO_DIR")"

    mkdir -p "$GENPROTO_DIR"

    echo "Generating protobuf and gRPC code for $PROTO_DIR..."


    protoc \
        --go_out="$GENPROTO_DIR" \
        --go-grpc_out="$GENPROTO_DIR" \
        "$PROTO_DIR"/*.proto

    echo "Code generated for $PROTO_DIR in $GENPROTO_DIR"
done

echo "All protos generated successfully."
