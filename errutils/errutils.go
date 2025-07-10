package errutils

func RootCause(err error) error {
	for {
		e, ok := err.(interface{ Unwrap() error })
		if !ok {
			return err
		}

		err = e.Unwrap()
		if err == nil {
			return nil
		}
	}
}
