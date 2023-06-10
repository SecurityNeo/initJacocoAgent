package service

import (
	"encoding/json"
	"fmt"
	"github.com/wI2L/jsondiff"
	"initJacocoAgent/common"
	adminssionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MutateSts(sts *appsv1.StatefulSet, protectNS []string) *adminssionv1.AdmissionResponse {
	objectMeta := &sts.ObjectMeta
	if !common.RequiredMutate(protectNS, objectMeta) {
		return &adminssionv1.AdmissionResponse{
			Allowed: true,
		}
	}

	newSts := sts.DeepCopy()
	newPodSpec := mutatePodSpec(newSts.Spec.Template.Spec)
	newSts.Spec.Template.Spec = newPodSpec
	patch, err := jsondiff.Compare(sts, newSts)
	if err != nil {
		return &adminssionv1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	patchBytes, err := json.MarshalIndent(patch, "", "	")
	if err != nil {
		return &adminssionv1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}
	fmt.Printf("\nAdmissionResponse: patch=%v\n", string(patchBytes))
	return &adminssionv1.AdmissionResponse{
		Allowed: true,
		Patch:   patchBytes,
		PatchType: func() *adminssionv1.PatchType {
			pt := adminssionv1.PatchTypeJSONPatch
			return &pt
		}(),
	}
}
