load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

proto_library(
    name = "gentest_proto",
    srcs = [
        "only_streaming.proto",
        "only_unary_exec.proto",
        "only_unary_query.proto",
        "unary_type_test.proto",
    ],
    visibility = ["//visibility:public"],
)

go_proto_library(
    name = "gentest_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_grpc"],
    importpath = "github.com/jackskj/protoc-gen-map/testdata/gentest",
    proto = ":gentest_proto",
    visibility = ["//visibility:public"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "only_streaming.pb.map.go",
        "only_unary_exec.pb.map.go",
        "only_unary_query.pb.map.go",
        "unary_type_test.pb.map.go",
    ],
    embed = [":gentest_go_proto"],
    importpath = "github.com/jackskj/protoc-gen-map/testdata/gentest",
    visibility = ["//visibility:public"],
    deps = [
        "//mapper:go_default_library",
        "//templates:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
        "@com_github_masterminds_sprig//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)
