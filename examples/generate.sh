# Sample command to generate .pb.go and pb.map.go files
protoc --map_out="sql=./sql:."   \
        --go_out="plugins=grpc:."   \
        -I=. \
        ./query.proto
