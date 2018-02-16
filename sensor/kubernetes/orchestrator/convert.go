package orchestrator

import (
	"regexp"
	"strings"

	"bitbucket.org/stack-rox/apollo/pkg/env"
	"bitbucket.org/stack-rox/apollo/pkg/orchestrators"
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	namespaceLabel    = `com.docker.stack.namespace`
	preventLabelValue = `prevent`
	serviceLabel      = `com.prevent.service-name`
)

var (
	invalidDNSLabelCharacter = regexp.MustCompile("[^A-Za-z0-9-]")
)

type serviceWrap struct {
	orchestrators.SystemService
	namespace string
}

type converter struct {
	serviceAccount   string
	imagePullSecrets []string
}

func newConverter() converter {
	ips := env.ImagePullSecrets.Setting()

	return converter{
		serviceAccount:   env.ServiceAccount.Setting(),
		imagePullSecrets: strings.Split(ips, ","),
	}
}

func (c converter) asDaemonSet(service *serviceWrap) *v1beta1.DaemonSet {
	return &v1beta1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "DaemonSet",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.Name,
			Namespace: service.namespace,
		},
		Spec: v1beta1.DaemonSetSpec{
			Template: c.asKubernetesPod(service),
		},
	}
}

func (c converter) asDeployment(service *serviceWrap) *v1beta1.Deployment {
	return &v1beta1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.Name,
			Namespace: service.namespace,
		},
		Spec: v1beta1.DeploymentSpec{
			Replicas: &[]int32{1}[0],
			Template: c.asKubernetesPod(service),
		},
	}
}

func (c converter) asKubernetesPod(service *serviceWrap) v1.PodTemplateSpec {
	return v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: service.namespace,
			Labels:    c.asKubernetesLabels(service.Name),
		},
		Spec: v1.PodSpec{
			Containers:         c.asContainers(service),
			ServiceAccountName: c.serviceAccount,
			ImagePullSecrets:   c.asImagePullSecrets(),
			RestartPolicy:      v1.RestartPolicyAlways,
			Volumes:            c.asVolumes(service),
			HostPID:            service.HostPID,
		},
	}
}

func (converter) asKubernetesLabels(name string) (labels map[string]string) {
	labels = make(map[string]string)

	labels[namespaceLabel] = preventLabelValue
	labels[serviceLabel] = name
	return
}

func (c converter) asImagePullSecrets() (result []v1.LocalObjectReference) {
	result = make([]v1.LocalObjectReference, len(c.imagePullSecrets))

	for i, s := range c.imagePullSecrets {
		result[i] = v1.LocalObjectReference{
			Name: s,
		}
	}
	return
}

func (c converter) asContainers(service *serviceWrap) []v1.Container {
	return []v1.Container{
		{
			Name:         service.Name,
			Env:          c.asEnv(service.Envs),
			Image:        service.Image,
			Command:      service.Command,
			VolumeMounts: c.asVolumeMounts(service),
		},
	}
}

func (c converter) asEnv(envs []string) (vars []v1.EnvVar) {
	for _, env := range envs {
		split := strings.SplitN(env, "=", 2)
		if len(split) == 2 {
			vars = append(vars, v1.EnvVar{
				Name:  split[0],
				Value: split[1],
			})
		}
	}

	return
}
