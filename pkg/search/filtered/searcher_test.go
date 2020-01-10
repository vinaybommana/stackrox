package filtered

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/badgerhelper"
	"github.com/stackrox/rox/pkg/dackbox/graph"
	"github.com/stackrox/rox/pkg/dackbox/graph/mocks"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/search"
	searchMocks "github.com/stackrox/rox/pkg/search/blevesearch/mocks"
	"github.com/stretchr/testify/suite"
)

var (
	prefix1 = []byte("pre1")
	prefix2 = []byte("pre2")
	prefix3 = []byte("namespace")
	prefix4 = []byte("cluster")

	id1 = []byte("id1")
	id2 = []byte("id2")
	id3 = []byte("id3")
	id4 = []byte("id4")
	id5 = []byte("id5")
	id6 = []byte("id6")
	id7 = []byte("id7")
	id8 = []byte("id8")

	prefixedID1 = badgerhelper.GetBucketKey(prefix1, id1)
	prefixedID2 = badgerhelper.GetBucketKey(prefix1, id2)
	prefixedID3 = badgerhelper.GetBucketKey(prefix2, id3)
	prefixedID4 = badgerhelper.GetBucketKey(prefix2, id4)
	prefixedID5 = badgerhelper.GetBucketKey(prefix3, id5)
	prefixedID6 = badgerhelper.GetBucketKey(prefix3, id6)
	prefixedID7 = badgerhelper.GetBucketKey(prefix4, id7)
	prefixedID8 = badgerhelper.GetBucketKey(prefix4, id8)

	// id1 -> id3 -> id5 (namespace) -> id7 (cluster)
	// id2 -> id4 -> id6 (namespace) -> id7 (cluster)
	toID1 = [][]byte{prefixedID3}
	toID2 = [][]byte{prefixedID4}
	toID3 = [][]byte{prefixedID5}
	toID4 = [][]byte{prefixedID6}
	toID5 = [][]byte{prefixedID7, prefixedID8}
	toID6 = [][]byte{prefixedID7}
)

func TestIDTransformation(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(idTransformationTestSuite))
}

type idTransformationTestSuite struct {
	suite.Suite

	mockRGraph         *mocks.MockRGraph
	mockUnsafeSearcher *searchMocks.MockUnsafeSearcher

	mockCtrl *gomock.Controller
}

func (s *idTransformationTestSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockRGraph = mocks.NewMockRGraph(s.mockCtrl)
	s.mockUnsafeSearcher = searchMocks.NewMockUnsafeSearcher(s.mockCtrl)
}

