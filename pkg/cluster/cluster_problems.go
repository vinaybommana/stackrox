package cluster

import (
	"fmt"

	"github.com/stackrox/rox/generated/storage"
)

type Problem interface {
	ToProto() *storage.ClusterProblem
}

type SensorPolicyVersionMismatch struct {
	CentralVersion       string
	CentralPolicyVersion string

	SensorVersion       string
	SensorPolicyVersion string
}

func (p SensorPolicyVersionMismatch) ToProto() *storage.ClusterProblem {
	descr := fmt.Sprintf("Policy version used by Sensor (%v) does not match the version used by Central (%v)." +
		" This is not necessarily critical because some versions can be downgraded.",
		p.SensorPolicyVersion, p.CentralPolicyVersion)

	remedy := "Running the same version of Central and Sensor guarantees there is no policy version mismatch."
	if p.SensorVersion < p.CentralVersion {
		remedy += fmt.Sprintf(" Upgrade Sensor to version %v.", p.CentralVersion)
	}
	if p.CentralVersion < p.SensorVersion {
		remedy += fmt.Sprintf(" Upgrade Central to version %v.", p.SensorVersion)
	}

	return &storage.ClusterProblem{
		ShortName:   "Sensor Policy Version Mismatch",
		Description: descr,
		Remedy:      remedy,
	}
}

type IncompatibleSensorPolicyVersion struct {
	CentralVersion       string
	CentralPolicyVersion string

	SensorPolicyVersion string
}

func (p IncompatibleSensorPolicyVersion) ToProto() *storage.ClusterProblem {
	descr := fmt.Sprintf("Policy version used by Sensor (%v) is incompatible with the version used by Central (%v)." +
		" This means that Sensor does not understand policies advertised by Central and hence cannot enforce them." +
		" The functionality of your StackRox installation is significantly reduced.",
		p.SensorPolicyVersion, p.CentralPolicyVersion)

	remedy := fmt.Sprintf("Upgrade Sensor to version %v.", p.CentralVersion)

	return &storage.ClusterProblem{
		ShortName:   "Incompatible Sensor Policy Version",
		Description: descr,
		Remedy:      remedy,
	}
}
