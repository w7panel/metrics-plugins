import { panelGet, trafficGet } from './runtime'

export const trafficApi = {
  health: () => trafficGet('health'),
  summary: (params) => trafficGet('summary', params),
  series: (params) => trafficGet('series', params),
  apps: (params) => trafficGet('apps', params),
  pods: (params) => trafficGet('pods', params),
  domains: (params) => trafficGet('domains', params),
  urls: (params) => trafficGet('urls', params),
  namespaces: () => panelGet('namespaces'),
}
