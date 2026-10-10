package notion

// IsZero reports whether dt is the zero value of DateOrDateTime.
//
// DateOrDateTime is an untagged anyOf union (at least one field is set after
// unmarshaling), so a value is zero only when every set field is itself zero:
//   - nil receiver
//   - both fields nil (the struct's zero value)
//   - only Date set and zero
//   - only Time set and zero
//   - both set and both zero
//
// In particular, if both fields are non-nil, both must be zero for this to
// return true.
func (dt *DateOrDateTime) IsZero() bool {
	return dt == nil ||
		(dt.Date == nil || dt.Date.IsZero()) &&
			(dt.Time == nil || dt.Time.IsZero())
}
