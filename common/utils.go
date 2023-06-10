package common

import (
	"fmt"
	"initJacocoAgent/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"strings"
)

var InitJacocoAgentImg = os.Getenv("INIT_JACOCO_AGENT_IMG")

func RequiredMutate(ignoredList []string, metadata *metav1.ObjectMeta) (required bool) {

	for _, ns := range ignoredList {
		if metadata.Namespace == ns {
			required = false
			fmt.Printf("Skip mutate for %v,because the namespace %v is protected.\n", metadata.Name, ns)
			return
		}
	}
	annotations := metadata.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}
	switch strings.ToLower(annotations[constants.AdmissionWebhookAnnotationMutateKey]) {
	case "true":
		required = true
	default:
		annotation := constants.AdmissionWebhookAnnotationMutateKey + ": " + "true"
		fmt.Printf("Skip mutate for %v,because the resource has no annotations [%v].\n", metadata.Name, annotation)
		required = false
	}
	return
}
