package constants

const (
	VolumeName                          = "init-jacoco-by-platform"
	VolumeMountPath                     = "/opt/JacocoAgent"
	JavaToolOptions                     = "-javaagent:/opt/JacocoAgent/jacocoagent.jar=includes=*,output=tcpserver,port=6300,address=0.0.0.0,append=true"
	InitContainerName                   = "init-jacoco-agent-by-platform"
	AdmissionWebhookAnnotationMutateKey = "enable-jacoco-agent.neo.com/mutate"
)
