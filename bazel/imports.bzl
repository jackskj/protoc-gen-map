load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

def dept_imports():
    protobuf_deps()
    go_rules_dependencies()
    go_register_toolchains()
    gazelle_dependencies()
