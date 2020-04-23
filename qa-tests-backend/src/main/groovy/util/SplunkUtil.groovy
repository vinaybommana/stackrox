package util

import static com.jayway.restassured.RestAssured.given

import com.google.gson.GsonBuilder
import objects.SplunkAlert
import objects.SplunkAlertRaw
import objects.SplunkAlerts
import objects.SplunkSearch

import com.google.gson.Gson
import com.jayway.restassured.response.Response

class SplunkUtil {
    static final private Gson GSON = new GsonBuilder().create()

    static List<SplunkAlert> getSplunkAlerts(String httpsEndPoint, String searchId) {
        Response response = getSearchResults(httpsEndPoint, searchId)
        SplunkAlerts alerts = GSON.fromJson(response.asString(), SplunkAlerts)

        def returnAlerts = []
        for (SplunkAlertRaw raw : alerts.results) {
            returnAlerts.add(GSON.fromJson(raw._raw, SplunkAlert))
        }
        return returnAlerts
    }

    static List<SplunkAlert> waitForSplunkAlerts(String httpsLoadBalancer, int timeoutSeconds) {
        int intervalSeconds = 3
        int iterations = timeoutSeconds / intervalSeconds
        List results = []
        Timer t = new Timer(iterations, intervalSeconds)
        while (results.size() == 0 && t.IsValid()) {
            def searchId = createSearch(httpsLoadBalancer)
            results = getSplunkAlerts(httpsLoadBalancer, searchId)
        }
        return results
    }

    static Response getSearchResults(String deploymentIP, String searchId) {
        Response response
        try {
            response = given().auth().basic("admin", "changeme")
                    .param("output_mode", "json")
                    .get("https://${deploymentIP}:8089/services/search/jobs/${searchId}/events")
            println("Querying loadbalancer ${deploymentIP}")
        }
        catch (Exception e) {
            println("catching unknownhost exception for KOPS and other intermittent connection issues" + e)
        }
        println "Printing response from ${deploymentIP} " + response?.prettyPrint()
        return response
    }

    static String createSearch(String deploymentIP) {
        Response response
        try {
            response = given().auth().basic("admin", "changeme")
                    .formParam("search", "search")
                    .param("output_mode", "json")
                    .post("https://${deploymentIP}:8089/services/search/jobs")
        }
        catch (Exception e) {
            println("catching unknownhost exception for KOPS and other intermittent connection issues" + e)
        }
        return GSON.fromJson(response.asString(), SplunkSearch).sid
    }
}
