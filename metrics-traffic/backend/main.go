package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/w7panel/metrics-plugins/metrics-traffic/internal/identity"
)

func trafficHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user, err := identity.Verify(r)
	if err != nil {
		identity.Error(w, http.StatusUnauthorized, err)
		return
	}
	client := NewQueryClient()
	if r.URL.Path == "/api/v1/traffic/health" {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		err := client.Health(ctx)
		message := ""
		if err != nil {
			message = err.Error()
		}
		identity.JSON(w, http.StatusOK, map[string]any{"logs": map[string]any{"available": err == nil, "message": message}, "hubble": map[string]any{"available": metricHasData(ctx, "hubble_flows_processed_total")}, "higressMetrics": map[string]any{"available": metricHasData(ctx, "envoy_cluster_upstream_rq_total")}})
		return
	}
	params, err := parseTrafficParams(r, user)
	if err != nil {
		identity.Error(w, http.StatusBadRequest, err)
		return
	}
	var rows []map[string]any
	switch r.URL.Path {
	case "/api/v1/traffic/summary":
		rows, err = client.Summary(r.Context(), params)
	case "/api/v1/traffic/series":
		rows, err = client.Series(r.Context(), params)
	case "/api/v1/traffic/apps":
		rows, err = client.Apps(r.Context(), params)
	case "/api/v1/traffic/pods":
		rows, err = client.Pods(r.Context(), params)
	case "/api/v1/traffic/domains":
		rows, err = client.Domains(r.Context(), params)
	case "/api/v1/traffic/urls":
		rows, err = client.URLs(r.Context(), params)
	default:
		http.NotFound(w, r)
		return
	}
	if err != nil {
		identity.Error(w, http.StatusBadGateway, err)
		return
	}
	switch r.URL.Path {
	case "/api/v1/traffic/summary":
		folded := FoldRows(rows)
		if len(folded) == 0 {
			identity.JSON(w, http.StatusOK, map[string]any{})
		} else {
			identity.JSON(w, http.StatusOK, folded[0])
		}
	case "/api/v1/traffic/series":
		identity.JSON(w, http.StatusOK, rows)
	default:
		keys := map[string][]string{
			"/api/v1/traffic/apps":    {"workload_name", "workload_kind", "workload_title", "workload_namespace"},
			"/api/v1/traffic/pods":    {"upstream_ip", "upstream_pod_name", "upstream_service", "workload_title", "upstream_namespace"},
			"/api/v1/traffic/domains": {"authority"},
			"/api/v1/traffic/urls":    {"authority", "method", "path"},
		}
		rows = FoldRows(rows, keys[r.URL.Path]...)
		SortRows(rows, params.Sort)
		switch r.URL.Path {
		case "/api/v1/traffic/apps":
			rows = SearchRows(rows, params.Search, "workload_name", "workload_kind", "workload_namespace")
		case "/api/v1/traffic/pods":
			resolvePodNames(r.Context(), params.Namespace, rows)
			rows = SearchRows(rows, params.Search, "pod_name", "upstream_pod_name", "upstream_ip", "upstream_service", "upstream_namespace")
		case "/api/v1/traffic/domains":
			rows = SearchRows(rows, params.Search, "authority")
		}
		list, total := Paginate(rows, params.Page, params.PageSize)
		identity.JSON(w, http.StatusOK, map[string]any{"list": list, "total": total, "page": params.Page, "pageSize": params.PageSize})
	}
}

