// Package sourcify provides tools for interacting with the Sourcify API.
package sourcify

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

type MethodParamType int
type MethodMatchType string

const (
	// MethodParamTypeUri denotes the type of parameter which is part of the URI.
	MethodParamTypeUri MethodParamType = iota // 0

	// MethodParamTypeQueryString denotes the type of parameter which is part of the query string.
	MethodParamTypeQueryString // 1

	// MethodMatchTypeFull denotes the type of match which is full.
	MethodMatchTypeFull MethodMatchType = "full"

	// MethodMatchTypePartial denotes the type of match which is partial.
	MethodMatchTypeAny MethodMatchType = "any"

	// MethodMatchTypePartial denotes the type of match which is partial.
	MethodMatchTypePartial MethodMatchType = "partial"
)

// String returns a string representation of the MethodParamType.
func (t MethodParamType) String() string {
	switch t {
	case MethodParamTypeUri:
		return "MethodParamTypeUri"
	case MethodParamTypeQueryString:
		return "MethodParamTypeQueryString"
	default:
		return fmt.Sprintf("Unknown MethodParamType (%d)", t)
	}
}

// MethodParam represents a parameter key-value pair.
type MethodParam struct {
	Key   string
	Value interface{}
}

// String returns a string representation of the MethodParam.
func (p MethodParam) String() string {
	return fmt.Sprintf("MethodParam{Key: %q, Value: %q}", p.Key, fmt.Sprintf("%v", p.Value))
}

// Method represents an API endpoint in the Sourcify service.
// It includes the name, the HTTP method, the URI, and any necessary parameters for requests to this endpoint.
type Method struct {
	Name           string
	Method         string
	URI            string
	MoreInfo       string
	ParamType      MethodParamType
	RequiredParams []string
	Params         []MethodParam
}

// GetParams returns a slice of the parameters for the API endpoint.
func (e Method) GetParams() []MethodParam {
	return e.Params
}

// GetQueryParams returns the query parameters for the API endpoint as a url.Values object.
func (e Method) GetQueryParams() url.Values {
	params := url.Values{}
	for _, param := range e.Params {
		if e.ParamType == MethodParamTypeQueryString {
			switch v := param.Value.(type) {
			case []string:
				params.Add(param.Key, strings.Join(v, ","))
			case []int:
				paramsString := []string{}
				for _, value := range v {
					paramsString = append(paramsString, fmt.Sprintf("%d", value))
				}
				params.Add(param.Key, strings.Join(paramsString, ","))
			case string:
				params.Set(param.Key, v)
			case int:
				params.Add(param.Key, fmt.Sprintf("%d", v))
			}
		}
	}
	return params
}

// SetParams allows setting parameters for the API endpoint using a variadic list of MethodParam values.
func (e *Method) SetParams(params ...MethodParam) {
	e.Params = params
}

// Verify checks if all the required parameters for the API endpoint are provided.
// It returns an error if any of the required parameters is missing.
func (e Method) Verify() error {
	for _, param := range e.RequiredParams {
		found := false
		for _, endpointParam := range e.Params {
			if param == endpointParam.Key {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("missing required parameter: %s", param)
		}
	}
	return nil
}

// ParseUri constructs a URL for the API endpoint, including any necessary query parameters or URI replacements.
// It can handle parameters of type string, int, []string, and []int. For []string and []int,
// the individual elements are joined with commas for the query string.
// Other types of parameters will trigger an error.
func (e Method) ParseUri() (string, error) {
	switch e.ParamType {
	case MethodParamTypeQueryString:
		var toReturn string

		// Add the parameters to the URL
		params := []string{}
		for _, param := range e.Params {
			switch v := param.Value.(type) {
			case []string:
				if len(v) > 0 {
					params = append(params, fmt.Sprintf("%s=%s", param.Key, strings.Join(v, ",")))
				}
			case []int:
				if len(v) > 0 {
					strs := []string{}
					for _, i := range v {
						strs = append(strs, fmt.Sprintf("%d", i))
					}
					params = append(params, fmt.Sprintf("%s=%s", param.Key, strings.Join(strs, ",")))
				}
			case string:
				if v != "" {
					params = append(params, fmt.Sprintf("%s=%s", param.Key, v))
				}
			case int:
				params = append(params, fmt.Sprintf("%s=%d", param.Key, v))
			default:
				// Return an error when encountering unsupported parameter type
				return "", ErrInvalidParamType(reflect.TypeOf(v).String())
			}
		}

		if len(params) > 0 {
			toReturn = fmt.Sprintf("%s?%s", toReturn, strings.Join(params, "&"))
		}

		return toReturn, nil

	case MethodParamTypeUri:
		toReturn := e.URI
		for _, param := range e.Params {
			switch v := param.Value.(type) {
			case string, int:
				toReturn = strings.ReplaceAll(toReturn, param.Key, fmt.Sprintf("%v", v))
			default:
				// Return an error when encountering unsupported parameter type
				return "", ErrInvalidParamType(reflect.TypeOf(v).String())
			}
		}
		return toReturn, nil

	default:
		return "", fmt.Errorf("invalid MethodParamType: %v", e.ParamType)
	}
}

// String returns a string representation of the Method struct.
func (m Method) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Method{\n    Name: %q,\n    Method: %q,\n    URI: %q,\n    MoreInfo: %q,\n    ParamType: %s,\n    RequiredParams: %v,\n    Params: [\n", m.Name, m.Method, m.URI, m.MoreInfo, methodParamTypeToString(m.ParamType), m.RequiredParams))
	for _, param := range m.Params {
		sb.WriteString(fmt.Sprintf("        %s,\n", param))
	}
	sb.WriteString("    ],\n}")
	return sb.String()
}

// methodParamTypeToString converts MethodParamType to a string representation.
func methodParamTypeToString(pt MethodParamType) string {
	switch pt {
	case MethodParamTypeUri:
		return "MethodParamTypeUri"
	case MethodParamTypeQueryString:
		return "MethodParamTypeQueryString"
	default:
		return "Unknown"
	}
}
