package booleanpolicy

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	gogoTypes "github.com/gogo/protobuf/types"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/image/policies"
	"github.com/stackrox/rox/pkg/defaults"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/fixtures"
	"github.com/stackrox/rox/pkg/images/types"
	policyUtils "github.com/stackrox/rox/pkg/policies"
	"github.com/stackrox/rox/pkg/protoconv"
	"github.com/stackrox/rox/pkg/readable"
	"github.com/stackrox/rox/pkg/set"
	"github.com/stackrox/rox/pkg/sliceutils"
	"github.com/stackrox/rox/pkg/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestDefaultPolicies(t *testing.T) {
	suite.Run(t, new(DefaultPoliciesTestSuite))
}

type DefaultPoliciesTestSuite struct {
	suite.Suite

	defaultPolicies map[string]*storage.Policy

	deployments             map[string]*storage.Deployment
	images                  map[string]*storage.Image
	deploymentsToImages     map[string][]*storage.Image
	deploymentsToIndicators map[string][]*storage.ProcessIndicator
}

func (suite *DefaultPoliciesTestSuite) SetupSuite() {
	suite.deployments = make(map[string]*storage.Deployment)
	suite.images = make(map[string]*storage.Image)
	suite.deploymentsToImages = make(map[string][]*storage.Image)
	suite.deploymentsToIndicators = make(map[string][]*storage.ProcessIndicator)

	defaults.PoliciesPath = policies.Directory()

	defaultPolicies, err := defaults.Policies()
	suite.Require().NoError(err)

	suite.defaultPolicies = make(map[string]*storage.Policy, len(defaultPolicies))
	for _, p := range defaultPolicies {
		suite.defaultPolicies[p.GetName()] = p
	}
}

func (suite *DefaultPoliciesTestSuite) imageIDFromDep(deployment *storage.Deployment) string {
	suite.Require().Len(deployment.GetContainers(), 1, "This function only supports deployments with exactly one container")
	id := deployment.GetContainers()[0].GetImage().GetId()
	suite.NotEmpty(id, "Deployment '%s' had no image id", proto.MarshalTextString(deployment))
	return id
}

func (suite *DefaultPoliciesTestSuite) TestNoDuplicatePolicyIDs() {
	ids := set.NewStringSet()
	for _, p := range suite.defaultPolicies {
		suite.True(ids.Add(p.GetId()))
	}
}

func (suite *DefaultPoliciesTestSuite) MustGetPolicy(name string) *storage.Policy {
	p, ok := suite.defaultPolicies[name]
	suite.Require().True(ok, "Policy %s not found", name)
	return p
}

func (suite *DefaultPoliciesTestSuite) addDepAndImages(deployment *storage.Deployment, images ...*storage.Image) {
	suite.deployments[deployment.GetId()] = deployment
	for _, i := range images {
		suite.images[i.GetId()] = i
		suite.deploymentsToImages[deployment.GetId()] = append(suite.deploymentsToImages[deployment.GetId()], i)
	}
}

func imageWithComponents(components []*storage.EmbeddedImageScanComponent) *storage.Image {
	return &storage.Image{
		Id:   uuid.NewV4().String(),
		Name: &storage.ImageName{FullName: "ASFASF"},
		Scan: &storage.ImageScan{
			Components: components,
		},
	}
}

func imageWithLayers(layers []*storage.ImageLayer) *storage.Image {
	return &storage.Image{
		Id: uuid.NewV4().String(),
		Metadata: &storage.ImageMetadata{
			V1: &storage.V1Metadata{
				Layers: layers,
			},
		},
	}
}

func deploymentWithImageAnyID(img *storage.Image) *storage.Deployment {
	return deploymentWithImage(uuid.NewV4().String(), img)
}

func deploymentWithImage(id string, img *storage.Image) *storage.Deployment {
	return &storage.Deployment{
		Id:         id,
		Containers: []*storage.Container{{Image: types.ToContainerImage(img)}},
	}
}

func (suite *DefaultPoliciesTestSuite) addIndicator(deploymentID, name, args, path string, lineage []string, uid uint32) *storage.ProcessIndicator {
	deployment := suite.deployments[deploymentID]
	if len(deployment.GetContainers()) == 0 {
		deployment.Containers = []*storage.Container{{Name: uuid.NewV4().String()}}
	}
	indicator := &storage.ProcessIndicator{
		Id:            uuid.NewV4().String(),
		DeploymentId:  deploymentID,
		ContainerName: deployment.GetContainers()[0].GetName(),
		Signal: &storage.ProcessSignal{
			Name:         name,
			Args:         args,
			ExecFilePath: path,
			Time:         gogoTypes.TimestampNow(),
			Lineage:      lineage,
			Uid:          uid,
		},
	}
	suite.deploymentsToIndicators[deploymentID] = append(suite.deploymentsToIndicators[deploymentID], indicator)
	return indicator
}

type testCase struct {
	policyName                string
	skip                      bool
	expectedViolations        map[string][]*storage.Alert_Violation
	expectedProcessViolations map[string][]*storage.ProcessIndicator

	// If shouldNotMatch is specified (which is the case for policies that check for the absence of something), we verify that
	// it matches everything except shouldNotMatch.
	// If sampleViolationForMatched is provided, we verify that all the matches are the string provided in sampleViolationForMatched.
	shouldNotMatch            map[string]struct{}
	sampleViolationForMatched string
}

func (suite *DefaultPoliciesTestSuite) getImagesForDeployment(deployment *storage.Deployment) []*storage.Image {
	images := suite.deploymentsToImages[deployment.GetId()]
	if len(images) == 0 {
		return make([]*storage.Image, len(deployment.GetContainers()))
	}
	suite.Equal(len(deployment.GetContainers()), len(images))
	return images
}

