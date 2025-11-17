package resolvers

// normalizePagination kapselt die ganzen Limit/Offset-Checks.
// defaultLimit: Standardwert, falls limit nil oder ungültig ist
// maxLimit: Hard-Cap nach oben (z.B. 100)
func normalizePagination(limit *int32, offset *int32, defaultLimit, maxLimit int) (int, int) {
	l := defaultLimit
	if limit != nil {
		if *limit > 0 && int(*limit) <= maxLimit {
			l = int(*limit)
		}
	}

	o := 0
	if offset != nil && *offset >= 0 {
		o = int(*offset)
	}

	return l, o
}
