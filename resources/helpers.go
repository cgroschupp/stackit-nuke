package resources

import "time"

func parseTimePtr(v interface{}) *time.Time {
	switch t := v.(type) {
	case time.Time:
		return &t
	case *time.Time:
		return t
	case string:
		if t == "" {
			return nil
		}
		if parsed, err := time.Parse(time.RFC3339, t); err == nil {
			return &parsed
		}
	}
	return nil
}

func deref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func safeName(name, fallback *string) string {
	if name != nil && *name != "" {
		return *name
	}
	if fallback != nil {
		return *fallback
	}
	return "<unknown>"
}
