package issummary

type Config struct {
	Port                int
	Token               string
	GitServiceBaseURL   string
	GitServiceType      string
	Organizations       []string
	SPLabelPrefix       string
	ClassLabelPrefix    string
	TargetLabelPrefixes []string
}
