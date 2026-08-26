package domain

func ValidateProfile(p Profile) error {
	if p.ID == "" || p.Name == "" {
		return ErrInactiveProfile
	}
	return nil
}
func EnsureUnique(ids []string) bool {
	seen := map[string]bool{}
	for _, id := range ids {
		if id == "" || seen[id] {
			return false
		}
		seen[id] = true
	}
	return true
}
func SumRecords(rs []Record) float64 {
	total := 0.0
	for _, r := range rs {
		total += r.Amount()
	}
	return total
}
