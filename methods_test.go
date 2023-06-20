package sourcify

import (
	"testing"
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
