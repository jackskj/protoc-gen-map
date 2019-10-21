load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

proto_library(
    name = "initdb_proto",
    srcs = ["initdb.proto"],
    visibility = ["//visibility:public"],
    deps = ["@com_google_protobuf//:timestamp_proto"],
)

go_proto_library(
    name = "initdb_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_grpc"],
    importpath = "github.com/jackskj/protoc-gen-map/testdata/initdb",
    proto = ":initdb_proto",
    visibility = ["//visibility:public"],
)

go_library(
    name = "go_default_library",
    srcs = ["initdb.pb.map.go"],
    embed = [":initdb_go_proto"],
    importpath = "github.com/jackskj/protoc-gen-map/testdata/initdb",
    visibility = ["//visibility:public"],
    deps = [
        "//mapper:go_default_library",
        "//templates:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
        "@com_github_masterminds_sprig//:go_default_library",
        "@io_bazel_rules_go//proto/wkt:timestamp_go_proto",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)
