# Audit log plugin

`frontend/` is the standalone audit MicroApp. `backend/` serves the three audit read APIs. The Panel server retains the login and operation event writer; it writes to the shared VictoriaLogs instance. The Chart deploys only the query API Service and Deployment.

Run `make package VERSION=x.y.z` for local release inputs, or log Docker in to `ccr.ccs.tencentyun.com`, set `ZPK_IDC_PASSWORD`, and run `make full VERSION=x.y.z` to build the UI and Helm package, push the API image, and publish both attachments to the `metrics-audit` ZPK application. `make build-ui`, `make build-helm`, `make build-image`, and `make push-zpk` run the individual stages. Configure the MicroApp backend URL to the Chart Service. The query API requires a valid Panel identity and uses server-side tenant and username filters.