func (s *idTransformationTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func (s *idTransformationTestSuite) TestGlobalAllowed() {
	s.mockUnsafeSearcher.EXPECT().Search(gomock.Any()).Return([]search.Result{
		{
			ID: string(id1),
		},
		{
			ID: string(id2),
		},
	}, nil)

	checker := sac.NewScopeChecker(sac.AllowAllAccessScopeChecker())

	filter, err := NewSACFilter(
		WithScopeChecker(checker),
	)
	s.NoError(err, "filter creation should have succeeded")

	searcher := Searcher(s.mockUnsafeSearcher, filter)
	results, err := searcher.Search(context.Background(), &v1.Query{})
	s.NoError(err, "search should have succeeded")
	s.Equal([]search.Result{
		{
			ID: string(id1),
		},
		{
			ID: string(id2),
		},
	}, results)
}

func (s *idTransformationTestSuite) TestGlobalDenied() {
	s.mockUnsafeSearcher.EXPECT().Search(gomock.Any()).Return([]search.Result{
		{
			ID: string(id1),
		},
		{
			ID: string(id2),
		},
	}, nil)

	checker := sac.NewScopeChecker(sac.DenyAllAccessScopeChecker())

	filter, err := NewSACFilter(
		WithScopeChecker(checker),
	)
	s.NoError(err, "filter creation should have succeeded")

	searcher := Searcher(s.mockUnsafeSearcher, filter)
	results, err := searcher.Search(context.Background(), &v1.Query{})
	s.NoError(err, "search should have succeeded")
	s.Equal([]search.Result{}, results)
}

func (s *idTransformationTestSuite) TestClusterScoped() {
	// Expect graph and search interactions.
	s.mockRGraph.EXPECT().GetRefsTo(prefixedID1).Return(toID1)
	s.mockRGraph.EXPECT().GetRefsTo(prefixedID2).Return(toID2)
	s.mockRGraph.EXPECT().GetRefsTo(prefixedID3).Return(toID3)
	s.mockRGraph.EXPECT().GetRefsTo(prefixedID4).Return(toID4)
	s.mockRGraph.EXPECT().GetRefsTo(prefixedID5).Return(toID5)
	s.mockRGraph.EXPECT().GetRefsTo(prefixedID6).Return(toID6)

	s.mockUnsafeSearcher.EXPECT().Search(gomock.Any()).Return([]search.Result{
		{
			ID: string(id1),
		},
		{
			ID: string(id2),
		},
	}, nil)

	checker := sac.NewScopeChecker(sac.OneStepSCC{
		sac.ClusterScopeKey("id7"): sac.AllowAllAccessScopeChecker(),
	})

	filter, err := NewSACFilter(
		WithScopeChecker(checker),
		WithGraphProvider(fakeGraphProvider{mg: s.mockRGraph}),
		WithClusterPath(prefix1, prefix2, prefix3, prefix4),
	)
	s.NoError(err, "filter creation should have succeeded")

	searcher := Searcher(s.mockUnsafeSearcher, filter)
	results, err := searcher.Search(context.Background(), &v1.Query{})
	s.NoError(err, "search should have succeeded")
	s.Equal([]search.Result{
		{
			ID: string(id1),
		},
		{
			ID: string(id2),
		},
	}, results)
}

func (s *idTransformationTestSuite) TestClusterScopedMultiCluster() {
	// Expect graph and search interactions.
	s.mockRGraph.EXPECT().GetRefsTo(prefixedID1).Return(toID1)
	s.mockRGraph.EXPECT().GetRefsTo(prefixedID2).Return(toID2)
	s.mockRGraph.EXPECT().GetRefsTo(prefixedID3).Return(toID3)
	s.mockRGraph.EXPECT().GetRefsTo(prefixedID4).Return(toID4)
	s.mockRGraph.EXPECT().GetRefsTo(prefixedID5).Return(toID5)
	s.mockRGraph.EXPECT().GetRefsTo(prefixedID6).Return(toID6)

	s.mockUnsafeSearcher.EXPECT().Search(gomock.Any()).Return([]search.Result{
		{
			ID: string(id1),
		},
		{
			ID: string(id2),
		},
	}, nil)

	// Allow first namespace and cluster
	checker := sac.NewScopeChecker(sac.OneStepSCC{
		sac.ClusterScopeKey("id7"): sac.DenyAllAccessScopeChecker(),
		sac.ClusterScopeKey("id8"): sac.AllowAllAccessScopeChecker(),
	})

	filter, err := NewSACFilter(
		WithScopeChecker(checker),
		WithGraphProvider(fakeGraphProvider{mg: s.mockRGraph}),
		WithClusterPath(prefix1, prefix2, prefix3, prefix4),
	)
	s.NoError(err, "filter creation should have succeeded")

	searcher := Searcher(s.mockUnsafeSearcher, filter)
	results, err := searcher.Search(context.Background(), &v1.Query{})
	s.NoError(err, "search should have succeeded")
	s.Equal([]search.Result{
		{
			ID: string(id1),
		},
	}, results)
}

func (s *idTransformationTestSuite) TestNamespaceScoped() {
	// Expect graph and search interactions.
	s.mockRGraph.EXPECT().GetRefsTo(prefixedID1).Return(toID1)
	s.mockRGraph.EXPECT().GetRefsTo(prefixedID2).Return(toID2)
	s.mockRGraph.EXPECT().GetRefsTo(prefixedID3).Return(toID3)
	s.mockRGraph.EXPECT().GetRefsTo(prefixedID4).Return(toID4)
	s.mockRGraph.EXPECT().GetRefsTo(prefixedID5).Return(toID5)
	s.mockRGraph.EXPECT().GetRefsTo(prefixedID6).Return(toID6)

	s.mockUnsafeSearcher.EXPECT().Search(gomock.Any()).Return([]search.Result{
		{
			ID: string(id1),
		},
		{
			ID: string(id2),
		},
	}, nil)

	// Allow first namespace and cluster
	checker := sac.NewScopeChecker(sac.OneStepSCC{
		sac.ClusterScopeKey("id7"): sac.OneStepSCC{
			sac.NamespaceScopeKey("id5"): sac.AllowAllAccessScopeChecker(),
			sac.NamespaceScopeKey("id6"): sac.DenyAllAccessScopeChecker(),
		},
	})

	filter, err := NewSACFilter(
		WithScopeChecker(checker),
		WithGraphProvider(fakeGraphProvider{mg: s.mockRGraph}),
		WithNamespacePath(prefix1, prefix2, prefix3),
		WithClusterPath(prefix1, prefix2, prefix3, prefix4),
	)
	s.NoError(err, "filter creation should have succeeded")

	searcher := Searcher(s.mockUnsafeSearcher, filter)
	results, err := searcher.Search(context.Background(), &v1.Query{})
	s.NoError(err, "search should have succeeded")
	s.Equal([]search.Result{
		{
			ID: string(id1),
		},
	}, results)
}

func (s *idTransformationTestSuite) TestNamespaceScopedMultiCluster() {
	// Expect graph and search interactions.
	s.mockRGraph.EXPECT().GetRefsTo(prefixedID1).Return(toID1)
	s.mockRGraph.EXPECT().GetRefsTo(prefixedID2).Return(toID2)
	s.mockRGraph.EXPECT().GetRefsTo(prefixedID3).Return(toID3)
	s.mockRGraph.EXPECT().GetRefsTo(prefixedID4).Return(toID4)
	s.mockRGraph.EXPECT().GetRefsTo(prefixedID5).Return(toID5)
	s.mockRGraph.EXPECT().GetRefsTo(prefixedID6).Return(toID6)

	s.mockUnsafeSearcher.EXPECT().Search(gomock.Any()).Return([]search.Result{
		{
			ID: string(id1),
		},
		{
			ID: string(id2),
		},
	}, nil)

	// Allow first namespace and cluster
	checker := sac.NewScopeChecker(sac.OneStepSCC{
		sac.ClusterScopeKey("id7"): sac.DenyAllAccessScopeChecker(),
		sac.ClusterScopeKey("id8"): sac.AllowAllAccessScopeChecker(),
	})

	filter, err := NewSACFilter(
		WithScopeChecker(checker),
		WithGraphProvider(fakeGraphProvider{mg: s.mockRGraph}),
		WithNamespacePath(prefix1, prefix2, prefix3),
		WithClusterPath(prefix1, prefix2, prefix3, prefix4),
	)
	s.NoError(err, "filter creation should have succeeded")

	searcher := Searcher(s.mockUnsafeSearcher, filter)
	results, err := searcher.Search(context.Background(), &v1.Query{})
	s.NoError(err, "search should have succeeded")
	s.Equal([]search.Result{
		{
			ID: string(id1),
		},
	}, results)
}

type fakeGraphProvider struct {
	mg *mocks.MockRGraph
}

func (fgp fakeGraphProvider) NewGraphView() graph.DiscardableRGraph {
	return graph.NewDiscardableGraph(fgp.mg, func() {})
}
