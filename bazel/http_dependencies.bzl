load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

dependencies = {
    # com_google_protobuf
    "rules_cc": {
        "sha256": "bb8320b0bc1d8d01dc8c8e8c50edced8553655c03776960c1287d03dfbcac3e5",
        "strip_prefix": "rules_cc-401380cd2279b83da0dcb86ecbac04a04805405b",
        "urls": [
            "https://mirror.bazel.build/github.com/bazelbuild/rules_cc/archive/401380cd2279b83da0dcb86ecbac04a04805405b.tar.gz",
            "https://github.com/bazelbuild/rules_cc/archive/401380cd2279b83da0dcb86ecbac04a04805405b.tar.gz",
        ],
    },
    # com_google_protobuf
    "rules_java": {
        "sha256": "4e2f33528a66e3a9909910eaa5a562fb22f5b422513cdc3816fd01fbb6e2d08d",
        "strip_prefix": "rules_java-166a046a27e118d578127759b413ee0b06aca3cd",
        "urls": [
            "https://mirror.bazel.build/github.com/bazelbuild/rules_java/archive/166a046a27e118d578127759b413ee0b06aca3cd.tar.gz",
            "https://github.com/bazelbuild/rules_java/archive/166a046a27e118d578127759b413ee0b06aca3cd.tar.gz",
        ],
    },
    # skylib_version 0.9.0 introduced a breaking change to rules_proto
    "bazel_skylib": {
        "sha256": "2ef429f5d7ce7111263289644d233707dba35e39696377ebab8b0bc701f7818e",
        "urls": [
            "https://github.com/bazelbuild/bazel-skylib/releases/download/0.8.0/bazel-skylib.0.8.0.tar.gz",
        ],
    },
    "io_bazel_rules_go": {
        "urls": [
            "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/rules_go/releases/download/0.19.1/rules_go-0.19.1.tar.gz",
            "https://github.com/bazelbuild/rules_go/releases/download/0.19.1/rules_go-0.19.1.tar.gz",
        ],
        "sha256": "8df59f11fb697743cbb3f26cfb8750395f30471e9eabde0d174c3aebc7a1cd39",
    },
    "io_bazel_rules_docker": {
        "sha256": "e513c0ac6534810eb7a14bf025a0f159726753f97f74ab7863c650d26e01d677",
        "strip_prefix": "rules_docker-0.9.0",
        "urls": ["https://github.com/bazelbuild/rules_docker/archive/v0.9.0.tar.gz"],
    },
    "com_google_protobuf": {
        "sha256": "c90d9e13564c0af85fd2912545ee47b57deded6e5a97de80395b6d2d9be64854",
        "strip_prefix": "protobuf-3.9.1",
        "urls": ["https://github.com/google/protobuf/archive/v3.9.1.zip"],
    },
    "rules_proto": {
        "sha256": "602e7161d9195e50246177e7c55b2f39950a9cf7366f74ed5f22fd45750cd208",
        "strip_prefix": "rules_proto-97d8af4dc474595af3900dd85cb3a29ad28cc313",
        "urls": [
            "https://mirror.bazel.build/github.com/bazelbuild/rules_proto/archive/97d8af4dc474595af3900dd85cb3a29ad28cc313.tar.gz",
            "https://github.com/bazelbuild/rules_proto/archive/97d8af4dc474595af3900dd85cb3a29ad28cc313.tar.gz",
        ],
    },
    "bazel_gazelle": {
        "urls": [
                "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/bazel-gazelle/releases/download/v0.19.0/bazel-gazelle-v0.19.0.tar.gz",
                "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.19.0/bazel-gazelle-v0.19.0.tar.gz",
        ],
        "sha256": "41bff2a0b32b02f20c227d234aa25ef3783998e5453f7eade929704dcff7cd4b",
        }
}

def http_depts():
    for name in dependencies:
        if name in native.existing_rules():
            continue

        http_archive(
            name = name,
            **dependencies[name]
        )
