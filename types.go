package cronjoborg

import "encoding/json"

// using mostly this regex to map types https://regex101.com/r/JCchDA/1

// Job
// https://docs.cron-job.org/rest-api.html#job
type Job struct {
	// Job identifier
	JobID int `json:"jobId"`

	// Whether the job is enabled (i.e. being executed) or not
	Enabled bool `json:"enabled"`

	// Job title
	Title string `json:"title"`

	// Whether to save job response header/body or not
	SaveResponses bool `json:"saveResponses"`

	// Job URL
	URL string `json:"url"`

	// Last execution status
	LastStatus JobStatus `json:"lastStatus"`

	// Last execution duration in milliseconds
	LastDuration Milliseconds `json:"lastDuration"`

	// Unix timestamp of last execution (in seconds)
	LastExecution Seconds `json:"lastExecution"`

	// Unix timestamp of predicted next execution (in seconds), null if no prediction available
	NextExecution Seconds `json:"nextExecution"`

	// Job type
	Type JobType `json:"type"`

	// Job timeout in seconds
	RequestTimeout Seconds `json:"requestTimeout"`

	// Job schedule
	Schedule Schedule `json:"schedule"`

	// HTTP request method
	RequestMethod RequestMethod `json:"requestMethod"`
}

// DetailedJob
// https://docs.cron-job.org/rest-api.html#jobauth
type DetailedJob struct {
	Job
	// HTTP authentication settings
	Auth JobAuth `json:"auth"`

	// Notification settings
	Notification JobNotificationSettings `json:"notification"`

	// Extended request data
	ExtendedData JobExtendedData `json:"extendedData"`
}

// JobAuth
// https://docs.cron-job.org/rest-api.html#jobauth
type JobAuth struct {
	// Whether to enable HTTP basic authentication or not.
	Enable bool `json:"enable"`

	// HTTP basic auth username
	User string `json:"user"`

	// HTTP basic auth password
	Password string `json:"password"`
}

// JobNotificationSettings
// https://docs.cron-job.org/rest-api.html#job
type JobNotificationSettings struct {
	// Whether to send a notification on job failure or not.
	OnFailure bool `json:"onFailure"`

	// Whether to send a notification when the job succeeds after a prior failure or not.
	OnSuccess bool `json:"onSuccess"`

	// Whether to send a notification when the job has been disabled automatically or not.
	OnDisable bool `json:"onDisable"`
}

// JobExtendedData
// https://docs.cron-job.org/rest-api.html#jobextendeddata
type JobExtendedData struct {
	// Request headers (key-value dictionary)
	Headers map[string]string `json:"headers"`

	// Request body data
	Body string `json:"body"`
}

// JobStatus
// https://docs.cron-job.org/rest-api.html#jobstatus
type JobStatus int

const (
	// Unknown / not executed yet
	JobStatusUnknown = iota

	// OK
	JobStatusOK

	// Failed (DNS error)
	JobStatusFailedDNS

	// Failed (could not connect to host)
	JobStatusFailedCouldNotConnect

	// Failed (HTTP error)
	JobStatusFailedHTTPError

	// Failed (timeout)
	JobStatusFailedTimeout

	// Failed (too much response data)
	JobStatusFailedTooMuchResponse

	// Failed (invalid URL)
	JobStatusFailedInvalidURL

	// Failed (internal errors)
	JobStatusFailedInternalError

	// Failed (unknown reason)
	JobStatusFailedUknown
)

// JobType
// https://docs.cron-job.org/rest-api.html#jobtype
type JobType int

const (
	// Default job
	JobTypeDefault = iota

	// Monitoring job (used in a status monitor)
	JobTypeMonitoring
)

// Schedule
// https://docs.cron-job.org/rest-api.html#jobschedule
type Schedule struct {
	// Schedule time zone (see here for a list of supported values)
	Timezone string `json:"timezone"`

	// Hours in which to execute the job (0-23; [-1] = every hour)
	Hours []int `json:"hours"`

	// Days of month in which to execute the job (1-31; [-1] = every day of month)
	Mdays []int `json:"mdays"`

	// Minutes in which to execute the job (0-59; [-1] = every minute)
	Minutes []int `json:"minutes"`

	// Months in which to execute the job (1-12; [-1] = every month)
	Month []int `json:"months"`

	// Days of week in which to execute the job (0-6; [-1] = every day of week)
	WDays []int `json:"wdays"`
}

// MarshalJSON implements the json.Marshaler interface
func (s Schedule) MarshalJSON() ([]byte, error) {
	if s.Timezone == "" {
		s.Timezone = "Europe/Paris"
	}
	if len(s.Hours) == 0 {
		s.Hours = []int{-1}
	}
	if len(s.Mdays) == 0 {
		s.Mdays = []int{-1}
	}
	if len(s.Minutes) == 0 {
		s.Minutes = []int{-1}
	}
	if len(s.Month) == 0 {
		s.Month = []int{-1}
	}
	if len(s.WDays) == 0 {
		s.WDays = []int{-1}
	}

	type alias Schedule
	return json.Marshal(alias(s))
}

// RequestMethod
// https://docs.cron-job.org/rest-api.html#requestmethod
type RequestMethod int

const (
	RequestMethodGet = iota
	RequestMethodPost
	RequestMethodOptions
	RequestMethodHead
	RequestMethodPut
	RequestMethodDelete
	RequestMethodTrace
	RequestMethodConnect
	RequestMethodPatch
)

// HistoryItem
// https://docs.cron-job.org/rest-api.html#historyitem
type HistoryItem struct {
	// Identifier of the associated cron job
	JobID int `json:"jobId"`

	// Identifier of the history item
	Identifier string `json:"identifier"`

	// Unix timestamp (in seconds) of the actual execution
	Date Seconds `json:"date"`

	// Unix timestamp (in seconds) of the planned/ideal execution
	DatePlanned Seconds `json:"datePlanned"`

	// Scheduling jitter in milliseconds
	Jitter Milliseconds `json:"jitter"`

	// Job URL at time of execution
	URL string `json:"url"`

	// Actual job duration in milliseconds
	Duration Milliseconds `json:"duration"`

	// Status of execution
	Status JobStatus `json:"status"`

	// Detailed job status Description
	StatusText string `json:"statusText"`

	// HTTP status code returned by the host, if any
	HTTPStatus int `json:"httpStatus"`

	// Raw response headers returned by the host (null if unavailable)
	Headers string `json:"headers,omitempty"`

	// Raw response body returned by the host (null if unavailable)
	Body string `json:"body"`

	// Additional timing information for this request
	Stats HistoryItemStats `json:"stats"`
}

// HistoryItemStats
// https://docs.cron-job.org/rest-api.html#historyitemstats
type HistoryItemStats struct {
	// Time from transfer start until name lookups completed (in microseconds)
	NameLookup Microseconds `json:"nameLookup"`

	// Time from transfer start until socket connect completed (in microseconds)
	Connect Microseconds `json:"connect"`

	// Time from transfer start until SSL handshake completed (n microseconds) - 0 if not using SSL
	AppConnect Microseconds `json:"appConnect"`

	// Time from transfer start until beginning of data transfer (in microseconds)
	PreTransfer Microseconds `json:"preTransfer"`

	// Time from transfer start until the first response byte is received (in microseconds)
	StartTransfer Microseconds `json:"startTransfer"`

	// Total transfer time (in microseconds)
	Total Microseconds `json:"total"`
}

// JobHistory is a cron-job.org job history
type JobHistory struct {
	History     HistoryItem `json:"history"`
	Predictions []Seconds   `json:"predictions"`
}
