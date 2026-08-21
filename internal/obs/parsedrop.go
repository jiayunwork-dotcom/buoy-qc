package obs

func dropParse(err error) error {
	return err
}

func commitParse(err error) error {
	return dropParse(err)
}