func parseTrafficParams(r *http.Request, user identity.User) (QueryParams, error) {
	q := r.URL.Query()
	timerange, err := ParseTimeRange(q.Get("start"), q.Get("end"), time.Now())
	if err != nil {
		return QueryParams{}, err
	}
	namespace := user.Tenant
	if user.IsAdmin {
		namespace = strings.TrimSpace(q.Get("namespace"))
		if namespace == "" {
			namespace = user.Tenant
		}
	} else {
		namespace = "k3k-" + user.Username
	}
	if namespace == "" {
		namespace = "default"
	}
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("pageSize"))
	params := QueryParams{Namespace: namespace, Domain: strings.TrimSpace(q.Get("domain")), UpstreamIP: strings.TrimSpace(q.Get("upstreamIp")), WorkloadName: strings.TrimSpace(q.Get("workloadName")), WorkloadKind: strings.TrimSpace(q.Get("workloadKind")), Method: strings.TrimSpace(q.Get("method")), Status: strings.TrimSpace(q.Get("status")), Keyword: strings.TrimSpace(q.Get("keyword")), Search: strings.TrimSpace(q.Get("search")), Sort: strings.TrimSpace(q.Get("sort")), Step: strings.TrimSpace(q.Get("step")), Page: page, PageSize: pageSize, Range: timerange}
	NormalizeParams(&params)
	return params, nil
}

func metricHasData(ctx context.Context, query string) bool {
	base := strings.TrimRight(os.Getenv("METRICS_BASE_URL"), "/")
	if base == "" {
		return false
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base+"/api/v1/query?query="+url.QueryEscape(query), nil)
	if err != nil {
		return false
	}
	resp, err := (&http.Client{Timeout: 4 * time.Second}).Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	var result struct {
		Data struct {
			Result []json.RawMessage `json:"result"`
		} `json:"data"`
	}
	return resp.StatusCode == http.StatusOK && json.NewDecoder(resp.Body).Decode(&result) == nil && len(result.Data.Result) > 0
}

func resolvePodNames(ctx context.Context, namespace string, rows []map[string]any) {
	byIP := map[string][2]string{}
	if host := os.Getenv("KUBERNETES_SERVICE_HOST"); host != "" {
		ca, caErr := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")
		token, tokenErr := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
		if caErr == nil && tokenErr == nil {
			roots := x509.NewCertPool()
			if roots.AppendCertsFromPEM(ca) {
				path := "/api/v1/pods"
				if namespace != "*" {
					path = "/api/v1/namespaces/" + url.PathEscape(namespace) + "/pods"
				}
				req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://"+host+":"+defaultPort(os.Getenv("KUBERNETES_SERVICE_PORT"))+path, nil)
				if err == nil {
					req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(string(token)))
					resp, err := (&http.Client{Timeout: 5 * time.Second, Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: roots}}}).Do(req)
					if err == nil {
						var pods struct {
							Items []struct {
								Metadata struct {
									Name      string `json:"name"`
									Namespace string `json:"namespace"`
								} `json:"metadata"`
								Status struct {
									PodIP string `json:"podIP"`
								} `json:"status"`
							} `json:"items"`
						}
						if resp.StatusCode == http.StatusOK {
							_ = json.NewDecoder(resp.Body).Decode(&pods)
						}
						_ = resp.Body.Close()
						for _, pod := range pods.Items {
							if pod.Status.PodIP != "" {
								byIP[pod.Status.PodIP] = [2]string{pod.Metadata.Name, pod.Metadata.Namespace}
							}
						}
					}
				}
			}
		}
	}
	for _, row := range rows {
		ip, _ := row["upstream_ip"].(string)
		if name := strings.TrimSpace(fmt.Sprint(row["upstream_pod_name"])); name != "" && name != "<nil>" {
			row["pod_name"] = name
			row["namespace"] = row["upstream_namespace"]
		} else if pod, ok := byIP[ip]; ok {
			row["pod_name"] = pod[0]
			row["namespace"] = pod[1]
		} else {
			row["pod_name"] = "已终止 Pod"
			row["namespace"] = row["upstream_namespace"]
		}
	}
}

func defaultPort(port string) string {
	if port == "" {
		return "443"
	}
	return port
}

func main() {
	if os.Getenv("LOGS_BASE_URL") == "" {
		panic("LOGS_BASE_URL is required")
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/traffic/", trafficHandler)
	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	if err := http.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}
}
