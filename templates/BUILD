load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "enum_value_maps.go",
        "funcs.go",
        "init.go",
        "server.go",
        "streaing_response.go",
        "templates.go",
        "unary_response.go",
    ],
    importpath = "github.com/jackskj/protoc-gen-map/templates",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_golang_protobuf//ptypes:go_default_library_gen",
        "@io_bazel_rules_go//proto/wkt:timestamp_go_proto",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["funcs_test.go"],
    embed = [":go_default_library"],
    deps = [
        "@com_github_golang_protobuf//ptypes:go_default_library_gen",
        "@com_github_masterminds_sprig//:go_default_library",
    ],
)
