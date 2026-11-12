# Traffic monitoring plugin

`ui/` is migrated from `w7panel-traffic`. `backend/` serves the seven traffic read APIs. The Chart owns Higress Scrape, Vector access-log collection, workload snapshots, and the query API. VictoriaLogs and its PVC remain in `w7panel-metrics`.

Run `make package VERSION=x.y.z` for local release inputs, or log Docker in to `ccr.ccs.tencentyun.com`, set `ZPK_IDC_PASSWORD`, and run `make full VERSION=x.y.z` to build the UI and Helm package, push the API image, and publish both attachments to the `metrics-traffic` ZPK application. `make build-ui`, `make build-helm`, `make build-image`, and `make push-zpk` run the individual stages. With a MicroApp backend URL, requests use the Chart Service; without one, they use the Panel traffic proxy route. The query API derives namespace access from the verified Panel identity; client-supplied `namespace` is ignored for normal users.

## 本地连接真实集群

在 `ui/` 执行 `npm ci`、`npm run dev`，打开 `http://172.16.1.3:8014/`。首次查询时输入当前 W7Panel 的 Panel Token；刷新页面可重新输入。Vite 将 `/panel-api` 请求转发到 `http://218.23.2.48:9090`，由 Panel 按当前用户身份访问集群中的插件 API。如本机局域网 IP 变化，设置 `METRICS_DEV_ORIGIN=http://<本机局域网IP>:8014`；`METRICS_DEV_PORT`、`METRICS_PANEL_ORIGIN` 可覆盖默认端口和 Panel 地址。测试前需确认 `w7panel-server` 的流量代理、`metrics-traffic` 后端及 VictoriaLogs 已运行。`npm run build` 和在本项目根目录执行 `go test ./...` 可做本地静态验证。
