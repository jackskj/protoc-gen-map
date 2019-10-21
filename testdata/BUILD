load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "testdata.go",
        "tests.pb.map.go",
    ],
    embed = [":tests_go_proto"],
    importpath = "github.com/jackskj/protoc-gen-map/testdata",
    visibility = ["//visibility:public"],
    deps = [
        "//examples:go_default_library",
        "//mapper:go_default_library",
        "//templates:go_default_library",
        "//testdata/initdb:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
        "@com_github_golang_protobuf//ptypes:go_default_library_gen",
        "@com_github_lib_pq//:go_default_library",
        "@com_github_masterminds_sprig//:go_default_library",
        "@io_bazel_rules_go//proto/wkt:timestamp_go_proto",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

proto_library(
    name = "tests_proto",
    srcs = ["tests.proto"],
    visibility = ["//visibility:public"],
)

go_proto_library(
    name = "tests_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_grpc"],
    importpath = "github.com/jackskj/protoc-gen-map/testdata",
    proto = ":testdata_proto",
    visibility = ["//visibility:public"],
    deps = ["//examples:go_default_library"],
)

proto_library(
    name = "testdata_proto",
    srcs = ["tests.proto"],
    visibility = ["//visibility:public"],
    deps = [
        "//examples:examples_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)
