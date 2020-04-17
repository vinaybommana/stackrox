package util

import orchestratormanager.OrchestratorTypes

import java.nio.file.Files
import java.nio.file.Paths
import java.nio.file.attribute.FileTime

class Env {

    private static final PROPERTIES_FILE = "qa-test-settings.properties"

    private static final DEFAULT_VALUES = [
            "API_HOSTNAME": "localhost",
            "API_PORT": "8000",
            "ROX_USERNAME": "admin",
    ]

    static final IN_CI = (System.getenv("CI") != null)

    private static final Env INSTANCE = new Env()

    static String get(String key, String defVal = null) {
        return INSTANCE.getInternal(key, defVal)
    }

    static String mustGet(String key) {
        return INSTANCE.mustGetInternal(key)
    }

    static String mustGetInCI(String key) {
        return INSTANCE.mustGetInCIInternal(key)
    }

    private final envVars = new Properties()

    private Env() {
        if (!IN_CI) {
            loadEnvVarsFromPropsFile()
        }
        envVars.putAll(System.getenv())
        if (!IN_CI) {
            assignFallbackValues()
        }
    }

    private loadEnvVarsFromPropsFile() {
        try {
            envVars.load(new FileInputStream(PROPERTIES_FILE))
        } catch (Exception ex) {
            print "Failed to load extra properties file: ${ex.toString()}"
        }
    }

    protected String getInternal(String key, String defVal) {
        return envVars.getOrDefault(key, defVal)
    }

    protected String mustGetInternal(String key) {
        def value = envVars.get(key)
        if (value == null) {
            throw new RuntimeException("No value assigned for required key ${key}")
        }
        return value
    }

    protected String mustGetInCIInternal(String key) {
        def value = envVars.get(key)
        if (value == null) {
            if (inCI) {
                throw new RuntimeException("No value assigned for required key ${key}")
            }
        }
        return value
    }

    private static OrchestratorTypes inferOrchestratorType() {
        // Infer the orchestrator type from the local deployment, by looking at which
        // `deploy/<orchestrator>/central-deploy/password` file was most recently written to.
        OrchestratorTypes selected = null
        FileTime mostRecent = null
        for (def orchestratorType : OrchestratorTypes.values()) {
            def passwordPath = "../deploy/${orchestratorType.toString().toLowerCase()}/central-deploy/password"
            try {
                def modTime = Files.getLastModifiedTime(Paths.get(passwordPath))
                if (mostRecent == null || modTime > mostRecent) {
                    selected = orchestratorType
                    mostRecent = modTime
                }
            } catch (Exception ex) {
                print "" // no-op
            }
        }

        return selected ?: OrchestratorTypes.K8S
    }

    private void assignFallbackValues() {
        for (def entry : DEFAULT_VALUES.entrySet()) {
            if (envVars.get(entry.key) == null) {
                envVars.put(entry.key, entry.value)
            }
        }

        if (envVars.get("ROX_PASSWORD") == null) {
            if (envVars.get("CLUSTER") == null) {
                envVars.put("CLUSTER", inferOrchestratorType().toString())
            }

            String password = null
            try {
                def passwordPath = "../deploy/${envVars.get("CLUSTER").toLowerCase()}/central-deploy/password"
                BufferedReader br = new BufferedReader(new FileReader(passwordPath))
                password = br.readLine()
            } catch (Exception ex) {
                println "Failed to load password for current deployment: ${ex.toString()}"
            }

            if (password != null) {
                envVars.put("ROX_PASSWORD", password)
            }
        }

        if (envVars.get("CLUSTER") == null) {
            envVars.put("CLUSTER", OrchestratorTypes.K8S.toString())
        }
    }

    static String mustGetUsername() {
        return mustGet("ROX_USERNAME")
    }

    static String mustGetPassword() {
        return mustGet("ROX_PASSWORD")
    }

    static int mustGetPort() {
        String portString = mustGet("API_PORT")
        int port
        try {
            port = Integer.parseInt(portString)
        } catch (NumberFormatException e) {
            throw new RuntimeException("API_PORT " + portString + " is not a valid number " + e.toString())
        }
        return port
    }

    static String mustGetHostname() {
        return mustGet("API_HOSTNAME")
    }

    static OrchestratorTypes mustGetOrchestratorType() {
        String cluster = mustGet("CLUSTER")
        OrchestratorTypes type
        try {
            type = OrchestratorTypes.valueOf(cluster)
        } catch (IllegalArgumentException e) {
            throw new RuntimeException("CLUSTER must be one of " + OrchestratorTypes.values() + " " + e.toString())
        }
        return type
    }

    static String mustGetKeystorePath() {
        return mustGet("KEYSTORE_PATH")
    }

    static String mustGetTruststorePath() {
        return mustGet("TRUSTSTORE_PATH")
    }

    static String mustGetClientCAPath() {
        return mustGet("CLIENT_CA_PATH")
    }

    static String mustGetCiJobName() {
        return mustGet("CI_JOB_NAME")
    }

    static String mustGetImageTag() {
        return mustGet("IMAGE_TAG")
    }

    static String mustGetResultsFilePath() {
        return mustGet("RESULTS_FILE_PATH")
    }

    static String mustGetTestRailPassword() {
        return mustGet("TESTRAIL_PASSWORD")
    }

    static String mustGetDockerIOUserName() {
        return mustGet("DOCKER_IO_PULL_USERNAME")
    }

    static String mustGetDockerIOPassword() {
        return mustGet("DOCKER_IO_PULL_PASSWORD")
    }
}

