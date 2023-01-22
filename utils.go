package tracelog

func getValuesFromMap(fieldsMap Fields, fields []string) []interface{} {
	vals := make([]interface{}, len(fields))
	for i, field := range fields {
		var ok bool
		vals[i], ok = fieldsMap[field]
		if !ok {
			vals[i] = "???"
		}
	}

	return vals
}
