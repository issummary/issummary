package issummary

type Config struct {
	Port              int
	Token             string
	GitServiceBaseURL string
	GitServiceType    string
	GIDs              []string
	SPLabelPrefix     string
	ClassLabelPrefix  string
}
