package sourcify

import (
	"testing"
)

func TestEndpoints_Validation(t *testing.T) {
	tests := []struct {
		name     string
		endpoint Endpoint
		wantErr  bool
	}{
		{
			name:     "Valid EndpointCheckByAddresses",
			endpoint: EndpointCheckByAddresses,
			wantErr:  false,
		},
		{
			name: "Valid EndpointCheckAllByAddresses",
			endpoint: Endpoint{
				Name:           "Check All By Addresses",
				URI:            "/check-all-by-addresses",
				MoreInfo:       "https://docs.sourcify.dev/docs/api/server/check-all-by-addresses/",
				Method:         "GET",
				ParamType:      EndpointParamTypeQueryString,
				RequiredParams: []string{"addresses", "chainIds"},
				Params: []EndpointParam{
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
			endpoint: Endpoint{
				Name:           "Invalid Endpoint",
				URI:            "/invalid-endpoint",
				MoreInfo:       "https://docs.sourcify.dev/docs/api/invalid",
				Method:         "GET",
				ParamType:      EndpointParamTypeQueryString,
				RequiredParams: []string{"missingParam"},
				Params: []EndpointParam{
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
			name: "Invalid Endpoint URI",
			endpoint: Endpoint{
				Name:           "Invalid Endpoint",
				URI:            "/aaa", // Empty URI
				MoreInfo:       "https://docs.sourcify.dev/docs/api/invalid",
				Method:         "GET",
				ParamType:      EndpointParamTypeQueryString,
				RequiredParams: []string{"test"},
				Params:         []EndpointParam{},
			},
			wantErr: true,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.endpoint.Verify()
			if (err != nil) != tt.wantErr {
				t.Errorf("Endpoint validation failed: %v, wantErr: %v", err, tt.wantErr)
			}
		})
	}
}
