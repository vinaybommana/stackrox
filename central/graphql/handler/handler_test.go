package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	clusterMocks "github.com/stackrox/rox/central/cluster/datastore/mocks"
	deploymentMocks "github.com/stackrox/rox/central/deployment/datastore/mocks"
	"github.com/stackrox/rox/central/graphql/resolvers"
	"github.com/stackrox/rox/central/graphql/schema"
	processMocks "github.com/stackrox/rox/central/processindicator/datastore/mocks"
	"github.com/stackrox/rox/pkg/grpc/authz/allow"
)

const (
	fakeClusterID = "fakeClusterId"
)

var (
	jsonPathPattern = regexp.MustCompile(`^(\[\d+\]|\.[^[.]+)`)
)

type mocks struct {
	cluster    *clusterMocks.MockDataStore
	deployment *deploymentMocks.MockDataStore
	process    *processMocks.MockDataStore
	resolver   *resolvers.Resolver
}

func mockResolver(t *testing.T) mocks {
	ctrl := gomock.NewController(t)
	cluster := clusterMocks.NewMockDataStore(ctrl)
	deployment := deploymentMocks.NewMockDataStore(ctrl)
	process := processMocks.NewMockDataStore(ctrl)

	resolver := &resolvers.Resolver{
		ClusterDataStore:      cluster,
		DeploymentDataStore:   deployment,
		ProcessIndicatorStore: process,
	}

	return mocks{
		cluster:    cluster,
		deployment: deployment,
		resolver:   resolver,
		process:    process,
	}
}

func assertJSONMatches(t *testing.T, buffer *bytes.Buffer, path string, expected string) (ok bool) {
	ok = false
	msg := &json.RawMessage{}
	json.Unmarshal(buffer.Bytes(), msg)
	for path != "" {
		if msg == nil {
			t.Errorf("No message found (remaining path %q)", path)
			return
		}
		indices := jsonPathPattern.FindStringIndex(path)
		if indices == nil {
			t.Errorf("Invalid path segment: %q", path)
			return
		}
		segment := path[indices[0]:indices[1]]
		path = path[indices[1]:]
		if segment[0] == '[' {
			index, err := strconv.ParseInt(segment[1:len(segment)-1], 10, 32)
			if err != nil {
				t.Errorf("Invalid array index: %q", segment)
				return
			}
			array := make([]*json.RawMessage, 0)
			err = json.Unmarshal([]byte(*msg), &array)
			if err != nil {
				t.Error(err)
				return
			}
			msg = array[index]
		} else {
			m := make(map[string]*json.RawMessage)
			err := json.Unmarshal([]byte(*msg), &m)
			if err != nil {
				t.Error(err)
				return
			}
			ok = true
			msg, ok = m[segment[1:]]
			if !ok {
				t.Errorf("Key not found: %q", segment)
				for k := range m {
					t.Errorf(" -- key %q", k)
				}
				ok = false
				return
			}
		}
	}
	actual := ""
	err := json.Unmarshal([]byte(*msg), &actual)
	ok = err == nil && expected == actual
	if err != nil {
		t.Error(err)
		return
	}
	if expected != actual {
		t.Errorf("Expected %q, actual %q", expected, actual)
	}
	return
}

type graphqlError struct {
	Errors []struct {
		Message   string `json:"message"`
		Locations []struct {
		} `json:"locations"`
	} `json:"errors"`
}

func assertNoErrors(t *testing.T, msg *bytes.Buffer) bool {
	errors := graphqlError{}
	err := json.Unmarshal(msg.Bytes(), &errors)
	if err != nil {
		t.Logf("Input is ok? %q", err)
		return true
	}
	for _, e := range errors.Errors {
		t.Fatalf("%s", e.Message)
	}
	return false
}

func executeTestQuery(t *testing.T, mocks mocks, query string) *httptest.ResponseRecorder {
	return executeTestQueryWithVariables(t, mocks, query, nil)
}

func executeTestQueryWithVariables(t *testing.T, mocks mocks, query string, variables map[string]string) *httptest.ResponseRecorder {
	ourSchema, err := graphql.ParseSchema(schema.Schema(), mocks.resolver)
	if err != nil {
		t.Fatal(err)
	}
	h := &relay.Handler{Schema: ourSchema}
	vals := map[string]interface{}{"query": query}
	if variables != nil {
		vals["variables"] = variables
	}
	b, err := json.Marshal(vals)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("POST", "/api/graphql", bytes.NewReader(b)).WithContext(
		resolvers.SetAuthorizerOverride(context.Background(), allow.Anonymous()))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	assertNoErrors(t, rec.Body)
	return rec
}

func TestSchemaValidates(t *testing.T) {
	s := schema.Schema()
	_, err := graphql.ParseSchema(s, mockResolver(t).resolver)
	if err != nil {
		t.Log(s)
		t.Error(err)
		t.Error("You might need to run `go generate .` in central/graphql/resolvers")
	}
}
