# Audit log plugin

`frontend/` is the standalone audit MicroApp. `backend/` serves the three audit read APIs. The Panel server retains the login and operation event writer; it writes to the shared VictoriaLogs instance. The Chart deploys only the query API Service and Deployment.

Run `make package VERSION=x.y.z` to build the local release inputs. Configure the MicroApp backend URL to the Chart Service and attach `frontend/dist.zip` plus the Helm package to the `metrics-audit` ZPK application. The query API requires a valid Panel identity and uses server-side tenant and username filters.
