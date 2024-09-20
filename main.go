package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    containerCPUUsage = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "container_cpu_usage_seconds_total",
            Help: "Total CPU usage of containers",
        },
        []string{"container_name"},
    )
    containerMemoryUsage = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "container_memory_usage_bytes",
            Help: "Memory usage of containers",
        },
        []string{"container_name"},
    )
    webappVisitorCount = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "webapp_visitor_count",
            Help: "Number of visitors to the webapp",
        },
        []string{"container_name"},
    )
)

func init() {
    prometheus.MustRegister(containerCPUUsage)
    prometheus.MustRegister(containerMemoryUsage)
    prometheus.MustRegister(webappVisitorCount)
}

func collectMetrics() {
    url := "http://localhost:2375/containers/json" // Docker API URL
    resp, err := http.Get(url)
    if err != nil {
        log.Printf("Error fetching container list: %v", err)
        return
    }
    defer resp.Body.Close()

    var containers []map[string]interface{}
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading response body: %v", err)
        return
    }
    if err := json.Unmarshal(body, &containers); err != nil {
        log.Printf("Error parsing container list JSON: %v", err)
        return
    }

    // Filter containers by names "webapp1" and "webapp2"
    targetContainers := map[string]bool{
        "webapp1": true,
        "webapp2": true,
    }

    for _, container := range containers {
        containerName := container["Names"].([]interface{})[0].(string)
        if containerName[0] == '/' {
            containerName = containerName[1:]
        }

        if !targetContainers[containerName] {
            continue
        }

        statsURL := fmt.Sprintf("http://localhost:2375/containers/%s/stats?stream=false", container["Id"])
        statsResp, err := http.Get(statsURL)
        if err != nil {
            log.Printf("Error fetching stats for container %s: %v", containerName, err)
            continue
        }
        defer statsResp.Body.Close()

        var stats map[string]interface{}
        statsBody, err := ioutil.ReadAll(statsResp.Body)
        if err != nil {
            log.Printf("Error reading stats response body for container %s: %v", containerName, err)
            continue
        }
        if err := json.Unmarshal(statsBody, &stats); err != nil {
            log.Printf("Error parsing stats JSON for container %s: %v", containerName, err)
            continue
        }

        cpuUsage := 0.0
        memoryUsage := 0.0

        if cpuStats, ok := stats["cpu_stats"].(map[string]interface{}); ok {
            if cpuUsageSecs, ok := cpuStats["cpu_usage"].(map[string]interface{})["total_usage"].(float64); ok {
                cpuUsage = cpuUsageSecs / 1e9 // Convert nanoseconds to seconds
            }
        }

        if memoryStats, ok := stats["memory_stats"].(map[string]interface{}); ok {
            if usage, ok := memoryStats["usage"].(float64); ok {
                memoryUsage = usage
            }
        }

        containerCPUUsage.WithLabelValues(containerName).Set(cpuUsage)
        containerMemoryUsage.WithLabelValues(containerName).Set(memoryUsage)
        log.Printf("Stats for container %s: CPU %f, Memory %f", containerName, cpuUsage, memoryUsage)
    }

    // Collect visitor counts from the /count endpoints
    for _, containerName := range []string{"webapp1", "webapp2"} {
        countURL := fmt.Sprintf("http://localhost:%d/count", getPort(containerName))
        countResp, err := http.Get(countURL)
        if err != nil {
            log.Printf("Error fetching visitor count for container %s: %v", containerName, err)
            continue
        }
        defer countResp.Body.Close()

        var countData map[string]float64
        countBody, err := ioutil.ReadAll(countResp.Body)
        if err != nil {
            log.Printf("Error reading count response body for container %s: %v", containerName, err)
            continue
        }
        if err := json.Unmarshal(countBody, &countData); err != nil {
            log.Printf("Error parsing count JSON for container %s: %v", containerName, err)
            continue
        }

        visitorCount, ok := countData["visitor_count"]
        if !ok {
            log.Printf("Visitor count not found for container %s", containerName)
            continue
        }

        webappVisitorCount.WithLabelValues(containerName).Set(visitorCount)
        log.Printf("Visitor count for container %s: %f", containerName, visitorCount)
    }
}

// Helper function to get the port based on container name
func getPort(containerName string) int {
    switch containerName {
    case "webapp1":
        return 8081
    case "webapp2":
        return 8082
    default:
        return 8080
    }
}

func main() {
    go func() {
        for {
            collectMetrics()
            time.Sleep(15 * time.Second) // Adjust the interval as needed
        }
    }()

    http.Handle("/metrics", promhttp.Handler())
    log.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

