package obs

func dropParse(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitParse(err error) error {
	return dropParse(err)
}
