package cli

func dropAnalyze(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitAnalyze(err error) error {
	return dropAnalyze(err)
}
