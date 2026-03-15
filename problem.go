package respond

import (
	"encoding/json"
	"net/http"
)

// ProblemDetails represents an RFC 9457 Problem Details object.
type ProblemDetails struct {
	// Type is a URI reference that identifies the problem type.
	Type string `json:"type,omitempty"`

	// Title is a short, human-readable summary of the problem type.
	Title string `json:"title,omitempty"`

	// Status is the HTTP status code for this problem occurrence.
	Status int `json:"status"`

	// Detail is a human-readable explanation specific to this occurrence.
	Detail string `json:"detail,omitempty"`

	// Instance is a URI reference that identifies the specific occurrence.
	Instance string `json:"instance,omitempty"`

	// Extensions holds additional members of the problem details object.
	Extensions map[string]any `json:"-"`
}

// MarshalJSON implements custom JSON marshalling for ProblemDetails to
// include extension fields as top-level members alongside the standard fields.
func (p ProblemDetails) MarshalJSON() ([]byte, error) {
	type plain ProblemDetails
	base, err := json.Marshal(plain(p))
	if err != nil {
		return nil, err
	}

	if len(p.Extensions) == 0 {
		return base, nil
	}

	ext, err := json.Marshal(p.Extensions)
	if err != nil {
		return nil, err
	}

	// Merge: strip trailing } from base, strip leading { from ext, join with comma.
	merged := make([]byte, 0, len(base)+len(ext))
	merged = append(merged, base[:len(base)-1]...)
	merged = append(merged, ',')
	merged = append(merged, ext[1:]...)

	return merged, nil
}

// ProblemOption is a functional option for configuring a ProblemDetails response.
type ProblemOption func(*ProblemDetails)

// WithType sets the type URI on the problem details.
func WithType(uri string) ProblemOption {
	return func(p *ProblemDetails) {
		p.Type = uri
	}
}

// WithTitle sets the title on the problem details.
func WithTitle(t string) ProblemOption {
	return func(p *ProblemDetails) {
		p.Title = t
	}
}

// WithDetail sets the detail message on the problem details.
func WithDetail(d string) ProblemOption {
	return func(p *ProblemDetails) {
		p.Detail = d
	}
}

// WithInstance sets the instance URI on the problem details.
func WithInstance(uri string) ProblemOption {
	return func(p *ProblemDetails) {
		p.Instance = uri
	}
}

// WithExtension adds a custom extension field to the problem details.
// Extension fields appear as top-level members in the JSON output.
func WithExtension(key string, value any) ProblemOption {
	return func(p *ProblemDetails) {
		if p.Extensions == nil {
			p.Extensions = make(map[string]any)
		}
		p.Extensions[key] = value
	}
}

// Problem writes an RFC 9457 Problem Details JSON response with the given
// status code and options. The Content-Type is set to application/problem+json.
// If marshalling fails, a 500 Internal Server Error is written as plain text.
func Problem(w http.ResponseWriter, status int, opts ...ProblemOption) {
	pd := &ProblemDetails{
		Status: status,
	}

	for _, opt := range opts {
		opt(pd)
	}

	body, err := json.Marshal(pd)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	w.Write(body)
}
