# Well-Known Plugins

This document discusses the [Well-Known Plugins Issue](https://github.com/uber/prototool/issues/2), a proposal of which is implemented on this branch.

## Proposal Result

Here are the changes to [example/idl/uber/prototool.yaml](example/idl/uber/prototool.yaml) this results in:

From:

```yaml
protoc:
  includes:
    - ../../../vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis
generate:
  go_options:
    import_path: github.com/uber/prototool/example/idl/uber
    extra_modifiers:
      google/api/annotations.proto: google.golang.org/genproto/googleapis/api/annotations
      google/api/http.proto: google.golang.org/genproto/googleapis/api/annotations
  plugins:
    - name: gogoslick
      type: gogo
      flags: plugins=grpc
      output: ../../gen/proto/go
    - name: yarpc-go
      type: gogo
      output: ../../gen/proto/go
    - name: java
      output: ../../gen/proto/java
    - name: grpc-gateway
      type: go
      output: ../../gen/proto/go
```

To:

```yaml
protoc:
  include_extended_wkt: true
generate:
  output: ../../gen/proto
  go_options:
    import_path: github.com/uber/prototool/example/idl/uber
  plugins:
    - alias: gogoslick-with-grpc
    - alias: yarpc-go
    - alias: java
    - alias: grpc-gateway
```

Here is the output of the new flag `--list-well-known-plugins` on the sub-command `prototool generate`:

```bash
$ prototool generate --list-well-known-plugins
cpp:
  name:    cpp
  output:  cpp
cpp-grpc:
  name:    cpp-grpc
  path:    grpc_cpp_plugin
  output:  cpp
csharp:
  name:    csharp
  output:  csharp
csharp-grpc:
  name:    csharp-grpc
  path:    grpc_csharp_plugin
  output:  csharp
go:
  name:    go
  type:    go
  output:  go
go-with-grpc:
  name:    go
  type:    go
  flags:   plugins=grpc
  output:  go
gofast:
  name:    gofast
  type:    gogo
  output:  go
gofast-with-grpc:
  name:    gofast
  type:    gogo
  flags:   plugins=grpc
  output:  go
gogo:
  name:    gogo
  type:    gogo
  output:  go
gogo-with-grpc:
  name:    gogo
  type:    gogo
  flags:   plugins=grpc
  output:  go
gogofast:
  name:    gogofast
  type:    gogo
  output:  go
gogofast-with-grpc:
  name:    gogofast
  type:    gogo
  flags:   plugins=grpc
  output:  go
gogofaster:
  name:    gogofaster
  type:    gogo
  output:  go
gogofaster-with-grpc:
  name:    gogofaster
  type:    gogo
  flags:   plugins=grpc
  output:  go
gogoslick:
  name:    gogoslick
  type:    gogo
  output:  go
gogoslick-with-grpc:
  name:    gogoslick
  type:    gogo
  flags:   plugins=grpc
  output:  go
grpc-gateway:
  name:    grpc-gateway
  type:    go
  output:  go
java:
  name:    java
  output:  java
js:
  name:    js
  output:  js
objc:
  name:    objc
  output:  objc
objc-grpc:
  name:    objc-grpc
  path:    grpc_objective_c_plugin
  output:  objc
php:
  name:    php
  output:  php
php-grpc:
  name:    php-grpc
  path:    grpc_php_plugin
  output:  php
python:
  name:    python
  output:  python
python-grpc:
  name:    python-grpc
  path:    grpc_python_plugin
  output:  python
ruby:
  name:    ruby
  output:  ruby
ruby-grpc:
  name:    ruby-grpc
  path:    grpc_ruby_plugin
  output:  ruby
yarpc-go:
  name:    yarpc-go
  type:    gogo
  output:  go
```

## Overview

Prototool is meant to take away the complexity of working with `protoc`. The `generate` sub-command is a big part of this, especially for
Golang - the additional steps required to get Golang generation have the correct imports are not trivial. However, `generate` feels
incomplete, as we disuss in the Well-Known Plugins Issue - users still need to do a lot of setup to get `generate` to work for common
plugins. Prototool takes a "dependency injection" approach in that it purposefully does not handle plugin management, rather saying
any plugin will work, and it is your responsiblity to set the plugins up. This approach was taken as opposed to my previous approach
prior to Uber in [Protoeasy](https://github.com/peter-edge/protoeasy-go), which in effect was responsible for "knowing the world", however
the flag `--extra-plugin` was added later. The "knowing the world" approach proved to not be scalable, and resulted in out-of-date plugins
and assumed configuration that was not easy to override.

This proposal takes an intermediate approach. It still leaves the user to install their own plugins - this is out of the scope of
Prototool's current functionality, and in my opinion, should stay that way. However, it declares Well-Known Plugins - those plugins
which are generally widely used in the Protobuf ecosystem - and effectively stores defaults for the values a user would fill in
in their `prototool.yaml` files for each plugin. It also adds the concept of Extended Well-Known Types, effectively just meaning
`google/api/annotations.proto` and `google/api/http.proto`, two files that while not part of the Well-Known Types, are very commonly
used and would be nice to manage with Prototool, especially if we consider `protoc-gen-grpc-gateway` a Well-Known Plugin.

The end result provides a much simpler configuration option for users that will likely also result in more consistency in the generated code
across all `prototool generate` usage.

## Commits

Probably the easiest way to discuss the changes here are on a per-commit basis. The commits are roughly split into separate logical changes.

### fa20182fa618ddb501922c35dd73dfca40aea93c Add extended Well-Known Types and write to Prototool cache

This commit adds the concept of Extended Well-Known Types. The files [google/api/annotations.proto](https://github.com/googleapis/googleapis/blob/master/google/api/annotations.proto)
and [google/api/http.proto](https://github.com/googleapis/googleapis/blob/master/google/api/http.proto) are part of the official Google
APIs, however, they are widely used and required for `protoc-gen-grpc-gateway`. In my opinion, we should not mix them with the Well-Known
Types - these are defined set that should not change - however these two files are so commonly used that it would be nice if Prototool
handled them. I don't expect the list of Extended Well-Known Types will go beyond these two files - if in the future there is a global
Protobuf standard libary beyond the Well-Known Types, perhaps this library should be added, but an Uber-specfic Protobuf standard library
should not, and we might even want a separate Standard Libary concept.

This commit takes the data of those two files and makes them string constants, substituting the long-form `go_package` values for short-form
`go_package` values. Note this substitution is not actually needed, but done for consistency (and substituting `go_package` values is a
somewhat common pattern). The data is then written to a new `extend/include` directory within the `protoc` cache directory, so it can be
separately included when Prototool calls `protoc`. Here's the resulting cache directory on my system:

```bash
[~/Library/Caches/prototool/Darwin/x86_64/protobuf/3.6.1]
$ tree
.
├── bin
│   └── protoc
├── extended
│   └── include
│       └── google
│           └── api
│               ├── annotations.proto
│               └── http.proto
├── include
│   └── google
│       └── protobuf
│           ├── any.proto
│           ├── api.proto
│           ├── compiler
│           │   └── plugin.proto
│           ├── descriptor.proto
│           ├── duration.proto
│           ├── empty.proto
│           ├── field_mask.proto
│           ├── source_context.proto
│           ├── struct.proto
│           ├── timestamp.proto
│           ├── type.proto
│           └── wrappers.proto
└── readme.txt
```

This commit also adds the github.com/golang/protobuf and github.com/gogo/protobuf modifier mappings as Golang maps. Note that
the modifier mappings are actually the same, but I kept them separate in case this were to change, although we couldn't change
this after release, so perhaps this separation is not needed.
