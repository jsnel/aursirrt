package datastorage

type persistor interface {
	PersistData(string,CommitData) error
}

func getPersistor(strategy string) persistor {
	switch strategy{
	case "log":
		return logPersistor{}
	}
	return nil
}
