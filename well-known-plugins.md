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

This commit also adds the `github.com/golang/protobuf` and `github.com/gogo/protobuf` modifier mappings as Golang maps. Note that
the modifier mappings are actually the same, but I kept them separate in case this were to change, although we couldn't change
this after release, so perhaps this separation is not needed.

### ca3d2d98ef87e9aaa51eb05901100caf8d790560 Add protoc.include_extended_wkt option and propagate to example

This commit adds the configuration option `protoc.include_extended_wkt`. If set to `true`, the Extended Well-Known Types will be includes
as part of `protoc` compilation, and modifiers will be added for plugins of type `go` and `gogo`. For reference, in the above
cache directory, the flag `-I /path/to/cache/dir/extended/include` will be added to `protoc` calls. This is similar to the now-deprecated
`include_wkt` configuration option, however we probably want to have users explicitly opt-in to to the Extended Well-Known Types, but
this should definitely be a discussion point.

This commit also updates the example to use this new option.

### 5f16686a658f6bda28014ece103d76c1f4214646 Allow relative plugin paths

This commit allows relative values to be set for the configuration option `gen.plugins.path`. This is something we probably want to change
about Prototool regardless of this feature - as of this proposal, you need to specify an absolute path such as
`/usr/local/bin/grpc_cpp_plugin` for the override to work, which is not independent of a given user's installation. This change
allows the above path to be specified as `grpc_cpp_plugin` instead, and Prototool will execute `which plugin_path` to determine
the absolute path on Darwin and Linux. Note that there is an explicit `switch` on `runtime.GOOS` to make sure this is being
executed on Darwin or Linux, but this same switch is in other places in Prototool. If we decide to add Windows compatibility,
this is another area that will need to be changed, however this is not the hardest problem in the world to solve on any platform.

### 33e9227fb99c26130e6f6ae8beecc36d764c0f29 Add Well-Known Plugin specifications

This commit adds the type `settings.WellKnownPlugin` and adds the values for the Well-Known Plugins. These are referenced by alias,
as opposed to our previous discussion of just overriding `gen.plugins.name`. The main idea here is that `name` is meant to reference
either a built-in "plugin" such as `cpp` or `java`, or an installed plugin such as `protoc-gen-go`, however there are situations
we want to handle, especially with gRPC, that do not map nicely to this concept. Going with the gRPc example, for many gRPC plugins,
what you really want is a dummy name along with an overriden `path` value. For Golang plugins, you want to add the flag `plugins=grpc`.
What this comes down to is that it may be confusing to the user to override the `name` concept, and instead it would be nice to say
"This is an alias for this set of default values," especially as `cpp, cpp-grpc` (most gRPC code generation requires two `*_out` calls)
is different from `go-with-grpc` (Golang gRPC code generation requires one `*_out` call). See the commentary in the code for more details.

This also adds the concept of relative output path to the `WellKnownPlugin`, something that will come up in a later commit. This is the
relative directory, relative to a base directory, that generated code will be placed in, if not overridden.

The following values of Well-Known Plugins are added:

- All the built-ins to `protoc` except `descriptor_set` (which we may want to add, especially for users of other tools that want this for
calling gRPC endpoints without the reflection API). See `protoc -h | grep _out` for the full list.
- gRPC values for all the built-ins for when you `brew install grpc`. If you `brew install grpc`, run `ls /usr/local/bin | grep grpc`
for the list. Note that `objc` changes to `grpc_objective_c_plugin`. Also note this does not include a value for Java, however this
should probably be added.
- Values for `github.com/golang/protobuf/protoc-gen-go` and the commonly-used plugins in `github.com/gogo/protobuf`.
- gRPC values for the Golang and Gogo values.
- Values for yarpc-go and grpc-gateway.

We might want to add `github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger` as well.

A nice addition might also be `URL` and `Description` fields to the `WellKnownPlugin` type that we can output for more details for the user.

Also note that grpc-gateway requires the Extended Well-Known Types, but we do not do any enforcement that this value is set. See the above
discussion on not opting-in to the Extended Well-Known Types for the general reasoning behind this.

### 09565135b4332a008a85fe541a96044de754df5c Add --list-well-known-plugins flag to generate

This adds `prototool generate --list-well-known-plugins` to print out the above specifications for the user. We might also want to add
things such as `prototool compile/generate --list-well-known-types` and `prototool compile/generate --list-extended-well-known-types`, for
example. The output happens to be in YAML format, but this was just icing on the cake, and was not the original intent of this flag.

### 98c5005661141726fbefb8e97e812375a067d1fd Wire everything up and add to cfginit

This commit is a bit overloaded and probably should have been split up for the purpose of this documentation.

This commit:

- Adds the `gen.output` configuration option. This value will be the base path for all plugins, and if a plugin output path is
relative, the two paths will be concatenated. Note there's no check as of now that a plugin output path is actually there, only that
at least one of `gen.output` and `gen.plugins.output` are set, we probably want to add this check.
- Adds the `gen.plugins.alias` configuration option. This references the Well-Known Plugin specifications. All other values of
`gen.plugins` will use the default values from the Well-Known Plugin specification, however all values can be overridden.
- Adds the logic to handle Well-Known Plugins in the `settings` package. Note that no other code needs to change to handle
Well-Known Plugins in Prototool - the `GenConfig` and `GenPlugin` structs remain the same (although I had added values for this, I remove
them in a future commit, more on that below).
- Updates the example in [example](example) to use the Well-Known Plugins.

### 9be25c78f8437859f1cc83206bea61deff96c187 Fix bugs and add GenPluginTypeUnset

This commit fixes some bugs that came up on testing, and also adds the `settings.GenPluginTypeUnset` constant. We should differentiate
between unset and `none` for the plugin type, so that a user could override for example the Well-Known Plugin `gogoslick` to not have
the type `gogo`, which has the effect of suppressing the generated modifiers. This feature will need to be tested

### The rest of the commits

The rest of the commits are for this Markdown file. I wanted to be able to see the output in GitHub as I wrote this.

The extra values I added to `GenConfig` and `GenPlugin` are also removed, as I realized through writing this document that they did
not end up being used.

## Where this proposal stops short

This proposal does not include the following:

- Testing for the Well-Known Plugins. I expect that there will be significant changes to this proposal, so this does not yet add
the testing for the Well-Known Plugins. In the event we actually used this proposal as the base for the final product, this would
be added.
- Documentation. This does not add documentation beyond this document.
