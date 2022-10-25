package constants

const (
	VolumeName                          = "init-sw-by-platform"
	VolumeMountPath                     = "/opt/ApmAgent"
	JavaToolOptions                     = "-javaagent:/opt/ApmAgent/skywalking-agent/skywalking-agent.jar"
	InitContainerName                   = "init-skywalking-agent-by-platform"
	AdmissionWebhookAnnotationMutateKey = "enable-sw.unicloud.com/mutate"
)
