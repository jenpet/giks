package meta

// CommitHashArg of the git revision
var CommitHashArg string

// VersionArg of the giks binary
var VersionArg string

func CommitHash() string {
	if CommitHashArg == "" {
		return "local-hash"
	}
	return CommitHashArg
}

func Version() string {
	if VersionArg == "" {
		return "local-version"
	}
	return VersionArg
}
