// Copyright (c) 2018 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package settings

import "sort"

var (
	// AliasToWellKnownPlugin is a map from alias to Well-Known Plugin.
	AliasToWellKnownPlugin = map[string]WellKnownPlugin{
		// Built-ins.
		"cpp": WellKnownPlugin{
			Name:          "cpp",
			RelOutputPath: "cpp",
		},
		"csharp": WellKnownPlugin{
			Name:          "csharp",
			RelOutputPath: "csharp",
		},
		"java": WellKnownPlugin{
			Name:          "java",
			RelOutputPath: "java",
		},
		"js": WellKnownPlugin{
			Name:          "js",
			RelOutputPath: "js",
		},
		"objc": WellKnownPlugin{
			Name:          "objc",
			RelOutputPath: "objc",
		},
		"php": WellKnownPlugin{
			Name:          "php",
			RelOutputPath: "php",
		},
		"python": WellKnownPlugin{
			Name:          "python",
			RelOutputPath: "python",
		},
		"ruby": WellKnownPlugin{
			Name:          "ruby",
			RelOutputPath: "ruby",
		},

		// Built-in gRPC variants.
		// These will need to be specified in addition to the base plugins,
		// ie for cpp+grpc, you need to specify for cpp and cpp-grpc.
		// Golang is different in that it has gRPC built in to the base plugins.
		// We could have WellKnownPluginSet that groups these, but that feels
		// like a concept that should be on top of WellKnownPlugin.
		"cpp-grpc": WellKnownPlugin{
			Name:          "cpp-grpc",
			Path:          "grpc_cpp_plugin",
			RelOutputPath: "cpp",
		},
		"csharp-grpc": WellKnownPlugin{
			Name:          "csharp-grpc",
			Path:          "grpc_csharp_plugin",
			RelOutputPath: "csharp",
		},
		"objc-grpc": WellKnownPlugin{
			Name:          "objc-grpc",
			Path:          "grpc_objective_c_plugin",
			RelOutputPath: "objc",
		},
		"php-grpc": WellKnownPlugin{
			Name:          "php-grpc",
			Path:          "grpc_php_plugin",
			RelOutputPath: "php",
		},
		"python-grpc": WellKnownPlugin{
			Name:          "python-grpc",
			Path:          "grpc_python_plugin",
			RelOutputPath: "python",
		},
		"ruby-grpc": WellKnownPlugin{
			Name:          "ruby-grpc",
			Path:          "grpc_ruby_plugin",
			RelOutputPath: "ruby",
		},

		// Golang base plugins.
		// Note we use the same relative output path for most of the Golang
		// plugins as you would expect the files to be in the same directory,
		// and the choice between go and gogo variants should not affect
		// output path.
		"go": WellKnownPlugin{
			Name:          "go",
			Type:          GenPluginTypeGo,
			RelOutputPath: "go",
		},
		"gogo": WellKnownPlugin{
			Name:          "gogo",
			Type:          GenPluginTypeGogo,
			RelOutputPath: "go",
		},
		"gofast": WellKnownPlugin{
			Name:          "gofast",
			Type:          GenPluginTypeGogo,
			RelOutputPath: "go",
		},
		"gogofast": WellKnownPlugin{
			Name:          "gogofast",
			Type:          GenPluginTypeGogo,
			RelOutputPath: "go",
		},
		"gogofaster": WellKnownPlugin{
			Name:          "gogofaster",
			Type:          GenPluginTypeGogo,
			RelOutputPath: "go",
		},
		"gogoslick": WellKnownPlugin{
			Name:          "gogoslick",
			Type:          GenPluginTypeGogo,
			RelOutputPath: "go",
		},

		// Golang base plugins with gRPC added.
		// These are different than other languages as you would specify
		// ie either "go" or "go-with-grpc".
		"go-with-grpc": WellKnownPlugin{
			Name:          "go",
			Type:          GenPluginTypeGo,
			Flags:         "plugins=grpc",
			RelOutputPath: "go",
		},
		"gogo-with-grpc": WellKnownPlugin{
			Name:          "gogo",
			Type:          GenPluginTypeGogo,
			Flags:         "plugins=grpc",
			RelOutputPath: "go",
		},
		"gofast-with-grpc": WellKnownPlugin{
			Name:          "gofast",
			Type:          GenPluginTypeGogo,
			Flags:         "plugins=grpc",
			RelOutputPath: "go",
		},
		"gogofast-with-grpc": WellKnownPlugin{
			Name:          "gogofast",
			Type:          GenPluginTypeGogo,
			Flags:         "plugins=grpc",
			RelOutputPath: "go",
		},
		"gogofaster-with-grpc": WellKnownPlugin{
			Name:          "gogofaster",
			Type:          GenPluginTypeGogo,
			Flags:         "plugins=grpc",
			RelOutputPath: "go",
		},
		"gogoslick-with-grpc": WellKnownPlugin{
			Name:          "gogoslick",
			Type:          GenPluginTypeGogo,
			Flags:         "plugins=grpc",
			RelOutputPath: "go",
		},

		// Special plugins we know about.
		// Note that grpc-gateway required extended Well-Known Types to be
		// present, we could enforce this, however users may want to specify
		// the paths to these themselves so we do no no enforcement within
		// the settings package.
		"yarpc-go": WellKnownPlugin{
			Name:          "yarpc-go",
			Type:          GenPluginTypeGogo,
			RelOutputPath: "go",
		},
		"grpc-gateway": WellKnownPlugin{
			Name:          "grpc-gateway",
			Type:          GenPluginTypeGo,
			RelOutputPath: "go",
		},
	}
)

// WellKnownPlugin is a Well-Known Plugin.
type WellKnownPlugin struct {
	// The name of the plugin. For example, if you want to use
	// protoc-gen-gogoslick, the name is "gogoslick".
	Name string
	// The path to the executable. For example, if the name is "grpc-cpp"
	// but the path to the executable "protoc-gen-grpc-cpp" is "/usr/local/bin/grpc_cpp_plugin",
	// then this will be "/usr/local/bin/grpc_cpp_plugin".
	Path string
	// The type, if any. This will be GenPluginTypeNone if
	// there is no specific type.
	Type GenPluginType
	// Extra flags to pass.
	// If there is an associated type, some flags may be generated,
	// for example plugins=grpc or Mfile=package modifiers.
	Flags string
	// The path to output to.
	// Must be relative in a config file.
	RelOutputPath string
}

// ForEachWellKnownPlugin iterates over the Well-Known Plugins in order, sorted by alias.
func ForEachWellKnownPlugin(f func(string, WellKnownPlugin) error) error {
	aliases := make([]string, 0, len(AliasToWellKnownPlugin))
	for alias := range AliasToWellKnownPlugin {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)
	for _, alias := range aliases {
		if err := f(alias, AliasToWellKnownPlugin[alias]); err != nil {
			return err
		}
	}
	return nil
}
