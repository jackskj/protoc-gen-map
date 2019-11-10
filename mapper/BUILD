load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "mapper.go",
        "reflect.go",
    ],
    importpath = "github.com/jackskj/protoc-gen-map/mapper",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_golang_protobuf//ptypes:go_default_library_gen",
        "@io_bazel_rules_go//proto/wkt:timestamp_go_proto",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "mapper_test.go",
        "reflect_test.go",
    ],
    data = [
        "//testdata:mapper.golden",
    ],
    embed = [":go_default_library"],
    deps = [
        "//examples:go_default_library",
        "//testdata:go_default_library",
        "//testdata/initdb:go_default_library",
        "@com_github_golang_protobuf//jsonpb:go_default_library_gen",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@com_github_yudai_gojsondiff//:go_default_library",
        "@com_github_yudai_gojsondiff//formatter:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//test/bufconn:go_default_library",
    ],
)
