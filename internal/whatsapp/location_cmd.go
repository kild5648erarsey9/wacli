package whatsapp

// LocationOptions holds the parameters for sending a location message
// via CLI or programmatic use.
type LocationOptions struct {
	// Recipient is the destination phone number in E.164 format.
	Recipient string

	// Latitude is the geographic latitude of the location.
	Latitude float64

	// Longitude is the geographic longitude of the location.
	Longitude float64

	// Name is an optional human-readable name for the location.
	Name string

	// Address is an optional street address for the location.
	Address string
}

// Validate checks that all required LocationOptions fields are set
// and that coordinates are within valid ranges.
func (o LocationOptions) Validate() error {
	if o.Recipient == "" {
		return errorf("recipient phone number must not be empty")
	}
	if o.Latitude < -90 || o.Latitude > 90 {
		return errorf("latitude must be between -90 and 90, got %f", o.Latitude)
	}
	if o.Longitude < -180 || o.Longitude > 180 {
		return errorf("longitude must be between -180 and 180, got %f", o.Longitude)
	}
	return nil
}

// errorf is a thin wrapper around fmt.Errorf to keep imports minimal
// within this file while reusing the package-level fmt dependency.
func errorf(format string, args ...interface{}) error {
	if len(args) == 0 {
		return &locationError{msg: format}
	}
	return &locationError{msg: fmt.Sprintf(format, args...)}
}

type locationError struct{ msg string }

func (e *locationError) Error() string { return e.msg }
