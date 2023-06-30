package report

type jobStatus string

const (
	Pending    jobStatus = "pending"
	InProgress jobStatus = "build in progress"
	Success    jobStatus = "build successful"
	Failed     jobStatus = "build failed"
)