func (suite *DefaultPoliciesTestSuite) TestDefaultPolicies() {
	fixtureDep := fixtures.GetDeployment()
	fixturesImages := fixtures.DeploymentImages()

	suite.addDepAndImages(fixtureDep, fixturesImages...)

	nginx110 := &storage.Image{
		Id: "SHANGINX110",
		Name: &storage.ImageName{
			Registry: "docker.io",
			Remote:   "library/nginx",
			Tag:      "1.10",
		},
	}

	nginx110Dep := deploymentWithImage("nginx110", nginx110)
	suite.addDepAndImages(nginx110Dep, nginx110)

	oldScannedTime := time.Now().Add(-31 * 24 * time.Hour)
	oldScannedImage := &storage.Image{
		Id: "SHAOLDSCANNED",
		Scan: &storage.ImageScan{
			ScanTime: protoconv.ConvertTimeToTimestamp(oldScannedTime),
		},
	}
	oldScannedDep := deploymentWithImage("oldscanned", oldScannedImage)
	suite.addDepAndImages(oldScannedDep, oldScannedImage)

	addDockerFileImg := imageWithLayers([]*storage.ImageLayer{
		{
			Instruction: "ADD",
			Value:       "deploy.sh",
		},
		{
			Instruction: "RUN",
			Value:       "deploy.sh",
		},
	})
	addDockerFileDep := deploymentWithImageAnyID(addDockerFileImg)
	suite.addDepAndImages(addDockerFileDep, addDockerFileImg)

	imagePort22Image := imageWithLayers([]*storage.ImageLayer{
		{
			Instruction: "EXPOSE",
			Value:       "22/tcp",
		},
	})
	imagePort22Dep := deploymentWithImageAnyID(imagePort22Image)
	suite.addDepAndImages(imagePort22Dep, imagePort22Image)

	insecureCMDImage := imageWithLayers([]*storage.ImageLayer{
		{
			Instruction: "CMD",
			Value:       "do an insecure thing",
		},
	})

	insecureCMDDep := deploymentWithImageAnyID(insecureCMDImage)
	suite.addDepAndImages(insecureCMDDep, insecureCMDImage)

	runSecretsImage := imageWithLayers([]*storage.ImageLayer{
		{
			Instruction: "VOLUME",
			Value:       "/run/secrets",
		},
	})
	runSecretsDep := deploymentWithImageAnyID(runSecretsImage)
	suite.addDepAndImages(runSecretsDep, runSecretsImage)

	oldImageCreationTime := time.Now().Add(-100 * 24 * time.Hour)
	oldCreatedImage := &storage.Image{
		Id: "SHA:OLDCREATEDIMAGE",
		Metadata: &storage.ImageMetadata{
			V1: &storage.V1Metadata{
				Created: protoconv.ConvertTimeToTimestamp(oldImageCreationTime),
			},
		},
	}
	oldImageDep := deploymentWithImage("oldimagedep", oldCreatedImage)
	suite.addDepAndImages(oldImageDep, oldCreatedImage)

	apkImage := imageWithComponents([]*storage.EmbeddedImageScanComponent{
		{Name: "apk", Version: "1.2"},
		{Name: "asfa", Version: "1.5"},
	})
	apkDep := deploymentWithImageAnyID(apkImage)
	suite.addDepAndImages(apkDep, apkImage)

	curlImage := imageWithComponents([]*storage.EmbeddedImageScanComponent{
		{Name: "curl", Version: "1.3"},
		{Name: "curlwithextra", Version: "0.9"},
	})
	curlDep := deploymentWithImageAnyID(curlImage)
	suite.addDepAndImages(curlDep, curlImage)

	componentDeps := make(map[string]*storage.Deployment)
	for _, component := range []string{"apt", "dnf", "wget"} {
		img := imageWithComponents([]*storage.EmbeddedImageScanComponent{
			{Name: component},
		})
		dep := deploymentWithImageAnyID(img)
		suite.addDepAndImages(dep, img)
		componentDeps[component] = dep
	}

	heartbleedDep := &storage.Deployment{
		Id: "HEARTBLEEDDEPID",
		Containers: []*storage.Container{
			{
				SecurityContext: &storage.SecurityContext{Privileged: true},
				Image:           &storage.ContainerImage{Id: "HEARTBLEEDDEPSHA"},
			},
		},
	}
	suite.addDepAndImages(heartbleedDep, &storage.Image{
		Id: "HEARTBLEEDDEPSHA",
		Scan: &storage.ImageScan{
			Components: []*storage.EmbeddedImageScanComponent{
				{Name: "heartbleed", Version: "1.2", Vulns: []*storage.EmbeddedVulnerability{
					{Cve: "CVE-2014-0160", Link: "https://heartbleed", Cvss: 6, SetFixedBy: &storage.EmbeddedVulnerability_FixedBy{FixedBy: "v1.2"}},
				}},
			},
		},
	})

	requiredImageLabel := &storage.Deployment{
		Id: "requiredImageLabel",
		Containers: []*storage.Container{
			{
				Image: &storage.ContainerImage{Id: "requiredImageLabelImage"},
			},
		},
	}
	suite.addDepAndImages(requiredImageLabel, &storage.Image{
		Id: "requiredImageLabelImage",
		Metadata: &storage.ImageMetadata{
			V1: &storage.V1Metadata{
				Labels: map[string]string{
					"required-label": "required-value",
				},
			},
		},
	})

	shellshockImage := imageWithComponents([]*storage.EmbeddedImageScanComponent{
		{Name: "shellshock", Version: "1.2", Vulns: []*storage.EmbeddedVulnerability{
			{Cve: "CVE-2014-6271", Link: "https://shellshock", Cvss: 6},
			{Cve: "CVE-ARBITRARY", Link: "https://notshellshock"},
		}},
	})
	shellshockDep := deploymentWithImageAnyID(shellshockImage)
	suite.addDepAndImages(shellshockDep, shellshockImage)

	strutsImage := imageWithComponents([]*storage.EmbeddedImageScanComponent{
		{Name: "struts", Version: "1.2", Vulns: []*storage.EmbeddedVulnerability{
			{Cve: "CVE-2017-5638", Link: "https://struts", Cvss: 8, SetFixedBy: &storage.EmbeddedVulnerability_FixedBy{FixedBy: "v1.3"}},
		}},
		{Name: "OTHER", Version: "1.3", Vulns: []*storage.EmbeddedVulnerability{
			{Cve: "CVE-1223-451", Link: "https://cvefake"},
		}},
	})
	strutsDep := deploymentWithImageAnyID(strutsImage)
	suite.addDepAndImages(strutsDep, strutsImage)

	depWithNonSeriousVulnsImage := imageWithComponents([]*storage.EmbeddedImageScanComponent{
		{Name: "NOSERIOUS", Version: "2.3", Vulns: []*storage.EmbeddedVulnerability{
			{Cve: "CVE-1234-5678", Link: "https://abcdefgh"},
			{Cve: "CVE-5678-1234", Link: "https://lmnopqrst"},
		}},
	})
	depWithNonSeriousVulns := deploymentWithImageAnyID(depWithNonSeriousVulnsImage)
	suite.addDepAndImages(depWithNonSeriousVulns, depWithNonSeriousVulnsImage)

	dockerSockDep := &storage.Deployment{
		Id: "DOCKERSOCDEP",
		Containers: []*storage.Container{
			{Volumes: []*storage.Volume{
				{Source: "/var/run/docker.sock", Name: "DOCKERSOCK"},
				{Source: "NOTDOCKERSOCK"},
			}},
		},
	}
	suite.addDepAndImages(dockerSockDep)

	containerPort22Dep := &storage.Deployment{
		Id: "CONTAINERPORT22DEP",
		Ports: []*storage.PortConfig{
			{Protocol: "TCP", ContainerPort: 22},
			{Protocol: "UDP", ContainerPort: 4125},
		},
	}
	suite.addDepAndImages(containerPort22Dep)

	secretEnvDep := &storage.Deployment{
		Id: "SECRETENVDEP",
		Containers: []*storage.Container{
			{Config: &storage.ContainerConfig{
				Env: []*storage.ContainerConfig_EnvironmentConfig{
					{Key: "THIS_IS_SECRET_VAR", Value: "stealthmode", EnvVarSource: storage.ContainerConfig_EnvironmentConfig_RAW},
					{Key: "HOME", Value: "/home/stackrox"},
				},
			}},
		},
	}
	suite.addDepAndImages(secretEnvDep)

	secretEnvSrcUnsetDep := &storage.Deployment{
		Id: "SECRETENVSRCUNSETDEP",
		Containers: []*storage.Container{
			{Config: &storage.ContainerConfig{
				Env: []*storage.ContainerConfig_EnvironmentConfig{
					{Key: "THIS_IS_SECRET_VAR", Value: "stealthmode"},
				},
			}},
		},
	}
	suite.addDepAndImages(secretEnvSrcUnsetDep)

	secretKeyRefDep := &storage.Deployment{
		Id: "SECRETKEYREFDEP",
		Containers: []*storage.Container{
			{Config: &storage.ContainerConfig{
				Env: []*storage.ContainerConfig_EnvironmentConfig{
					{Key: "THIS_IS_SECRET_VAR", EnvVarSource: storage.ContainerConfig_EnvironmentConfig_SECRET_KEY},
					{Key: "HOME", Value: "/home/stackrox"},
				},
			}},
		},
	}
	suite.addDepAndImages(secretKeyRefDep)

	// Fake deployment that shouldn't match anything, just to make sure
	// that none of our queries will accidentally match it.
	suite.addDepAndImages(&storage.Deployment{Id: "FAKEID", Name: "FAKENAME"})

	depWithGoodEmailAnnotation := &storage.Deployment{
		Id: "GOODEMAILDEPID",
		Annotations: map[string]string{
			"email": "vv@stackrox.com",
		},
	}
	suite.addDepAndImages(depWithGoodEmailAnnotation)

	depWithOwnerAnnotation := &storage.Deployment{
		Id: "OWNERANNOTATIONDEP",
		Annotations: map[string]string{
			"owner": "IOWNTHIS",
			"blah":  "Blah",
		},
	}
	suite.addDepAndImages(depWithOwnerAnnotation)

	depWitharbitraryAnnotations := &storage.Deployment{
		Id: "ARBITRARYANNOTATIONDEPID",
		Annotations: map[string]string{
			"emailnot": "vv@stackrox.com",
			"notemail": "vv@stackrox.com",
			"ownernot": "vv",
			"nowner":   "vv",
		},
	}
	suite.addDepAndImages(depWitharbitraryAnnotations)

	depWithBadEmailAnnotation := &storage.Deployment{
		Id: "BADEMAILDEPID",
		Annotations: map[string]string{
			"email": "NOTANEMAIL",
		},
	}
	suite.addDepAndImages(depWithBadEmailAnnotation)

	sysAdminDep := &storage.Deployment{
		Id: "SYSADMINDEPID",
		Containers: []*storage.Container{
			{
				SecurityContext: &storage.SecurityContext{
					AddCapabilities: []string{"CAP_SYS_ADMIN"},
				},
			},
		},
	}
	suite.addDepAndImages(sysAdminDep)

	depWithAllResourceLimitsRequestsSpecified := &storage.Deployment{
		Id: "ALLRESOURCESANDLIMITSDEP",
		Containers: []*storage.Container{
			{Resources: &storage.Resources{
				CpuCoresRequest: 0.1,
				CpuCoresLimit:   0.3,
				MemoryMbLimit:   100,
				MemoryMbRequest: 1251,
			}},
		},
	}
	suite.addDepAndImages(depWithAllResourceLimitsRequestsSpecified)

	depWithEnforcementBypassAnnotation := &storage.Deployment{
		Id: "ENFORCEMENTBYPASS",
		Annotations: map[string]string{
			"admission.stackrox.io/break-glass": "ticket-1234",
		},
	}
	suite.addDepAndImages(depWithEnforcementBypassAnnotation)

	hostMountDep := &storage.Deployment{
		Id: "HOSTMOUNT",
		Containers: []*storage.Container{
			{Volumes: []*storage.Volume{
				{Source: "/etc/passwd", Name: "HOSTMOUNT"},
				{Source: "/var/lib/kubelet", Name: "KUBELET"},
			}},
		},
	}
	suite.addDepAndImages(hostMountDep)

	// Index processes
	bashLineage := []string{"/bin/bash"}
	fixtureDepAptIndicator := suite.addIndicator(fixtureDep.GetId(), "apt", "", "/usr/bin/apt", bashLineage, 1)
	sysAdminDepAptIndicator := suite.addIndicator(sysAdminDep.GetId(), "apt", "install blah", "/usr/bin/apt", bashLineage, 1)

	kubeletIndicator := suite.addIndicator(containerPort22Dep.GetId(), "curl", "https://12.13.14.15:10250", "/bin/curl", bashLineage, 1)
	kubeletIndicator2 := suite.addIndicator(containerPort22Dep.GetId(), "wget", "https://heapster.kube-system/metrics", "/bin/wget", bashLineage, 1)

	nmapIndicatorfixtureDep1 := suite.addIndicator(fixtureDep.GetId(), "nmap", "blah", "/usr/bin/nmap", bashLineage, 1)
	nmapIndicatorfixtureDep2 := suite.addIndicator(fixtureDep.GetId(), "nmap", "blah2", "/usr/bin/nmap", bashLineage, 1)
	nmapIndicatorNginx110Dep := suite.addIndicator(nginx110Dep.GetId(), "nmap", "", "/usr/bin/nmap", bashLineage, 1)

	javaLineage := []string{"/bin/bash", "/mnt/scripts/run_server.sh", "/bin/java"}
	fixtureDepJavaIndicator := suite.addIndicator(fixtureDep.GetId(), "/bin/bash", "-attack", "/bin/bash", javaLineage, 0)

	deploymentTestCases := []testCase{
		{
			policyName: "Latest tag",
			expectedViolations: map[string][]*storage.Alert_Violation{
				fixtureDep.GetId(): {
					{
						Message: "Image tag 'latest' matched latest",
					},
				},
			},
		},
		{
			policyName: "DockerHub NGINX 1.10",
			expectedViolations: map[string][]*storage.Alert_Violation{
				fixtureDep.GetId(): {
					{
						Message: "Image tag '1.10' matched 1.10",
					},
					{
						Message: "Image registry 'docker.io' matched docker.io",
					},
					{
						Message: "Image remote 'library/nginx' matched nginx",
					},
				},
				nginx110Dep.GetId(): {
					{
						Message: "Image tag '1.10' matched 1.10",
					},
					{
						Message: "Image registry 'docker.io' matched docker.io",
					},
					{
						Message: "Image remote 'library/nginx' matched nginx",
					},
				},
			},
		},
		{
			policyName: "Alpine Linux Package Manager (apk) in Image",
			expectedViolations: map[string][]*storage.Alert_Violation{
				apkDep.GetId(): {
					{
						Message: "Component name 'apk' matched apk",
					},
				},
			},
		},
		{
			policyName: "Ubuntu Package Manager in Image",
			expectedViolations: map[string][]*storage.Alert_Violation{
				componentDeps["apt"].GetId(): {
					{
						Message: "Component name 'apt' matched apt|dpkg",
					},
				},
			},
		},
		{
			policyName: "Curl in Image",
			expectedViolations: map[string][]*storage.Alert_Violation{
				curlDep.GetId(): {
					{
						Message: "Component name 'curl' matched curl",
					},
				},
			},
		},
		{
			policyName: "Red Hat Package Manager in Image",
			expectedViolations: map[string][]*storage.Alert_Violation{
				componentDeps["dnf"].GetId(): {
					{
						Message: "Component name 'dnf' matched rpm|dnf|yum",
					},
				},
			},
		},
		{
			policyName: "Wget in Image",
			expectedViolations: map[string][]*storage.Alert_Violation{
				componentDeps["wget"].GetId(): {
					{
						Message: "Component name 'wget' matched wget",
					},
				},
			},
		},
		{
			policyName: "Mount Docker Socket",
			expectedViolations: map[string][]*storage.Alert_Violation{
				dockerSockDep.GetId(): {
					{
						Message: "Volume source '/var/run/docker.sock' matched /var/run/docker.sock",
					},
				},
			},
		},
		{
			policyName: "90-Day Image Age",
			skip:       true,
			expectedViolations: map[string][]*storage.Alert_Violation{
				oldImageDep.GetId(): {
					{
						Message: fmt.Sprintf("Time of image creation '%s' was more than 90 days ago", readable.Time(oldImageCreationTime)),
					},
				},
			},
		},
		{
			policyName: "30-Day Scan Age",
			skip:       true,
			expectedViolations: map[string][]*storage.Alert_Violation{
				oldScannedDep.GetId(): {
					{
						Message: fmt.Sprintf("Time of last scan '%s' was more than 30 days ago", readable.Time(oldScannedTime)),
					},
				},
			},
		},
		{
			policyName: "Secure Shell (ssh) Port Exposed in Image",
			expectedViolations: map[string][]*storage.Alert_Violation{
				imagePort22Dep.GetId(): {
					{
						Message: "Dockerfile Line 'EXPOSE 22/tcp' matches the rule EXPOSE (22/tcp|\\s+22/tcp)",
					},
				},
			},
		},
		{
			policyName: "Secure Shell (ssh) Port Exposed",
			expectedViolations: map[string][]*storage.Alert_Violation{
				containerPort22Dep.GetId(): {
					{
						Message: "Port '22' matched 22",
					},
					{
						Message: "Protocol 'tcp' matched tcp",
					},
				},
			},
		},
		{
			policyName: "Privileged Container",
			expectedViolations: map[string][]*storage.Alert_Violation{
				fixtureDep.GetId(): {
					{
						Message: "Privileged container found",
					},
				},
				heartbleedDep.GetId(): {
					{
						Message: "Privileged container found",
					},
				},
			},
		},
		{
			policyName: "Container using read-write root filesystem",
			expectedViolations: map[string][]*storage.Alert_Violation{
				heartbleedDep.GetId(): {
					{
						Message: "Container using read-write root filesystem found",
					},
				},
				fixtureDep.GetId(): {
					{
						Message: "Container using read-write root filesystem found",
					},
				},
				sysAdminDep.GetId(): {
					{
						Message: "Container using read-write root filesystem found",
					},
				},
			},
		},
		{
			policyName: "Insecure specified in CMD",
			expectedViolations: map[string][]*storage.Alert_Violation{
				insecureCMDDep.GetId(): {
					{
						Message: "Dockerfile Line 'CMD do an insecure thing' matches the rule CMD .*insecure.*",
					},
				},
			},
		},
		{
			policyName: "Improper Usage of Orchestrator Secrets Volume",
			expectedViolations: map[string][]*storage.Alert_Violation{
				runSecretsDep.GetId(): {
					{
						Message: "Dockerfile Line 'VOLUME /run/secrets' matches the rule VOLUME /run/secrets",
					},
				},
			},
		},
		{
			policyName: "Images with no scans",
			skip:       true,
			shouldNotMatch: map[string]struct{}{
				// These deployments have scans on their images.
				fixtureDep.GetId():    {},
				oldScannedDep.GetId(): {},
				// The rest of the deployments have no images!
				"FAKEID":                                          {},
				containerPort22Dep.GetId():                        {},
				dockerSockDep.GetId():                             {},
				secretEnvDep.GetId():                              {},
				secretEnvSrcUnsetDep.GetId():                      {},
				secretKeyRefDep.GetId():                           {},
				depWithOwnerAnnotation.GetId():                    {},
				depWithGoodEmailAnnotation.GetId():                {},
				depWithBadEmailAnnotation.GetId():                 {},
				depWitharbitraryAnnotations.GetId():               {},
				sysAdminDep.GetId():                               {},
				depWithAllResourceLimitsRequestsSpecified.GetId(): {},
				depWithEnforcementBypassAnnotation.GetId():        {},
				hostMountDep.GetId():                              {},
			},
			sampleViolationForMatched: "Image has not been scanned",
		},
		{
			policyName:                "Required Label: Email",
			skip:                      true,
			shouldNotMatch:            map[string]struct{}{fixtureDep.GetId(): {}},
			sampleViolationForMatched: "Required label not found (key = 'email', value = '[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+')",
		},
		{
			policyName:                "Required Annotation: Email",
			skip:                      true,
			shouldNotMatch:            map[string]struct{}{depWithGoodEmailAnnotation.GetId(): {}},
			sampleViolationForMatched: "Required annotation not found (key = 'email', value = '[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+')",
		},
		{
			policyName:                "Required Label: Owner",
			skip:                      true,
			shouldNotMatch:            map[string]struct{}{fixtureDep.GetId(): {}},
			sampleViolationForMatched: "Required label not found (key = 'owner', value = '.+')",
		},
		{
			policyName:                "Required Annotation: Owner",
			skip:                      true,
			shouldNotMatch:            map[string]struct{}{depWithOwnerAnnotation.GetId(): {}},
			sampleViolationForMatched: "Required annotation not found (key = 'owner', value = '.+')",
		},
		{
			policyName: "CAP_SYS_ADMIN capability added",
			expectedViolations: map[string][]*storage.Alert_Violation{
				sysAdminDep.GetId(): {
					{
						Message: "CAP_SYS_ADMIN was in the ADD CAPABILITIES list",
					},
				},
			},
		},
		{
			policyName: "Shellshock: Multiple CVEs",
			expectedViolations: map[string][]*storage.Alert_Violation{
				shellshockDep.GetId(): {
					{
						Message: "CVE CVE-2014-6271 matched regex 'CVE-2014-(6271|6277|6278|7169|7186|7187)'",
					},
				},
				fixtureDep.GetId(): {
					{
						Message: "CVE CVE-2014-6271 matched regex 'CVE-2014-(6271|6277|6278|7169|7186|7187)'",
					},
				},
			},
		},
		{
			policyName: "Apache Struts: CVE-2017-5638",
			expectedViolations: map[string][]*storage.Alert_Violation{
				strutsDep.GetId(): {
					{
						Message: "CVE CVE-2017-5638 matched regex 'CVE-2017-5638'",
					},
				},
			},
		},
		{
			policyName: "Heartbleed: CVE-2014-0160",
			expectedViolations: map[string][]*storage.Alert_Violation{
				heartbleedDep.GetId(): {
					{
						Message: "CVE CVE-2014-0160 matched regex 'CVE-2014-0160'",
					},
				},
			},
		},
		{
			policyName: "No resource requests or limits specified",
			skip:       true,
			expectedViolations: map[string][]*storage.Alert_Violation{
				fixtureDep.GetId(): {
					{Message: "The CPU resource limit of 0 is equal to the threshold of 0.00"},
					{Message: "The memory resource limit of 0 is equal to the threshold of 0.00"},
					{Message: "The memory resource request of 0 is equal to the threshold of 0.00"},
				},
			},
		},
		{
			policyName: "Environment Variable Contains Secret",
			skip:       true,
			expectedViolations: map[string][]*storage.Alert_Violation{
				secretEnvDep.GetId(): {
					{
						Message: "Container Environment (key='THIS_IS_SECRET_VAR', value='stealthmode') matched environment policy (key = '.*SECRET.*|.*PASSWORD.*', value from = 'RAW')",
					},
				},
			},
		},
		{
			policyName: "Secret Mounted as Environment Variable",
			skip:       true,
			expectedViolations: map[string][]*storage.Alert_Violation{
				secretKeyRefDep.GetId(): {
					{
						Message: "Container Environment (key='THIS_IS_SECRET_VAR', value='') matched environment policy (value from = 'SECRET_KEY')",
					},
				},
			},
		},
		{
			policyName: "Fixable CVSS >= 6 and Privileged",
			expectedViolations: map[string][]*storage.Alert_Violation{
				heartbleedDep.GetId(): {
					{
						Message: "Found a CVSS score of 6 (greater than or equal to 6.0) (cve: CVE-2014-0160) that is fixable",
					},
					{
						Message: "Privileged container found",
					},
				},
			},
		},
		{
			policyName: "Fixable CVSS >= 7",
			expectedViolations: map[string][]*storage.Alert_Violation{
				strutsDep.GetId(): {
					{
						Message: "Found a CVSS score of 8 (greater than or equal to 7.0) (cve: CVE-2017-5638) that is fixable",
					},
				},
			},
		},
		{
			policyName: "ADD Command used instead of COPY",
			expectedViolations: map[string][]*storage.Alert_Violation{
				addDockerFileDep.GetId(): {
					{
						Message: "Dockerfile Line 'ADD deploy.sh' matches the rule ADD .*",
					},
				},
				fixtureDep.GetId(): {
					{
						Message: "Dockerfile Line 'ADD FILE:blah' matches the rule ADD .*",
					},
					{
						Message: "Dockerfile Line 'ADD file:4eedf861fb567fffb2694b65ebdd58d5e371a2c28c3863f363f333cb34e5eb7b in /' matches the rule ADD .*",
					},
				},
			},
		},
		{
			policyName: "nmap Execution",
			expectedProcessViolations: map[string][]*storage.ProcessIndicator{
				fixtureDep.GetId():  {nmapIndicatorfixtureDep1, nmapIndicatorfixtureDep2},
				nginx110Dep.GetId(): {nmapIndicatorNginx110Dep},
			},
		},
		{
			policyName: "Process Targeting Cluster Kubelet Endpoint",
			expectedProcessViolations: map[string][]*storage.ProcessIndicator{
				containerPort22Dep.GetId(): {kubeletIndicator, kubeletIndicator2},
			},
		},
		{
			policyName: "Ubuntu Package Manager Execution",
			expectedProcessViolations: map[string][]*storage.ProcessIndicator{
				fixtureDep.GetId():  {fixtureDepAptIndicator},
				sysAdminDep.GetId(): {sysAdminDepAptIndicator},
			},
		},
		{
			policyName: "Process with UID 0",
			expectedProcessViolations: map[string][]*storage.ProcessIndicator{
				fixtureDep.GetId(): {fixtureDepJavaIndicator},
			},
		},
		{
			policyName: "Shell Spawned by Java Application",
			expectedProcessViolations: map[string][]*storage.ProcessIndicator{
				fixtureDep.GetId(): {fixtureDepJavaIndicator},
			},
		},
		{
			policyName: "Emergency Deployment Annotation",
			skip:       true,
			expectedViolations: map[string][]*storage.Alert_Violation{
				depWithEnforcementBypassAnnotation.GetId(): {
					{Message: "Disallowed annotation found (key = 'admission.stackrox.io/break-glass')"},
				},
			},
		},
		{
			policyName: "Mounting Sensitive Host Directories",
			expectedViolations: map[string][]*storage.Alert_Violation{
				hostMountDep.GetId(): {
					{Message: "Volume source '/etc/passwd' matched (/etc/.*|/sys/.*|/dev/.*|/proc/.*|/var/.*)"},
					{Message: "Volume source '/var/lib/kubelet' matched (/etc/.*|/sys/.*|/dev/.*|/proc/.*|/var/.*)"},
				},
				dockerSockDep.GetId(): {
					{Message: "Volume source '/var/run/docker.sock' matched (/etc/.*|/sys/.*|/dev/.*|/proc/.*|/var/.*)"},
				},
			},
		},
	}

	for _, c := range deploymentTestCases {
		p := suite.MustGetPolicy(c.policyName)
		// Skip unsupported tests, but ensure we don't continue to do this by the time we enable
		// the flag.
		if c.skip && !features.BooleanPolicyLogic.Enabled() {
			continue
		}
		suite.T().Run(fmt.Sprintf("%s (on deployments)", c.policyName), func(t *testing.T) {
			if len(c.shouldNotMatch) == 0 {
				assert.True(t, (c.expectedViolations != nil) != (c.expectedProcessViolations != nil), "Every test case must "+
					"contain exactly one of expectedViolations and expectedProcessViolations")
			} else {
				assert.Nil(t, c.expectedViolations, "Cannot specify shouldNotMatch AND expectedViolations")
				assert.Nil(t, c.expectedProcessViolations, "Cannot specify shouldNotMatch AND expectedProcessViolations")
			}

			convertedP, err := CloneAndEnsureConverted(p)
			require.NoError(t, err)
			m, err := BuildDeploymentMatcher(convertedP)
			require.NoError(t, err)

			if c.expectedProcessViolations != nil {
				for deploymentID, processes := range c.expectedProcessViolations {
					expectedProcesses := set.NewStringSet(sliceutils.Map(processes, func(p *storage.ProcessIndicator) string {
						return p.GetId()
					}).([]string)...)
					deployment := suite.deployments[deploymentID]

					for _, process := range suite.deploymentsToIndicators[deploymentID] {
						match, err := m.MatchDeployment(context.Background(), deployment, suite.getImagesForDeployment(deployment), process)
						require.NoError(t, err)
						if expectedProcesses.Contains(process.GetId()) {
							assert.NotNil(t, match.ProcessViolation, "process %+v should match", process)
						} else {
							assert.Nil(t, match.ProcessViolation, "process %+v should not match", process)
						}
					}
				}
				return
			}

			actualViolations := make(map[string][]*storage.Alert_Violation)
			for id, deployment := range suite.deployments {
				violationsForDep, err := m.MatchDeployment(context.Background(), deployment, suite.getImagesForDeployment(deployment), nil)
				require.NoError(t, err)
				assert.Nil(t, violationsForDep.ProcessViolation)
				if alertViolations := violationsForDep.AlertViolations; len(alertViolations) > 0 {
					actualViolations[id] = alertViolations
				}
			}
			if len(c.shouldNotMatch) > 0 {
				for shouldNotMatchID := range c.shouldNotMatch {
					assert.NotContains(t, actualViolations, shouldNotMatchID)
				}
				for id := range suite.deployments {
					if _, shouldNotMatch := c.shouldNotMatch[id]; !shouldNotMatch {
						assert.Contains(t, actualViolations, id)
					}
				}
				return
			}
			for id := range suite.deployments {
				_, expected := c.expectedViolations[id]
				if expected {
					assert.Contains(t, actualViolations, id)
				} else {
					assert.NotContains(t, actualViolations, id)
				}
			}

		})
	}

	imageTestCases := []testCase{
		{
			policyName: "Latest tag",
			expectedViolations: map[string][]*storage.Alert_Violation{
				fixtureDep.GetContainers()[1].GetImage().GetId(): {
					{Message: "Image tag 'latest' matched latest"},
				},
			},
		},
		{
			policyName: "DockerHub NGINX 1.10",
			expectedViolations: map[string][]*storage.Alert_Violation{
				fixtureDep.GetContainers()[0].GetImage().GetId(): {
					{
						Message: "Image tag '1.10' matched 1.10",
					},
					{
						Message: "Image registry 'docker.io' matched docker.io",
					},
					{
						Message: "Image remote 'library/nginx' matched nginx",
					},
				},
				suite.imageIDFromDep(nginx110Dep): {
					{
						Message: "Image tag '1.10' matched 1.10",
					},
					{
						Message: "Image registry 'docker.io' matched docker.io",
					},
					{
						Message: "Image remote 'library/nginx' matched nginx",
					},
				},
			},
		},
		{
			policyName: "Alpine Linux Package Manager (apk) in Image",
			expectedViolations: map[string][]*storage.Alert_Violation{
				suite.imageIDFromDep(apkDep): {
					{
						Message: "Component name 'apk' matched apk",
					},
				},
			},
		},
		{
			policyName: "Ubuntu Package Manager in Image",
			expectedViolations: map[string][]*storage.Alert_Violation{
				suite.imageIDFromDep(componentDeps["apt"]): {
					{
						Message: "Component name 'apt' matched apt|dpkg",
					},
				},
			},
		},
		{
			policyName: "Curl in Image",
			expectedViolations: map[string][]*storage.Alert_Violation{
				suite.imageIDFromDep(curlDep): {
					{
						Message: "Component name 'curl' matched curl",
					},
				},
			},
		},
		{
			policyName: "Red Hat Package Manager in Image",
			expectedViolations: map[string][]*storage.Alert_Violation{
				suite.imageIDFromDep(componentDeps["dnf"]): {
					{
						Message: "Component name 'dnf' matched rpm|dnf|yum",
					},
				},
			},
		},
		{
			policyName: "Wget in Image",
			expectedViolations: map[string][]*storage.Alert_Violation{
				suite.imageIDFromDep(componentDeps["wget"]): {
					{
						Message: "Component name 'wget' matched wget",
					},
				},
			},
		},
		{
			policyName: "90-Day Image Age",
			skip:       true,
			expectedViolations: map[string][]*storage.Alert_Violation{
				suite.imageIDFromDep(oldImageDep): {
					{
						Message: fmt.Sprintf("Time of image creation '%s' was more than 90 days ago", readable.Time(oldImageCreationTime)),
					},
				},
			},
		},
		{
			policyName: "30-Day Scan Age",
			skip:       true,
			expectedViolations: map[string][]*storage.Alert_Violation{
				suite.imageIDFromDep(oldScannedDep): {
					{
						Message: fmt.Sprintf("Time of last scan '%s' was more than 30 days ago", readable.Time(oldScannedTime)),
					},
				},
			},
		},
		{
			policyName: "Secure Shell (ssh) Port Exposed in Image",
			expectedViolations: map[string][]*storage.Alert_Violation{
				suite.imageIDFromDep(imagePort22Dep): {
					{
						Message: "Dockerfile Line 'EXPOSE 22/tcp' matches the rule EXPOSE (22/tcp|\\s+22/tcp)",
					},
				},
			},
		},
		{
			policyName: "Insecure specified in CMD",
			expectedViolations: map[string][]*storage.Alert_Violation{
				suite.imageIDFromDep(insecureCMDDep): {
					{
						Message: "Dockerfile Line 'CMD do an insecure thing' matches the rule CMD .*insecure.*",
					},
				},
			},
		},
		{
			policyName: "Improper Usage of Orchestrator Secrets Volume",
			expectedViolations: map[string][]*storage.Alert_Violation{
				suite.imageIDFromDep(runSecretsDep): {
					{
						Message: "Dockerfile Line 'VOLUME /run/secrets' matches the rule VOLUME /run/secrets",
					},
				},
			},
		},
		{
			policyName: "Images with no scans",
			skip:       true,
			shouldNotMatch: map[string]struct{}{
				fixtureDep.GetContainers()[0].GetImage().GetId(): {},
				fixtureDep.GetContainers()[1].GetImage().GetId(): {},
				suite.imageIDFromDep(oldScannedDep):              {},
			},
			sampleViolationForMatched: "Image has not been scanned",
		},
		{
			policyName: "Shellshock: Multiple CVEs",
			expectedViolations: map[string][]*storage.Alert_Violation{
				suite.imageIDFromDep(shellshockDep): {
					{
						Message: "CVE CVE-2014-6271 matched regex 'CVE-2014-(6271|6277|6278|7169|7186|7187)'",
					},
				},
				fixtureDep.GetContainers()[1].GetImage().GetId(): {
					{
						Message: "CVE CVE-2014-6271 matched regex 'CVE-2014-(6271|6277|6278|7169|7186|7187)'",
					},
				},
			},
		},
		{
			policyName: "Apache Struts: CVE-2017-5638",
			expectedViolations: map[string][]*storage.Alert_Violation{
				suite.imageIDFromDep(strutsDep): {
					{
						Message: "CVE CVE-2017-5638 matched regex 'CVE-2017-5638'",
					},
				},
			},
		},
		{
			policyName: "Heartbleed: CVE-2014-0160",
			expectedViolations: map[string][]*storage.Alert_Violation{
				suite.imageIDFromDep(heartbleedDep): {
					{
						Message: "CVE CVE-2014-0160 matched regex 'CVE-2014-0160'",
					},
				},
			},
		},
		{
			policyName: "Fixable CVSS >= 7",
			expectedViolations: map[string][]*storage.Alert_Violation{
				suite.imageIDFromDep(strutsDep): {
					{
						Message: "Found a CVSS score of 8 (greater than or equal to 7.0) (cve: CVE-2017-5638) that is fixable",
					},
				},
			},
		},
		{
			policyName: "ADD Command used instead of COPY",
			expectedViolations: map[string][]*storage.Alert_Violation{
				suite.imageIDFromDep(addDockerFileDep): {
					{
						Message: "Dockerfile Line 'ADD deploy.sh' matches the rule ADD .*",
					},
				},
				fixtureDep.GetContainers()[0].GetImage().GetId(): {
					{
						Message: "Dockerfile Line 'ADD FILE:blah' matches the rule ADD .*",
					},
				},
				fixtureDep.GetContainers()[1].GetImage().GetId(): {
					{
						Message: "Dockerfile Line 'ADD file:4eedf861fb567fffb2694b65ebdd58d5e371a2c28c3863f363f333cb34e5eb7b in /' matches the rule ADD .*",
					},
				},
			},
		},
		{
			policyName: "Required Image Label",
			skip:       true,
			shouldNotMatch: map[string]struct{}{
				"requiredImageLabelImage": {},
			},
		},
	}

	for _, c := range imageTestCases {
		// Temporarily skip unsupported tests.
		if !features.BooleanPolicyLogic.Enabled() && c.skip {
			continue
		}
		p := suite.MustGetPolicy(c.policyName)
		suite.T().Run(fmt.Sprintf("%s (on images)", c.policyName), func(t *testing.T) {
			assert.Nil(t, c.expectedProcessViolations)

			convertedP, err := CloneAndEnsureConverted(p)
			require.NoError(t, err)
			m, err := BuildImageMatcher(convertedP)
			require.NoError(t, err)

			actualViolations := make(map[string][]*storage.Alert_Violation)
			for id, image := range suite.images {
				violationsForImg, err := m.MatchImage(context.Background(), image)
				suite.Require().NoError(err)
				suite.Nil(violationsForImg.ProcessViolation)
				if alertViolations := violationsForImg.AlertViolations; len(alertViolations) > 0 {
					actualViolations[id] = alertViolations
				}
			}

			for id := range c.expectedViolations {
				assert.Contains(t, actualViolations, id)
			}
			if len(c.shouldNotMatch) > 0 {
				for shouldNotMatchID := range c.shouldNotMatch {
					assert.NotContains(t, actualViolations, shouldNotMatchID)
				}
				for id := range suite.images {
					if _, shouldNotMatch := c.shouldNotMatch[id]; !shouldNotMatch {
						assert.Contains(t, actualViolations, id)
					}
				}
			}
		})
	}
}

