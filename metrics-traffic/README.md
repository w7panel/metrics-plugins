# Traffic monitoring plugin

`frontend/` is migrated from `w7panel-traffic`. `backend/` serves the seven traffic read APIs. The Chart owns Higress Scrape, Vector access-log collection, workload snapshots, and the query API. VictoriaLogs and its PVC remain in `w7panel-metrics`.

Run `make package VERSION=x.y.z` for local release inputs, or log Docker in to `ccr.ccs.tencentyun.com`, set `ZPK_IDC_PASSWORD`, and run `make full VERSION=x.y.z` to build the UI and Helm package, push the API image, and publish both attachments to the `metrics-traffic` ZPK application. `make build-ui`, `make build-helm`, `make build-image`, and `make push-zpk` run the individual stages. Configure the MicroApp backend URL to the Chart Service. The query API derives namespace access from the verified Panel identity; client-supplied `namespace` is ignored for normal users.
