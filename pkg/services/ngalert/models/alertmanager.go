package models

const AlertConfigurationVersion = 1

// AlertConfiguration represents a single version of the Alerting Engine Configuration.
type AlertConfiguration struct {
	ID int64 `xorm:"pk autoincr 'id'"`

	AlertmanagerConfiguration string
	ConfigurationHash         string
	ConfigurationVersion      string
	CreatedAt                 int64 `xorm:"created"`
	Default                   bool
	OrgID                     int64 `xorm:"org_id"`
	SuccessfullyAppliedAt     int64 `xorm:"successfully_applied_at"`
}

// GetLatestAlertmanagerConfigurationQuery is the query to get the latest alertmanager configuration.
type GetLatestAlertmanagerConfigurationQuery struct {
	OrgID  int64
	Result *AlertConfiguration
}

// GetSuccessfulAlertmanagerConfigurationsQuery is the query
// to get configurations that were successfully applied in the past for a given organization.
type GetSuccessfulAlertmanagerConfigurationsQuery struct {
	OrgID  int64
	Limit  int
	Result []*AlertConfiguration
}

// SaveAlertmanagerConfigurationCmd is the command to save an alertmanager configuration.
type SaveAlertmanagerConfigurationCmd struct {
	AlertmanagerConfiguration string
	FetchedConfigurationHash  string
	ConfigurationVersion      string
	Default                   bool
	OrgID                     int64
	ResultID                  int64
}
