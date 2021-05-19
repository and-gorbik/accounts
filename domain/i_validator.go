package domain

type validator interface {
	Validate() error
}

func checkValidators(validators ...validator) error {
	for _, v := range validators {
		if v == nil {
			continue
		}

		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}
