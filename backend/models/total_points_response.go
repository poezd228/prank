// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// TotalPointsResponse total points response
//
// swagger:model TotalPointsResponse
type TotalPointsResponse struct {

	// value
	// Example: 100000
	Value int64 `json:"value,omitempty"`
}

// Validate validates this total points response
func (m *TotalPointsResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this total points response based on context it is used
func (m *TotalPointsResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *TotalPointsResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TotalPointsResponse) UnmarshalBinary(b []byte) error {
	var res TotalPointsResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
