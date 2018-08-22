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
