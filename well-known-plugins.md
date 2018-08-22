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
