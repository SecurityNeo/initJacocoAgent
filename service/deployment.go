package service

import (
	"encoding/json"
	"fmt"
	"github.com/wI2L/jsondiff"
	"initJacocoAgent/common"
	adminssionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func MutateDeploy(deploy *appsv1.Deployment) *adminssionv1.AdmissionResponse {
	resourceNS, resourceName, objectMeta := deploy.Namespace, deploy.Name, &deploy.ObjectMeta
	fmt.Printf("\n----PreCheck----")
	if !common.RequiredMutate(ignoredNamespaces, objectMeta) {
		fmt.Printf("\nSkip mutate for %v/%v", resourceNS, resourceName)
		return &adminssionv1.AdmissionResponse{
			Allowed: true,
		}
	}
	fmt.Printf("\n----EndCheck----")

	newDp := deploy.DeepCopy()
	newPodSpec := mutatePodSpec(newDp.Spec.Template.Spec, resourceNS, resourceName)
	newDp.Spec.Template.Spec = newPodSpec
	fmt.Printf("\n----BeginMutateYaml----")
	bytes, err := json.Marshal(newPodSpec)
	if err == nil {
		yamlStr, err := yaml.JSONToYAML(bytes)
		if err == nil {
			fmt.Printf("\n----YamlContent----\n" + string(yamlStr))
		}
	}
	fmt.Printf("\n----EndMutateYaml----")
	patch, err := jsondiff.Compare(deploy, newDp)
	if err != nil {
		fmt.Printf("\nCompare patch error: %v", err)
		return &adminssionv1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	patchBytes, err := json.MarshalIndent(patch, "", "	")
	if err != nil {
		fmt.Printf("\nPatch error: %v", err)
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
