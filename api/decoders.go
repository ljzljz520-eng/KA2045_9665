package api

import (
	"encoding/json"
	"io"
	"workpay/domain"
)

func DecodeRecord(r io.Reader) (domain.Record, error) {
	var v domain.Record
	e := json.NewDecoder(r).Decode(&v)
	return v, e
}
func DecodeProfile(r io.Reader) (domain.Profile, error) {
	var v domain.Profile
	e := json.NewDecoder(r).Decode(&v)
	return v, e
}
func EncodeError(err error) map[string]string {
	if err == nil {
		return map[string]string{}
	}
	return map[string]string{"error": err.Error()}
}
func StatusCode(err error) int {
	if err == nil {
		return 200
	}
	if err == domain.ErrPeriodClosed {
		return 409
	}
	return 400
}
