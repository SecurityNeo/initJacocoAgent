package service

import (
	"fmt"
	"initSkywalkingAgent/common"
	"initSkywalkingAgent/constants"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

var (
	ignoredNamespaces = []string{
		metav1.NamespaceSystem,
		metav1.NamespaceDefault,
		metav1.NamespacePublic,
	}
)

func mutatePodSpec(podSpec corev1.PodSpec, ns string, name string) corev1.PodSpec {
	fmt.Printf("\n----PreMutate----")
	SwBackEnd := os.Getenv(constants.SWBackendSvcKey)
	addVolume := corev1.Volume{
		Name: constants.VolumeName,
	}
	addFlag := true
	for _, v := range podSpec.Volumes {
		if v.Name == constants.VolumeName {
			addFlag = false
		}
	}
	if addFlag {
		podSpec.Volumes = append(podSpec.Volumes, addVolume)
		addVolumeMount := corev1.VolumeMount{
			Name:      constants.VolumeName,
			MountPath: constants.VolumeMountPath,
		}
		var addEnvs = []corev1.EnvVar{
			{
				Name:  constants.SWBackendSvcKey,
				Value: SwBackEnd,
			},
			{
				Name:  "SW_AGENT_NAMESPACE",
				Value: ns,
			},
			{
				Name:  "SW_AGENT_NAME",
				Value: name,
			},
			{
				Name:  "JAVA_TOOL_OPTIONS",
				Value: constants.JavaToolOptions,
			},
		}
		containers := podSpec.Containers
		for i, _ := range containers {
			flag := true
			for _, v := range containers[i].VolumeMounts {
				if v.Name == constants.VolumeName {
					flag = false
				}
			}
			if flag {
				containers[i].VolumeMounts = append(containers[i].VolumeMounts, addVolumeMount)
			}
			for _, addEnv := range addEnvs {
				flag := true
				for _, v := range containers[i].Env {
					if v.Name == addEnv.Name {
						flag = false
					}
				}
				if flag {
					containers[i].Env = append(containers[i].Env, addEnv)
				}
			}
		}
		podSpec.Containers = containers
	}

	addInitContainerFlag := true
	for _, v := range podSpec.InitContainers {
		if v.Name == constants.InitContainerName {
			addInitContainerFlag = false
		}
	}
	if addInitContainerFlag {
		initContainer := corev1.Container{
			Name:    constants.InitContainerName,
			Image:   common.InitSwAgentImg,
			Command: []string{"/bin/sh"},
			Args:    []string{"-c", "cp -r /opt/skywalking-agent " + constants.VolumeMountPath},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      constants.VolumeName,
					MountPath: constants.VolumeMountPath,
				},
			},
		}
		initContainerReq := make(map[corev1.ResourceName]resource.Quantity)
		initContainerReq[corev1.ResourceCPU] = *resource.NewMilliQuantity(100, resource.DecimalSI)
		initContainerReq[corev1.ResourceMemory] = *resource.NewQuantity(100*1024*1024, resource.BinarySI)
		initContainer.Resources.Requests = initContainerReq
		podSpec.InitContainers = append(podSpec.InitContainers, initContainer)
	}
	fmt.Printf("\n----EndMutate----")
	return podSpec
}
