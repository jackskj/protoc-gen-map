load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "generate.go",
        "imports.go",
        "plugin.go",
        "templates.go",
    ],
    importpath = "github.com/jackskj/protoc-gen-map/plugin",
    visibility = ["//visibility:public"],
    deps = [
        "//templates:go_default_library",
        "@com_github_gogo_protobuf//protoc-gen-gogo/descriptor:go_default_library",
        "@com_github_gogo_protobuf//protoc-gen-gogo/generator:go_default_library",
        "@com_github_masterminds_sprig//:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["generate_test.go"],
    embed = [":go_default_library"],
    deps = [
        "//examples:go_default_library",
        "//testdata:go_default_library",
        "//testdata/gentest:go_default_library",
        "//testdata/initdb:go_default_library",
    ],
)
