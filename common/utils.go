package common

import (
	"bytes"
	"fmt"
	"initJacocoAgent/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"strings"
)

var InitJacocoAgentImg = os.Getenv("INIT_JACOCO_AGENT_IMG")

func RequiredMutate(ignoredList []string, metadata *metav1.ObjectMeta) (required bool) {
	var log *bytes.Buffer
	for _, ns := range ignoredList {
		if metadata.Namespace == ns {
			log.WriteString(fmt.Sprintf("\nSkip validate for [#{metadata.Name}]"))
			required = false
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
		required = false
	}
	return
}
