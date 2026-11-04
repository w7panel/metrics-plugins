# Traffic monitoring plugin

`frontend/` is migrated from `w7panel-traffic`. `backend/` serves the seven traffic read APIs. The Chart owns Higress Scrape, Vector access-log collection, workload snapshots, and the query API. VictoriaLogs and its PVC remain in `w7panel-metrics`.

Run `make package VERSION=x.y.z` to build the local release inputs. Configure the MicroApp backend URL to the Chart Service and attach `frontend/dist.zip` plus the Helm package to the `metrics-traffic` ZPK application. The query API derives namespace access from the verified Panel identity; client-supplied `namespace` is ignored for normal users.
