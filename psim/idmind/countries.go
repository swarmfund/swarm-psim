package idmind

// TODO
func convertToISO(countryFromEnum string) string {
	if len(countryFromEnum) == 2 {
		// In case country is already of type ISO (2 letter)
		return countryFromEnum
	}

	// TODO
	return countryFromEnum[:2]
}