func (suite *DefaultPoliciesTestSuite) TestMapPolicyMatchOne() {
	if !features.BooleanPolicyLogic.Enabled() {
		return
	}
	noAnnotation := &storage.Deployment{
		Id: "noAnnotation",
	}
	suite.addDepAndImages(noAnnotation)

	validAnnotation := &storage.Deployment{
		Id: "validAnnotation",
		Annotations: map[string]string{
			"email": "joseph@rules.gov",
		},
	}
	suite.addDepAndImages(validAnnotation)

	policy := suite.defaultPolicies["Required Annotation: Email"]
	m, err := BuildDeploymentMatcher(policy)
	suite.NoError(err)

	matched, err := m.MatchDeployment(context.Background(), noAnnotation, nil, nil)
	suite.NoError(err)
	suite.Len(matched.AlertViolations, 1)

	matched, err = m.MatchDeployment(context.Background(), validAnnotation, nil, nil)
	suite.NoError(err)
	suite.Empty(matched.AlertViolations)
}

func (suite *DefaultPoliciesTestSuite) TestRuntimePolicyFieldsCompile() {
	for _, p := range suite.defaultPolicies {
		if policyUtils.AppliesAtRunTime(p) && p.GetFields().GetProcessPolicy() != nil {
			processPolicy := p.GetFields().GetProcessPolicy()
			if processPolicy.GetName() != "" {
				regexp.MustCompile(processPolicy.GetName())
			}
			if processPolicy.GetArgs() != "" {
				regexp.MustCompile(processPolicy.GetArgs())
			}
			if processPolicy.GetAncestor() != "" {
				regexp.MustCompile(processPolicy.GetAncestor())
			}
		}
	}
}

