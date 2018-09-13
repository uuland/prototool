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

package compat

import (
	"fmt"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/uber/prototool/internal/text"
	"go.uber.org/zap"
)

type runner struct {
	logger   *zap.Logger
	checkers []Checker
}

func newRunner(options ...RunnerOption) *runner {
	runner := &runner{
		logger:   zap.NewNop(),
		checkers: AllCheckers,
	}
	for _, option := range options {
		option(runner)
	}
	return runner
}

func (r *runner) Run(from []*descriptor.FileDescriptorSet, to []*descriptor.FileDescriptorSet) ([]*text.Failure, error) {
	fromMap, err := getPackageToFileDescriptorSet(from)
	if err != nil {
		return nil, err
	}
	toMap, err := getPackageToFileDescriptorSet(to)
	if err != nil {
		return nil, err
	}
	for pkg := range fromMap {
		if _, ok := toMap[pkg]; !ok {
			delete(fromMap, pkg)
		}
	}

	var failures []*text.Failure
	for pkg, fromFileDescriptorSet := range fromMap {
		toFileDescriptorSet, ok := toMap[pkg]
		if !ok {
			return nil, fmt.Errorf("mismatch in file descriptor set maps, this should never happen: %v", pkg)
		}
		for _, checker := range r.checkers {
			var checkerFailures []*text.Failure
			if err := checker.Check(
				func(failure *text.Failure) {
					checkerFailures = append(checkerFailures, failure)
				},
				fromFileDescriptorSet,
				toFileDescriptorSet,
			); err != nil {
				return nil, err
			}
			for _, checkerFailure := range checkerFailures {
				checkerFailure.LintID = checker.ID
				failures = append(failures, checkerFailure)
			}
		}
	}
	return failures, nil
}

func getPackageToFileDescriptorSet(all []*descriptor.FileDescriptorSet) (map[string]*descriptor.FileDescriptorSet, error) {
	packageToFileDescriptorProtos := make(map[string]map[string]*descriptor.FileDescriptorProto)

	for _, fileDescriptorSet := range all {
		for _, fileDescriptorProto := range fileDescriptorSet.File {
			pkg := fileDescriptorProto.GetPackage()
			if pkg == "" {
				return nil, fmt.Errorf("compat requires all Protobuf files to have a package but %v does not", fileDescriptorProto)
			}
			m := packageToFileDescriptorProtos[pkg]
			if m == nil {
				m = make(map[string]*descriptor.FileDescriptorProto)
				packageToFileDescriptorProtos[pkg] = m
			}
			name := fileDescriptorProto.GetName()
			if name == "" {
				return nil, fmt.Errorf("compat requires all FileDescriptorProtos to have a file name but %v does not", fileDescriptorProto)
			}
			// this overwrites duplicates on purpose
			m[name] = fileDescriptorProto
		}
	}
	packageToFileDescriptorSet := make(map[string]*descriptor.FileDescriptorSet)
	for pkg, fileDescriptorProtos := range packageToFileDescriptorProtos {
		var files []*descriptor.FileDescriptorProto
		for _, file := range fileDescriptorProtos {
			files = append(files, file)
		}
		packageToFileDescriptorSet[pkg] = &descriptor.FileDescriptorSet{
			File: files,
		}
	}
	return packageToFileDescriptorSet, nil
}
