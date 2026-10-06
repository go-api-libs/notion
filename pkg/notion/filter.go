package notion

// Flatten returns the filter with every compound filter within one of the same kind joined into it, an "or" in an
// "or" or an "and" in an "and", and every compound of one filter replaced by that filter. It allows the same entries,
// nested as little as it can be: a view's filter may nest deeper than the two levels a query allows.
func (f Filter) Flatten() Filter {
	switch {
	case f.FilterOr != nil:
		return compound(join(f.FilterOr.Or, func(c Filter) []Filter {
			if c.FilterOr != nil {
				return c.FilterOr.Or
			}

			return nil
		}), func(fs []Filter) Filter { return Filter{FilterOr: &FilterOr{Or: fs}} })
	case f.FilterAnd != nil:
		return compound(join(f.FilterAnd.And, func(c Filter) []Filter {
			if c.FilterAnd != nil {
				return c.FilterAnd.And
			}

			return nil
		}), func(fs []Filter) Filter { return Filter{FilterAnd: &FilterAnd{And: fs}} })
	default:
		return f
	}
}

// join flattens each of filters, and lists, in place of each that same holds filters of, those filters.
func join(filters []Filter, same func(Filter) []Filter) []Filter {
	joined := make([]Filter, 0, len(filters))

	for _, c := range filters {
		c = c.Flatten()
		if inner := same(c); inner != nil {
			joined = append(joined, inner...)
		} else {
			joined = append(joined, c)
		}
	}

	return joined
}

// compound is the one filter of filters, or else the compound of them that of makes.
func compound(filters []Filter, of func([]Filter) Filter) Filter {
	if len(filters) == 1 {
		return filters[0]
	}

	return of(filters)
}
