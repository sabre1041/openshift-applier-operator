package util

import (
	"fmt"
	"regexp"
	"strings"

	copapi "github.com/redhat-cop/openshift-applier-operator/pkg/apis/cop/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"

	batchv1 "k8s.io/api/batch/v1"
)

const (
	scmContainerName     = "applier-git-cloner"
	applierContainerName = "applier"
	applierImage         = "quay.io/redhat-cop/openshift-applier"
	sharedVolumeMount    = "/app"
	sharedVolumeName     = "repo"
	masterBranch         = "master"
	applierImageHome     = "/openshift-applier"
	sshVolumeName        = "ssh"
	inventoryPath        = "INVENTORY_PATH"
	repoURL              = "REPO_URL"
	repoRef              = "REPO_REF"
	destination          = "DESTINATION"
	prepSSHCommand       = "mkdir -p ~/.ssh && cp --no-preserve=mode /secret/id_rsa ~/.ssh/ && chmod 600 ~/.ssh/id_rsa && echo -e \"Host *\n   StrictHostKeyChecking no\" > ~/.ssh/config && "
)

var (
	initContainerCommand = "git clone -b $REPO_REF $REPO_URL $DESTINATION"
)

func GenerateJob(applier *copapi.Applier) (*batchv1.Job, error) {

	re := regexp.MustCompile("[^A-Za-z0-9]")
	jobName := strings.ToLower(fmt.Sprintf("%s-%s", re.ReplaceAllString(applier.Name, "-"), rand.String(4)))

	var repoBranch string

	if applier.Spec.Source.Git.Ref != "" {
		repoBranch = applier.Spec.Source.Git.Ref
	} else {
		repoBranch = masterBranch
	}

	job := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: applier.Namespace,
		},
	}

	job.Spec.Template.Spec.RestartPolicy = corev1.RestartPolicyNever

	jobVolumes := []corev1.Volume{
		{
			Name: sharedVolumeName,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
	}

	initContainerVolumeMounts := []corev1.VolumeMount{
		{
			Name:      sharedVolumeName,
			MountPath: sharedVolumeMount,
		},
	}

	if applier.Spec.Source.Git.SecretName != "" {
		initContainerVolumeMounts = append(initContainerVolumeMounts, corev1.VolumeMount{
			Name:      sshVolumeName,
			MountPath: "/secret",
		})

		jobVolumes = append(jobVolumes, corev1.Volume{
			Name: sshVolumeName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: applier.Spec.Source.Git.SecretName,
				},
			},
		})

		initContainerCommand = prepSSHCommand + initContainerCommand
	}

	initContainer := corev1.Container{
		Name:  scmContainerName,
		Image: applierImage,
		Args:  []string{"/bin/bash", "-c", initContainerCommand},
		Env: []corev1.EnvVar{
			{
				Name:  repoURL,
				Value: applier.Spec.Source.Git.URI,
			},
			{
				Name:  repoRef,
				Value: repoBranch,
			},
			{
				Name:  destination,
				Value: sharedVolumeMount,
			},
		},
	}

	initContainer.VolumeMounts = initContainerVolumeMounts

	job.Spec.Template.Spec.InitContainers = []corev1.Container{initContainer}

	container := corev1.Container{
		Name:  applierContainerName,
		Image: applierImage,
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      sharedVolumeName,
				MountPath: sharedVolumeMount,
			},
		},
	}

	containerEnv := []corev1.EnvVar{}

	var inventoryLocation string

	if applier.Spec.Source.Git.InventoryDir != "" {
		inventoryLocation = sharedVolumeMount + "/" + applier.Spec.Source.Git.InventoryDir
	} else {
		inventoryLocation = sharedVolumeMount
	}

	containerEnv = append(containerEnv, corev1.EnvVar{
		Name:  inventoryPath,
		Value: inventoryLocation,
	})

	container.Env = containerEnv

	job.Spec.Template.Spec.Containers = []corev1.Container{container}

	job.Spec.Template.Spec.Volumes = jobVolumes

	if applier.Spec.ServiceAccount != "" {
		job.Spec.Template.Spec.ServiceAccountName = applier.Spec.ServiceAccount
	}

	return job, nil
}

func ParseQueryString(querystring string) (string, string, error) {

	var finalArray []string

	for _, value := range strings.Split(querystring, "/") {
		if len(strings.TrimSpace(value)) > 0 {
			finalArray = append(finalArray, value)
		}
	}

	if len(finalArray) < 2 {
		return "", "", fmt.Errorf("")
	}

	return finalArray[len(finalArray)-2], finalArray[len(finalArray)-1], nil
}

func IsErrorMessage(err error, message string) bool {
	if err.Error() == message {
		return true
	}

	return false
}
