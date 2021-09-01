package forms

type errors map[string][]string

func (errors errors) Add(field, message string) {
	errors[field] = append(errors[field], message)
}

func (errors errors) Get(field string) string {
	errorMessages := errors[field]

	if len(errorMessages) == 0 {
		return ""
	}

	return errorMessages[0]
}
