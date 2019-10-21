load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:build_file_name BUILD

gazelle(
    name = "gazelle_update",
    prefix = "github.com/jackskj/protoc-gen-map",
    command = "update",
    extra_args = [],
)

gazelle(
    name = "gazelle_resolve",
    prefix = "github.com/jackskj/protoc-gen-map",
    extra_args = [],
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "github.com/jackskj/protoc-gen-map",
    visibility = ["//visibility:private"],
    deps = [
        "//plugin:go_default_library",
        "@com_github_gogo_protobuf//vanity/command:go_default_library",
    ],
)

go_binary(
    name = "protoc-gen-map",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)
