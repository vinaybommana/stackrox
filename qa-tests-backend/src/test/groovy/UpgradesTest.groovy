import com.google.protobuf.UnknownFieldSet
import groups.Upgrade
import com.google.protobuf.Timestamp
import io.stackrox.proto.api.v1.AlertServiceOuterClass
import io.stackrox.proto.api.v1.DeploymentServiceOuterClass
import io.stackrox.proto.api.v1.SearchServiceOuterClass
import io.stackrox.proto.api.v1.SummaryServiceOuterClass
import io.stackrox.proto.storage.ClusterOuterClass
import io.stackrox.proto.storage.ProcessIndicatorOuterClass
import org.junit.Assume
import org.junit.experimental.categories.Category
import services.AlertService
import services.ClusterService
import services.ConfigService
import services.DeploymentService
import services.ImageService
import services.PolicyService
import services.ProcessService
import services.SecretService
import services.SummaryService
import util.Env

class UpgradesTest extends BaseSpecification {
    private final static String CLUSTERID = Env.mustGet("UPGRADE_CLUSTER_ID")

    @Category(Upgrade)
    def "Verify cluster exists and that field values are retained"() {
        given:
        "Only run on specific upgrade from 2.4.16"
        Assume.assumeTrue(CLUSTERID=="260e11a3-cbea-464c-95f0-588fa7695b49")

        expect:
        def clusters = ClusterService.getClusters()
        clusters.size() == 1
        def expectedCluster = ClusterOuterClass.Cluster.newBuilder()
                .setId(CLUSTERID)
                .setName("remote")
                .setType(ClusterOuterClass.ClusterType.KUBERNETES_CLUSTER)
                .setPriority(1)
                .setMainImage("stackrox/main:2.4.16.4")
                .setCentralApiEndpoint("central.stackrox:443")
                .setCollectionMethod(ClusterOuterClass.CollectionMethod.KERNEL_MODULE)
                .setRuntimeSupport(true)
                .setTolerationsConfig(ClusterOuterClass.TolerationsConfig.newBuilder()
                        .setDisabled(true)
                        .build())
                .setStatus(ClusterOuterClass.ClusterStatus.newBuilder()
                        .setLastContact(Timestamp.newBuilder().setSeconds(1551412107).setNanos(857477786).build())
                        .setProviderMetadata(ClusterOuterClass.ProviderMetadata.newBuilder()
                                .setGoogle(ClusterOuterClass.GoogleProviderMetadata.newBuilder()
                                        .setProject("ultra-current-825")
                                        .setClusterName("setup-devde6c6")
                                        .build())
                                .setRegion("us-west1")
                                .setZone("us-west1-c")
                                .build())
                        .setOrchestratorMetadata(ClusterOuterClass.OrchestratorMetadata.newBuilder()
                                .setVersion("v1.11.7-gke.4")
                                .setBuildDate(Timestamp.newBuilder().setSeconds(1549394549).build())
                                .build())
                        .build())
                .build()

        def cluster = ClusterOuterClass.Cluster.newBuilder(clusters.get(0))
                .setUnknownFields(UnknownFieldSet.defaultInstance)
                .build()
        cluster == expectedCluster
    }

    @Category(Upgrade)
    def "Verify process indicators have cluster IDs and namespaces added"() {
        given:
        "Only run on specific upgrade from 2.4.16"
        Assume.assumeTrue(CLUSTERID=="260e11a3-cbea-464c-95f0-588fa7695b49")

        expect:
        "Migrated ProcessIndicators to have a cluster ID and a namespace"
        def processIndicators = ProcessService.getProcessIndicatorsByDeployment("33b3eb66-3bd4-11e9-b563-42010a8a0101")
        processIndicators.size() > 0
        for (ProcessIndicatorOuterClass.ProcessIndicator indicator : processIndicators) {
            assert(indicator.getClusterId() == CLUSTERID)
            assert(indicator.getNamespace() != "")
        }
    }

    @Category(Upgrade)
    def "Verify private config contains the correct retention duration for alerts and images"() {
        given:
        "Only run on specific upgrade from 2.4.16"
        Assume.assumeTrue(CLUSTERID=="260e11a3-cbea-464c-95f0-588fa7695b49")

        expect:
        "Alert retention duration is nil, image rentention duration is 7 days"
        def config = ConfigService.getConfig()
        config != null
        config.getPrivateConfig().getAlertConfig() != null
        config.getPrivateConfig().getAlertConfig().getAllRuntimeRetentionDurationDays() == 0
        config.getPrivateConfig().getAlertConfig().getResolvedDeployRetentionDurationDays() == 0
        config.getPrivateConfig().getAlertConfig().getDeletedRuntimeRetentionDurationDays() == 0
        config.getPrivateConfig().getImageRetentionDurationDays() == 7
    }

    @Category(Upgrade)
    def "Verify that deployments are searchable post upgrade"() {
        expect:
        "Deployments should be searchable after the upgrade"
        DeploymentServiceOuterClass.ListDeploymentsResponse resp = DeploymentService.listDeploymentsSearch(
                SearchServiceOuterClass.RawQuery.newBuilder().setQuery("Cluster ID:${CLUSTERID}").build())
        assert resp.deploymentsList.size() > 0
    }

    @Category(Upgrade)
    def "Verify that images are searchable post upgrade"() {
        expect:
        "Images should be searchable after the upgrade"
        def imageList = ImageService.getImages(
                SearchServiceOuterClass.RawQuery.newBuilder().setQuery("Cluster ID:${CLUSTERID}").build()
        )
        assert imageList.size() > 0
    }

    @Category(Upgrade)
    def "Verify that alerts are searchable post upgrade"() {
        expect:
        "Alerts should be searchable after the upgrade"
        def alertList = AlertService.getViolations(
                AlertServiceOuterClass.ListAlertsRequest.newBuilder().setQuery("Cluster ID:${CLUSTERID}").build())
        assert alertList.size() > 0
    }

    @Category(Upgrade)
    def "Verify that secrets are searchable post upgrade"() {
        expect:
        "Secrets should be searchable after the upgrade"
        def secretList = SecretService.getSecrets(
                SearchServiceOuterClass.RawQuery.newBuilder().setQuery("Cluster ID:${CLUSTERID}").build()
        )
        assert secretList.size() > 0
    }

    @Category(Upgrade)
    def "Verify that policies are searchable post upgrade"() {
        expect:
        "Policies should be searchable after the upgrade"
        def policyList = PolicyService.getPolicies(
                SearchServiceOuterClass.RawQuery.newBuilder().setQuery("Policy:Latest Tag").build()
        )
        assert policyList.size() > 0
    }

    @Category(Upgrade)
    def "Verify that summary API returns non-zero values on upgrade"() {
        expect:
        "Summary API returns non-zero values on upgrade"
        SummaryServiceOuterClass.SummaryCountsResponse resp = SummaryService.getCounts()
        assert resp.numAlerts != 0
        assert resp.numDeployments != 0
        assert resp.numSecrets != 0
        assert resp.numClusters != 0
        assert resp.numImages != 0
        assert resp.numNodes != 0
    }

    // TODO
    // network flow edges
    // compliance
    // clairify integration
    // slack integration
}
