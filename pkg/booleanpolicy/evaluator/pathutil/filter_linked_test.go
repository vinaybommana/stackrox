package pathutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type simplePathHolder Path

func (p *simplePathHolder) GetPath() *Path {
	return (*Path)(p)
}

func pathFromSteps(t *testing.T, steps ...interface{}) PathHolder {
	return (*simplePathHolder)(PathFromSteps(t, steps...))
}

func TestFilterLinked(t *testing.T) {
	for _, testCase := range []struct {
		desc                  string
		input, expectedOutput map[string][]PathHolder
	}{
		{
			desc: "Simple case, no arrays at all",
			input: map[string][]PathHolder{
				"Namespace": {pathFromSteps(t, "Namespace")},
			},
			expectedOutput: map[string][]PathHolder{
				"Namespace": {pathFromSteps(t, "Namespace")},
			},
		},
		{
			desc: "One array, same level, no match",
			input: map[string][]PathHolder{
				"VolumeName":   {pathFromSteps(t, "Volumes", 2, "Name")},
				"VolumeSource": {pathFromSteps(t, "Volumes", 1, "Source")},
			},
			expectedOutput: nil,
		},
		{
			desc: "One array, same level",
			input: map[string][]PathHolder{
				"VolumeName":   {pathFromSteps(t, "Volumes", 0, "Name"), pathFromSteps(t, "Volumes", 2, "Name"), pathFromSteps(t, "Volumes", 3, "Name")},
				"VolumeSource": {pathFromSteps(t, "Volumes", 0, "Source"), pathFromSteps(t, "Volumes", 1, "Source"), pathFromSteps(t, "Volumes", 3, "Source")},
			},
			expectedOutput: map[string][]PathHolder{
				"VolumeName":   {pathFromSteps(t, "Volumes", 0, "Name"), pathFromSteps(t, "Volumes", 3, "Name")},
				"VolumeSource": {pathFromSteps(t, "Volumes", 0, "Source"), pathFromSteps(t, "Volumes", 3, "Source")},
			},
		},
		{
			desc: "Complex case, multi-level linking but no branching, no match",
			input: map[string][]PathHolder{
				"Namespace":     {pathFromSteps(t, "Namespace")},
				"ContainerName": {pathFromSteps(t, "Containers", 0, "Name")},
				"VolumeName":    {pathFromSteps(t, "Containers", 1, "Volumes", 0, "Name")},
				"VolumeSource":  {pathFromSteps(t, "Containers", 1, "Volumes", 0, "Source"), pathFromSteps(t, "Containers", 1, "Volumes", 1, "Source")},
			},
			expectedOutput: nil,
		},
		{
			desc: "Complex case, multi-level linking but no branching",
			input: map[string][]PathHolder{
				"Namespace":     {pathFromSteps(t, "Namespace")},
				"ContainerName": {pathFromSteps(t, "Containers", 0, "Name"), pathFromSteps(t, "Containers", 1, "Name")},
				"VolumeName":    {pathFromSteps(t, "Containers", 1, "Volumes", 0, "Name")},
				"VolumeSource":  {pathFromSteps(t, "Containers", 1, "Volumes", 0, "Source"), pathFromSteps(t, "Containers", 1, "Volumes", 1, "Source")},
			},
			expectedOutput: map[string][]PathHolder{
				"Namespace":     {pathFromSteps(t, "Namespace")},
				"ContainerName": {pathFromSteps(t, "Containers", 1, "Name")},
				"VolumeName":    {pathFromSteps(t, "Containers", 1, "Volumes", 0, "Name")},
				"VolumeSource":  {pathFromSteps(t, "Containers", 1, "Volumes", 0, "Source")},
			},
		},
		{
			desc: "Complex case, multi-level linking plus branching no match",
			input: map[string][]PathHolder{
				"Namespace":     {pathFromSteps(t, "Namespace")},
				"ContainerName": {pathFromSteps(t, "Containers", 0, "Name"), pathFromSteps(t, "Containers", 1, "Name")},
				"VolumeName":    {pathFromSteps(t, "Containers", 1, "Volumes", 0, "Name")},
				"VolumeSource":  {pathFromSteps(t, "Containers", 1, "Volumes", 0, "Source"), pathFromSteps(t, "Containers", 1, "Volumes", 1, "Source")},
				"PortName":      {pathFromSteps(t, "Containers", 0, "Ports", 0, "Name"), pathFromSteps(t, "Containers", 1, "Ports", 0, "Name")},
				"PortProtocol":  {pathFromSteps(t, "Containers", 0, "Ports", 0, "Protocol"), pathFromSteps(t, "Containers", 1, "Ports", 1, "Protocol")},
			},
			expectedOutput: nil,
		},
		{
			desc: "Complex case, multi-level linking plus branching",
			input: map[string][]PathHolder{
				"Namespace":     {pathFromSteps(t, "Namespace")},
				"ContainerName": {pathFromSteps(t, "Containers", 0, "Name"), pathFromSteps(t, "Containers", 1, "Name")},
				"VolumeName":    {pathFromSteps(t, "Containers", 1, "Volumes", 0, "Name")},
				"VolumeSource":  {pathFromSteps(t, "Containers", 1, "Volumes", 0, "Source"), pathFromSteps(t, "Containers", 1, "Volumes", 1, "Source")},
				"PortName":      {pathFromSteps(t, "Containers", 0, "Ports", 0, "Name"), pathFromSteps(t, "Containers", 1, "Ports", 0, "Name")},
				"PortProtocol":  {pathFromSteps(t, "Containers", 0, "Ports", 0, "Protocol"), pathFromSteps(t, "Containers", 1, "Ports", 0, "Protocol")},
			},
			expectedOutput: map[string][]PathHolder{
				"Namespace":     {pathFromSteps(t, "Namespace")},
				"ContainerName": {pathFromSteps(t, "Containers", 1, "Name")},
				"VolumeName":    {pathFromSteps(t, "Containers", 1, "Volumes", 0, "Name")},
				"VolumeSource":  {pathFromSteps(t, "Containers", 1, "Volumes", 0, "Source")},
				"PortName":      {pathFromSteps(t, "Containers", 1, "Ports", 0, "Name")},
				"PortProtocol":  {pathFromSteps(t, "Containers", 1, "Ports", 0, "Protocol")},
			},
		},
	} {
		c := testCase
		t.Run(c.desc, func(t *testing.T) {
			t.Parallel()
			output, matched, err := FilterPathsToLinkedMatches(c.input)
			require.NoError(t, err)
			if c.expectedOutput == nil {
				assert.Empty(t, output)
				assert.False(t, matched)
				return
			}
			assert.True(t, matched)
			assert.Equal(t, c.expectedOutput, output)
		})
	}
}
