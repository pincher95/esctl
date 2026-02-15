package utils

// FlattenSettingsMap flattens nested settings into dot-delimited keys.
// Example: {"index":{"number_of_shards":1}} -> {"index.number_of_shards":1}
func FlattenSettingsMap(settings map[string]any) map[string]any {
	flat := make(map[string]any)
	if settings == nil {
		return flat
	}
	flattenSettings("", settings, flat)
	return flat
}

func flattenSettings(prefix string, settings map[string]any, out map[string]any) {
	for key, value := range settings {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}
		switch v := value.(type) {
		case map[string]any:
			flattenSettings(fullKey, v, out)
		case map[string]string:
			nested := make(map[string]any, len(v))
			for k, val := range v {
				nested[k] = val
			}
			flattenSettings(fullKey, nested, out)
		default:
			out[fullKey] = value
		}
	}
}
