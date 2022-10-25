package constants

const (
	VolumeName                          = "init-sw-by-platform"
	VolumeMountPath                     = "/opt/ApmAgent"
	JavaToolOptions                     = "-javaagent:/opt/ApmAgent/skywalking-agent/skywalking-agent.jar"
	InitContainerName                   = "init-skywalking-agent-by-platform"
	AdmissionWebhookAnnotationMutateKey = "enable-sw-agent.neo.com/mutate"
	SWBackendSvcKey                     = "SW_AGENT_COLLECTOR_BACKEND_SERVICES"
)