func policyWithSingleKeyValue(fieldName, value string, negate bool) *storage.Policy {
	group := &storage.PolicyGroup{FieldName: fieldName, Values: []*storage.PolicyValue{{Value: value}}, Negate: negate}
	return &storage.Policy{
		PolicyVersion:  Version,
		Name:           uuid.NewV4().String(),
		PolicySections: []*storage.PolicySection{{PolicyGroups: []*storage.PolicyGroup{group}}},
	}
}

func (suite *DefaultPoliciesTestSuite) TestK8sRBAC() {
	deployments := make(map[string]*storage.Deployment)
	for permissionLevelStr, permissionLevel := range storage.PermissionLevel_value {
		dep := fixtures.GetDeployment().Clone()
		dep.ServiceAccountPermissionLevel = storage.PermissionLevel(permissionLevel)
		deployments[permissionLevelStr] = dep
	}

	for _, testCase := range []struct {
		value           string
		negate          bool
		expectedMatches []string
	}{
		{
			"DEFAULT",
			false,
			[]string{"DEFAULT", "ELEVATED_IN_NAMESPACE", "ELEVATED_CLUSTER_WIDE", "CLUSTER_ADMIN"},
		},
		{
			"ELEVATED_CLUSTER_WIDE",
			false,
			[]string{"ELEVATED_CLUSTER_WIDE", "CLUSTER_ADMIN"},
		},
		{
			"CLUSTER_ADMIN",
			false,
			[]string{"CLUSTER_ADMIN"},
		},
		{
			"ELEVATED_CLUSTER_WIDE",
			true,
			[]string{"NONE", "DEFAULT", "ELEVATED_IN_NAMESPACE"},
		},
	} {
		c := testCase
		suite.T().Run(fmt.Sprintf("%+v", c), func(t *testing.T) {
			matcher, err := BuildDeploymentMatcher(policyWithSingleKeyValue(MinimumRBACPermissions, c.value, c.negate))
			require.NoError(t, err)
			matched := set.NewStringSet()
			for depRef, dep := range deployments {
				violations, err := matcher.MatchDeployment(context.Background(), dep, suite.getImagesForDeployment(dep), nil)
				require.NoError(t, err)
				if len(violations.AlertViolations) > 0 {
					matched.Add(depRef)
				}
			}
			assert.ElementsMatch(t, matched.AsSlice(), c.expectedMatches, "Got %v, expected: %v", matched.AsSlice(), c.expectedMatches)
		})
	}
}
