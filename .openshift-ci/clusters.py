#!/usr/bin/env python3

"""
Clusters used in test
"""

import os
import signal
import subprocess
import tempfile
import time

from common import popen_graceful_kill


class NullCluster:
    def provision(self):
        pass

    def teardown(self):
        pass


class GKECluster:
    # Provisioning timeout is tightly coupled to the time it may take gke.sh to
    # create a cluster.
    PROVISION_TIMEOUT = 140 * 60
    WAIT_TIMEOUT = 20 * 60
    TEARDOWN_TIMEOUT = 5 * 60
    # separate script names used for testability - test_clusters.py
    PROVISION_PATH = "scripts/ci/gke.sh"
    WAIT_PATH = "scripts/ci/gke.sh"
    REFRESH_PATH = "scripts/ci/gke.sh"
    TEARDOWN_PATH = "scripts/ci/gke.sh"

    def __init__(self, cluster_id, num_nodes=3, machine_type="e2-standard-4"):
        self.cluster_id = cluster_id
        self.num_nodes = num_nodes
        self.machine_type = machine_type
        self.refresh_token_cmd = None
        self.cluster_name = None

    def provision(self):
        with subprocess.Popen(
            [
                GKECluster.PROVISION_PATH,
                "provision_gke_cluster",
                self.cluster_id,
                str(self.num_nodes),
                self.machine_type,
            ]
        ) as cmd:

            try:
                exitstatus = cmd.wait(GKECluster.PROVISION_TIMEOUT)
                if exitstatus != 0:
                    raise RuntimeError(f"Cluster provision failed: exit {exitstatus}")
            except subprocess.TimeoutExpired as err:
                popen_graceful_kill(cmd)
                raise err

        # OpenShift CI sends a SIGINT when tests are canceled
        signal.signal(signal.SIGINT, self.sigint_handler)

        subprocess.run(
            [GKECluster.WAIT_PATH, "wait_for_cluster"],
            check=True,
            timeout=GKECluster.WAIT_TIMEOUT,
        )

        # pylint: disable=consider-using-with
        self.refresh_token_cmd = subprocess.Popen(
            [GKECluster.REFRESH_PATH, "refresh_gke_token"]
        )

        self.cluster_name = os.environ["CLUSTER_NAME"]

        return self

    def teardown(self):
        while os.path.exists("/tmp/hold-cluster"):
            print("Pausing teardown because /tmp/hold-cluster exists")
            time.sleep(60)

        if self.refresh_token_cmd is not None:
            print("Terminating GKE token refresh")
            try:
                popen_graceful_kill(self.refresh_token_cmd)
            except Exception as err:
                print(f"Could not terminate the token refresh: {err}")

        subprocess.run(
            [GKECluster.TEARDOWN_PATH, "teardown_gke_cluster", self.cluster_name],
            check=True,
            timeout=GKECluster.TEARDOWN_TIMEOUT,
        )

        return self

    def sigint_handler(self, signum, frame):
        print("Tearing down the cluster due to SIGINT", signum, frame)
        self.teardown()


class AutomationFlavorsCluster:
    KUBECTL_TIMEOUT = 5 * 60

    def provision(self):
        kubeconfig = os.environ["KUBECONFIG"]

        print(f"Using kubeconfig from {kubeconfig}")

        print("Nodes:")
        subprocess.run(
            ["kubectl", "get", "nodes", "-o", "wide"],
            check=True,
            timeout=AutomationFlavorsCluster.KUBECTL_TIMEOUT,
        )

        return self

    def teardown(self):
        pass

class OpenShiftScaleWorkersCluster:
    SCALE_CHANGE_TIMEOUT = 15 * 60

    def __init__(self, increment=1):
        self.increment = increment

    def provision(self):
        print("Scaling worker nodes")
        subprocess.run(
            ["scripts/ci/openshift.sh", "scale_worker_nodes", str(self.increment)],
            check=True,
            timeout=OpenShiftScaleWorkersCluster.SCALE_CHANGE_TIMEOUT,
        )

        return self

    def teardown(self):
        pass

class SeparateClusters:
    """
    SeparateClusters - central and sensor are deployed to separate clusters. If
    either of the two kubeconfig args are not passed a GKE cluster is created.
    """

    def __init__(self, cluster_id, central_cluster_kubeconfig="", sensor_cluster_kubeconfig=""):
        self.cluster_id = cluster_id
        self.central_cluster_kubeconfig = central_cluster_kubeconfig
        self.sensor_cluster_kubeconfig = sensor_cluster_kubeconfig
        self.central_cluster = None
        self.sensor_cluster = None

    def provision(self):
        if self.central_cluster_kubeconfig == "":
            kubeconfig = tempfile.NamedTemporaryFile(delete=False)
            kubeconfig.close()
            os.environ["KUBECONFIG"] = kubeconfig.name
            self.central_cluster = GKECluster(self.cluster_id + "-central")
            self.central_cluster.provision()
            os.environ["CENTRAL_CLUSTER_KUBECONFIG"] = kubeconfig.name
            self.central_cluster_kubeconfig = kubeconfig.name

        if self.sensor_cluster_kubeconfig == "":
            kubeconfig = tempfile.NamedTemporaryFile(delete=False)
            kubeconfig.close()
            os.environ["KUBECONFIG"] = kubeconfig.name
            self.sensor_cluster = GKECluster(self.cluster_id + "-central")
            self.sensor_cluster.provision()
            os.environ["SENSOR_CLUSTER_KUBECONFIG"] = kubeconfig.name
            self.sensor_cluster_kubeconfig = kubeconfig.name

        return self

    def teardown(self):
        if self.central_cluster is not None:
            self.central_cluster.teardown()

        if self.sensor_cluster is not None:
            self.sensor_cluster.teardown()

        return self
