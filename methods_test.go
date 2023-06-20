package sourcify

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMethods_Validation(t *testing.T) {
	tests := []struct {
		name     string
		endpoint Method
		wantErr  bool
	}{
		{
			name:     "Valid MethodCheckByAddresses",
			endpoint: MethodCheckByAddresses,
			wantErr:  false,
		},
		{
			name: "Valid MethodCheckAllByAddresses",
			endpoint: Method{
				Name:           "Check All By Addresses",
				URI:            "/check-all-by-addresses",
				MoreInfo:       "https://docs.sourcify.dev/docs/api/server/check-all-by-addresses/",
				Method:         "GET",
				ParamType:      MethodParamTypeQueryString,
				RequiredParams: []string{"addresses", "chainIds"},
				Params: []MethodParam{
					{
						Key:   "addresses",
						Value: []string{},
					},
					{
						Key:   "chainIds",
						Value: []int{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Missing Required Parameter",
			endpoint: Method{
				Name:           "Invalid Method",
				URI:            "/invalid-endpoint",
				MoreInfo:       "https://docs.sourcify.dev/docs/api/invalid",
				Method:         "GET",
				ParamType:      MethodParamTypeQueryString,
				RequiredParams: []string{"missingParam"},
				Params: []MethodParam{
					{
						Key:   "existingParam",
						Value: "value",
					},
					// Missing required parameter "missingParam"
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid Method URI",
			endpoint: Method{
				Name:           "Invalid Method",
				URI:            "/aaa", // Empty URI
				MoreInfo:       "https://docs.sourcify.dev/docs/api/invalid",
				Method:         "GET",
				ParamType:      MethodParamTypeQueryString,
				RequiredParams: []string{"test"},
				Params:         []MethodParam{},
			},
			wantErr: true,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.endpoint.Verify()
			if (err != nil) != tt.wantErr {
				t.Errorf("Method validation failed: %v, wantErr: %v", err, tt.wantErr)
			}
		})
	}
}

func TestMethod_GetParams(t *testing.T) {
	params := []MethodParam{
		{Key: "param1", Value: "value1"},
		{Key: "param2", Value: 123},
	}

	method := Method{
		Params: params,
	}

	assert.Equal(t, params, method.GetParams())
}

func TestMethod_GetQueryParams(t *testing.T) {
	params := []MethodParam{
		{Key: "param1", Value: []string{"value1", "value2"}},
		{Key: "param2", Value: []int{1, 2, 3}},
	}

	method := Method{
		Params:         params,
		ParamType:      MethodParamTypeQueryString,
		RequiredParams: []string{"param1", "param2"},
	}

	expected := url.Values{
		"param1": []string{"value1,value2"},
		"param2": []string{"1,2,3"},
	}

	assert.Equal(t, expected, method.GetQueryParams())
}

func TestMethod_GetQueryParams_NoParams(t *testing.T) {
	method := Method{
		Params:         []MethodParam{},
		ParamType:      MethodParamTypeQueryString,
		RequiredParams: []string{},
	}

	expected := url.Values{}
	assert.Equal(t, expected, method.GetQueryParams())
}

func TestMethod_SetParams(t *testing.T) {
	method := Method{}

	params := []MethodParam{
		{Key: "param1", Value: "value1"},
		{Key: "param2", Value: 123},
	}
	method.SetParams(params...)

	assert.Equal(t, params, method.Params)
}

func TestMethod_Verify(t *testing.T) {
	requiredParams := []string{"param1", "param2"}

	// Missing param1
	method := Method{
		Params:         []MethodParam{{Key: "param2", Value: "value2"}},
		RequiredParams: requiredParams,
	}
	err := method.Verify()
	assert.EqualError(t, err, "missing required parameter: param1")

	// Missing param2
	method = Method{
		Params:         []MethodParam{{Key: "param1", Value: "value1"}},
		RequiredParams: requiredParams,
	}
	err = method.Verify()
	assert.EqualError(t, err, "missing required parameter: param2")

	// All required params present
	method = Method{
		Params:         []MethodParam{{Key: "param1", Value: "value1"}, {Key: "param2", Value: "value2"}},
		RequiredParams: requiredParams,
	}
	err = method.Verify()
	assert.NoError(t, err)
}

func TestMethod_ParseUri_QueryString(t *testing.T) {
	method := Method{
		ParamType: MethodParamTypeQueryString,
		Params: []MethodParam{
			{Key: "param1", Value: "value1"},
			{Key: "param2", Value: 123},
		},
	}
	expected, err := url.ParseQuery("param1=value1&param2=123")
	assert.NoError(t, err)
	parsed, err := method.ParseUri()
	assert.NoError(t, err)
	assert.Equal(t, "?"+expected.Encode(), parsed)
}

func TestMethod_ParseUri_Uri(t *testing.T) {
	method := Method{
		ParamType: MethodParamTypeUri,
		URI:       "/files/:chain/:address",
		Params: []MethodParam{
			{Key: ":chain", Value: "mainnet"},
			{Key: ":address", Value: "0x1234567890abcdef"},
		},
	}

	expected := "/files/mainnet/0x1234567890abcdef"
	parsed, err := method.ParseUri()
	assert.NoError(t, err)
	assert.Equal(t, expected, parsed)
}

func TestMethod_ParseUri_InvalidParamType(t *testing.T) {
	method := Method{
		ParamType: MethodParamTypeQueryString,
		Params: []MethodParam{
			{Key: "param1", Value: true},
		},
	}

	_, err := method.ParseUri()
	assert.EqualError(t, err, "encountered a parameter of invalid type: bool")
}

func TestMethod_ParseUri_InvalidMethodParamType(t *testing.T) {
	method := Method{
		ParamType: MethodParamType(123),
		Params:    []MethodParam{},
	}

	_, err := method.ParseUri()
	assert.EqualError(t, err, "invalid MethodParamType: Unknown MethodParamType (123)")
}

func TestMethod_ParseUri_EmptyParams(t *testing.T) {
	method := Method{
		ParamType: MethodParamTypeQueryString,
		Params:    []MethodParam{},
	}

	expected := ""
	parsed, err := method.ParseUri()
	assert.NoError(t, err)
	assert.Equal(t, expected, parsed)
}

func TestMethod_ParseUri_EmptyValueParams(t *testing.T) {
	method := Method{
		ParamType: MethodParamTypeUri,
		URI:       "/files/:chain/:address",
		Params: []MethodParam{
			{Key: ":chain", Value: ""},
			{Key: ":address", Value: ""},
		},
	}

	expected := "/files//"
	parsed, err := method.ParseUri()
	assert.NoError(t, err)
	assert.Equal(t, expected, parsed)
}

func TestMethod_ParseUri_UnsupportedParamType(t *testing.T) {
	method := Method{
		ParamType: MethodParamTypeQueryString,
		Params: []MethodParam{
			{Key: "param1", Value: []bool{true, false}},
		},
	}

	_, err := method.ParseUri()
	assert.EqualError(t, err, "encountered a parameter of invalid type: []bool")
}

func TestMethodParamType(t *testing.T) {
	assert.Equal(t, MethodParamType(0), MethodParamTypeUri)
	assert.Equal(t, MethodParamType(1), MethodParamTypeQueryString)
}

func TestMethodParamType_String(t *testing.T) {
	assert.Equal(t, "MethodParamTypeUri", MethodParamTypeUri.String())
	assert.Equal(t, "MethodParamTypeQueryString", MethodParamTypeQueryString.String())
}

func TestErrInvalidParamType_Error(t *testing.T) {
	err := ErrInvalidParamType("int")
	assert.EqualError(t, err, "encountered a parameter of invalid type: int")
}

func TestMethodParam_String(t *testing.T) {
	param := MethodParam{
		Key:   "param1",
		Value: "value1",
	}
	expected := `MethodParam{Key: "param1", Value: "value1"}`
	assert.Equal(t, expected, param.String())
}

func TestMethod_String(t *testing.T) {
	method := Method{
		Name:           "Test Method",
		Method:         "GET",
		URI:            "/test",
		MoreInfo:       "https://example.com",
		ParamType:      MethodParamTypeQueryString,
		RequiredParams: []string{"param1", "param2"},
		Params: []MethodParam{
			{
				Key:   "param1",
				Value: "value1",
			},
			{
				Key:   "param2",
				Value: "123",
			},
		},
	}

	expected := `Method{
    Name: "Test Method",
    Method: "GET",
    URI: "/test",
    MoreInfo: "https://example.com",
    ParamType: MethodParamTypeQueryString,
    RequiredParams: [param1 param2],
    Params: [
        MethodParam{Key: "param1", Value: "value1"},
        MethodParam{Key: "param2", Value: "123"},
    ],
}`
	assert.Equal(t, expected, method.String())
}
