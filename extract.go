/*
   Copyright 2021 - protosam
   Source can be found at https://github.com/protosam/opts

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

// Package opts provides a genericised means to pass options to functions.
//
// opts uses custom type names to fill struct fields based on the value of a tag
// called optname. So if an option of type WithUsername is passed and the struct
// has a field called Username tagged optname:"WithUsername", this will populate
// Username with the value of WithUsername.
package opts

import (
	"fmt"
	"reflect"
	"strings"
)

// Extract options into dest struct. Options not in dest are skipped.
func Extract(dest interface{}, options ...interface{}) error {
	return extract(dest, false, options...)
}

// Extract options into dest struct. Options not in dest result in error.
func MustExtract(dest interface{}, options ...interface{}) error {
	return extract(dest, true, options...)
}

// Underlying extract function.
func extract(dest interface{}, mustFind bool, options ...interface{}) error {
	// reflection of destination
	optionStruct := reflect.ValueOf(dest)
	// the destination must be addressable to make changes
	if optionStruct.Kind() == reflect.Ptr || optionStruct.Kind() == reflect.Interface {
		optionStruct = optionStruct.Elem()
	}

	// it must be a struct
	if optionStruct.Kind() != reflect.Struct {
		return fmt.Errorf("dest must be a struct")
	}

	// map all the optnames to struct field names
	fieldNameMap := make(map[string]string)
	for i := 0; i < optionStruct.NumField(); i++ {
		// use optname tags
		optname := optionStruct.Type().Field(i).Tag.Get("optname")
		if optname == "" {
			continue
		}
		// make sure this option is not already in use
		if _, found := fieldNameMap[optname]; found {
			return fmt.Errorf("option name %s has multiple tagged fields", optname)
		}
		// store for assignments
		fieldNameMap[optname] = optionStruct.Type().Field(i).Name
	}

	// iterate the options to assign them
	for i := 0; i < len(options); i++ {
		// reflect the option
		optionValue := reflect.ValueOf(options[i])
		// transform the tag to just the type without package name
		extracter := strings.Split(optionValue.Type().String(), ".")
		optname := extracter[len(extracter)-1]

		// find the fieldName
		fieldName, found := fieldNameMap[optname]
		if !found {
			// skip this value when finding it is not required
			if !mustFind {
				continue
			}
			return fmt.Errorf("invalid option %s", optname)
		}

		// fit the optionValue as exact match
		if optionStruct.FieldByName(fieldName).Type().Kind() == optionValue.Kind() {
			optionValue = optionValue.Convert(optionStruct.FieldByName(fieldName).Type())
			optionStruct.FieldByName(fieldName).Set(optionValue)
			// fit has been made, skip to next
			continue
		}

		// fit the optionValue by appending into a slice
		if optionStruct.FieldByName(fieldName).Type().Kind() == reflect.Slice && optionStruct.FieldByName(fieldName).Type().Elem().Kind() == optionValue.Kind() {
			optionValue = optionValue.Convert(optionStruct.FieldByName(fieldName).Type().Elem())
			optionStruct.FieldByName(fieldName).Set(reflect.Append(optionStruct.FieldByName(fieldName), optionValue))
			// fit has been made, skip to next
			continue
		}

		// failed to find fit
		return fmt.Errorf("failed to set %s when fitting %s into %s", optname, optionStruct.FieldByName(fieldName).Type().Kind().String(), optionValue.Kind().String())

	}
	return nil
}
