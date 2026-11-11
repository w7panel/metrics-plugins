# W7Panel monitoring plugins

This repository contains two independently packaged monitoring applications:

| Directory | UI | API | Collection |
| --- | --- | --- | --- |
| `metrics-audit` | login and operation audit logs | VictoriaLogs audit queries | Panel continues to write audit events |
| `metrics-traffic` | traffic dashboard | VictoriaLogs traffic queries | Higress `VMPodScrape`, Vector access logs and workload snapshots |

Both use the existing `w7panel-metrics` VictoriaLogs service and PVC. Install `w7panel-metrics` first. The plugins default to the `default` namespace and the existing service names; override the Chart URLs for other installations.

Each directory is an independent Go module and Docker build context. Run `make package VERSION=x.y.z` inside either directory to produce a frontend `dist.zip`, Helm `.tgz`, and local API image. For a release, log Docker in to `ccr.ccs.tencentyun.com`, set `ZPK_IDC_PASSWORD` (or `PASSWORD`), and run `make full VERSION=x.y.z` inside each directory. The stages are `build-ui`, `build-helm`, `build-image` (build and push the API image), and `push-zpk` (attach the frontend zip and Helm package). Override `IMAGE_REPO` if publishing to another registry; set the Chart's `image.repository` to the same value. Configure each MicroApp backend URL to `http://metrics-audit.default.svc:8080` or `http://metrics-traffic.default.svc:8080`. Frontends prefer the host `microappProxy`; older hosts use the preserved `/panel-api/v1/audit/*` or `/panel-api/v1/traffic/*` compatibility routes.

The plugin API accepts a Panel token only after the host validates it via `/panel-api/v1/metrics/plugin-identity`. The API services are intended for cluster-internal access. The traffic API ServiceAccount may list Pods to preserve the existing Pod name fallback, while query filters still enforce the verified user's namespace.

## Upgrade order

1. Build and publish both API images, frontend archives, and Charts. Set the Chart image repositories to the published images. Install `metrics-audit`; install `metrics-traffic` with `higressAccessLogs.enabled=false` and `higressMetrics.enabled=false` so its query API is ready without starting a second collector.
2. Upgrade the Panel server to provide the identity endpoint and proxy the old read URLs to these APIs. Verify audit and traffic queries, including a normal user's namespace restrictions.
3. Upgrade `w7panel-metrics` to remove its Higress Scrape and Vector resources. Then enable the two collector switches in `metrics-traffic`. The old Vector hostPath checkpoint is retained to avoid replaying logs; this short handoff can leave a collection gap.
4. Install both MicroApp frontends and upgrade the main Panel UI to remove its native audit page. Verify the new menus and old read URLs before considering the rollout complete.

Do not run the old and new Vector DaemonSets together: both would ingest the same access logs. Uninstalling either plugin does not delete VictoriaLogs or its PVC.
