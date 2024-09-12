// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// AuthResponse auth response
//
// swagger:model AuthResponse
type AuthResponse struct {

	// Код ответа
	Code int64 `json:"code,omitempty"`

	// Детали ответа
	Detail string `json:"detail,omitempty"`

	// Сообщение на английском языке
	EnMessage string `json:"en_message,omitempty"`

	// Сообщение
	Message string `json:"message,omitempty"`
}

// Validate validates this auth response
func (m *AuthResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this auth response based on context it is used
func (m *AuthResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *AuthResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AuthResponse) UnmarshalBinary(b []byte) error {
	var res AuthResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
