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
)

func checkMessagesNotDeleted(addFailure func(*text.Failure), from *descriptor.FileDescriptorSet, to *descriptor.FileDescriptorSet) error {
	fromMessages, err := getMessages(from)
	if err != nil {
		return err
	}
	toMessages, err := getMessages(to)
	if err != nil {
		return err
	}
	for messageName, fromMessage := range fromMessages {
		if _, ok := toMessages[messageName]; !ok {
			addFailure(newTextFailuref(fromMessage.FileDescriptorProto, "Message %q was deleted.", messageName))
		}
	}
	return nil
}

func checkMessageFieldsNotDeleted(addFailure func(*text.Failure), from *descriptor.FileDescriptorSet, to *descriptor.FileDescriptorSet) error {
	fromMessages, err := getMessages(from)
	if err != nil {
		return err
	}
	toMessages, err := getMessages(to)
	if err != nil {
		return err
	}
	for _, messageName := range getIntersectingMessageNames(fromMessages, toMessages) {
		fromMessage := fromMessages[messageName]
		toMessage := toMessages[messageName]
		fromFields, err := getFields(fromMessage)
		if err != nil {
			return err
		}
		toFields, err := getFields(toMessage)
		if err != nil {
			return err
		}
		for fieldTag, fromField := range fromFields {
			if _, ok := toFields[fieldTag]; !ok {
				addFailure(newTextFailuref(fromField.Message.FileDescriptorProto, "Field %d on message %q was deleted.", fieldTag, messageName))
			}
		}
	}
	return nil
}

func checkMessageFieldsHaveSameType(addFailure func(*text.Failure), from *descriptor.FileDescriptorSet, to *descriptor.FileDescriptorSet) error {
	fromMessages, err := getMessages(from)
	if err != nil {
		return err
	}
	toMessages, err := getMessages(to)
	if err != nil {
		return err
	}
	for _, messageName := range getIntersectingMessageNames(fromMessages, toMessages) {
		fromMessage := fromMessages[messageName]
		toMessage := toMessages[messageName]
		fromFields, err := getFields(fromMessage)
		if err != nil {
			return err
		}
		toFields, err := getFields(toMessage)
		if err != nil {
			return err
		}
		for _, fieldTag := range getIntersectingFieldTags(fromFields, toFields) {
			fromField := fromFields[fieldTag]
			toField := toFields[fieldTag]
			fromFieldType := fromField.GetType()
			toFieldType := toField.GetType()
			if fromFieldType != toFieldType {
				addFailure(
					newTextFailuref(
						fromField.Message.FileDescriptorProto,
						"Field %d on message %q changed type from %q to %q.",
						fieldTag,
						messageName,
						fromFieldType,
						toFieldType,
					),
				)
			}
			if fromFieldType == descriptor.FieldDescriptorProto_TYPE_MESSAGE || fromFieldType == descriptor.FieldDescriptorProto_TYPE_ENUM {
				fromFieldTypeName := fromField.GetTypeName()
				toFieldTypeName := toField.GetTypeName()
				if fromFieldTypeName != toFieldTypeName {
					addFailure(
						newTextFailuref(
							fromField.Message.FileDescriptorProto,
							"Field %d on message %q changed type from %q to %q.",
							fieldTag,
							messageName,
							fromFieldTypeName,
							toFieldTypeName,
						),
					)
				}
			}
		}
	}
	return nil
}

// *** HELPER TYPES ***

type message struct {
	*descriptor.DescriptorProto

	FileDescriptorProto *descriptor.FileDescriptorProto
}

type field struct {
	*descriptor.FieldDescriptorProto

	Message *message
}

// *** HELPERS ***

// getMessages gets all Messages from the fileDescriptorSet.
//
// Returns a map from message name to message, or error if there is a system error.
func getMessages(fileDescriptorSet *descriptor.FileDescriptorSet) (map[string]*message, error) {
	messageNameToMessage := make(map[string]*message)
	for _, fileDescriptorProto := range fileDescriptorSet.File {
		if err := populateMessagesFromDescriptorProtos(messageNameToMessage, "", fileDescriptorProto, fileDescriptorProto.GetMessageType()); err != nil {
			return nil, err
		}
	}
	return messageNameToMessage, nil
}

// getIntersectingMessageNames gets the message names in both from and to.
//
// The checkMessagesNotDeleted check function  will verify the existence of messages that should exist.
// This is for other check functions.
func getIntersectingMessageNames(from map[string]*message, to map[string]*message) []string {
	var intersection []string
	for messageName := range from {
		if _, ok := to[messageName]; ok {
			intersection = append(intersection, messageName)
		}
	}
	return intersection
}

// getFields gets all fields from the message.
//
// Returns a map from tag to field, or error if there is a system error.
// Does not recurse into nested messages.
func getFields(message *message) (map[int32]*field, error) {
	fieldTagToField := make(map[int32]*field)
	for _, fieldDescriptorProto := range message.Field {
		tag := fieldDescriptorProto.GetNumber()
		if tag == 0 {
			return nil, fmt.Errorf("tag empty")
		}
		fieldTagToField[tag] = &field{
			FieldDescriptorProto: fieldDescriptorProto,
			Message:              message,
		}
	}
	return fieldTagToField, nil
}

// getIntersectingFieldTags gets the field tags in both from and to.
//
// The checkMessageFieldsNotDeleted check function will verify the existence of fields that should exist.
// This is for other check functions.
func getIntersectingFieldTags(from map[int32]*field, to map[int32]*field) []int32 {
	var intersection []int32
	for fieldTag := range from {
		if _, ok := to[fieldTag]; ok {
			intersection = append(intersection, fieldTag)
		}
	}
	return intersection
}

// newTextFailuref is a helper for check functions.
//
// The LintID will be populated by the Runner.
func newTextFailuref(fileDescriptorProto *descriptor.FileDescriptorProto, format string, args ...interface{}) *text.Failure {
	return &text.Failure{
		Filename: fileDescriptorProto.GetName(),
		Message:  fmt.Sprintf(format, args...),
	}
}

// *** SUB-HELPERS *** //
// *** SHOULD ONLY BE CALLED BY HELPERS *** //

func populateMessagesFromDescriptorProtos(
	messageNameToMessage map[string]*message,
	nestedName string,
	fileDescriptorProto *descriptor.FileDescriptorProto,
	descriptorProtos []*descriptor.DescriptorProto) error {
	for _, descriptorProto := range descriptorProtos {
		name := descriptorProto.GetName()
		if name == "" {
			// we do sanity checks like this just to confirm assumptions
			// this should never be hit
			return fmt.Errorf("name empty")
		}
		if nestedName != "" {
			name = nestedName + "." + name
		}
		messageNameToMessage[name] = &message{
			DescriptorProto:     descriptorProto,
			FileDescriptorProto: fileDescriptorProto,
		}
		if err := populateMessagesFromDescriptorProtos(messageNameToMessage, name, fileDescriptorProto, descriptorProto.GetNestedType()); err != nil {
			return err
		}
	}
	return nil
}
