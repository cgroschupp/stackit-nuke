package common

import "fmt"

type versionInfo struct {
	Version string
	Commit  string
	Date    string
	Summary string
}

var AppVersion = func() versionInfo {
	v := versionInfo{
		Version: "0.1.0-alpha",
		Commit:  "dev",
		Date:    "local",
	}
	v.Summary = fmt.Sprintf("%s (%s, %s)", v.Version, v.Commit, v.Date)
	return v
}()
