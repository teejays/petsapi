package route

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRoutes(t *testing.T) {
	tests := []struct {
		name string
		want []Route
	}{
		{
			name: "should return all the routes",
			want: routes,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetRoutes()
			assert.Equal(t, len(tt.want), len(got))
		})
	}
}

func TestRoute_GetPattern(t *testing.T) {
	type fields struct {
		Method      string
		Version     int
		Path        string
		HandlerFunc http.HandlerFunc
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "should work",
			fields: fields{
				Method:      http.MethodGet,
				Version:     2,
				Path:        "someresource",
				HandlerFunc: nil,
			},
			want: "/v2/someresource",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Route{
				Method:      tt.fields.Method,
				Version:     tt.fields.Version,
				Path:        tt.fields.Path,
				HandlerFunc: tt.fields.HandlerFunc,
			}
			got := r.GetPattern()
			assert.Equal(t, tt.want, got)
		})
	}
}
