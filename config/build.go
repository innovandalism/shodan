package config

const (
	// VersionMajor is the major component of the semver
	VersionMajor int = 1
	// VersionMinor is the minor component of the semver
	VersionMinor int = 0
	// VersionRevision is the revision component of the semver
	VersionRevision int = 5
)

var (
	// SingleGuildMode is a string that is either empty of the snowflake of a Discord guild.
	// If it is the later, the API will only serve that one server.
	SingleGuildMode = ""

	// OEMName is the name of the bot shown in the version command
	OEMName = "SHODAN"
	// OEMDescription is the description of the bot shown in the version command
	OEMDescription = "Core build of SHODAN, the Discord Bot Framework"
	// OEMVendor is the vendor name shown in the version command
	OEMVendor = "Innovandalism"
	// OEMURI is the URI for the name shown in the version command
	OEMURI = "https://github.com/innovandalism/shodan/"
	// OEMVendorURI is the URI shown for the vendor name in the version command
	OEMVendorURI = "https://innovandalism.eu"
	// OEMBuildfile is the URI to the buildfile shown in the version command
	OEMBuildfile = "https://github.com/innovandalism/shodan/"
	// OEMAnecdote is a witty line shown at the bottom of the version command
	OEMAnecdote = "This build of Shodan has super-kapow-powers."
)
