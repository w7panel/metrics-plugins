const isWujie = () => Boolean(window.__POWERED_BY_WUJIE__)

export function microProps() {
  return isWujie() ? (window.$wujie?.props || {}) : (window.microData?.cd || {})
}

export function isMicroapp() {
  return isWujie() || Boolean(window.microData?.cd)
}

export function userRole() {
  return String(microProps().role || 'normal').toLowerCase()
}

async function readPayload(response) {
  const payload = await response.json().catch(() => null)
  if (!response.ok) throw new Error(payload?.msg || payload?.message || `请求失败 (${response.status})`)
  return payload?.data ?? payload
}

export async function panelGet(path, params = {}) {
  const query = new URLSearchParams()
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') query.set(key, String(value))
  })
  const suffix = query.size ? `?${query}` : ''
  const props = microProps()
  if (typeof props.panelProxy === 'function') {
    return readPayload(await props.panelProxy(path.replace(/^\//, '') + suffix))
  }
  if (isMicroapp()) throw new Error('主面板未注入 panelProxy，请升级主面板后重试')
  return readPayload(await fetch(`/panel-api/v1/${path.replace(/^\//, '')}${suffix}`, { credentials: 'same-origin' }))
}

export async function trafficGet(path, params = {}) {
  const query = new URLSearchParams()
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') query.set(key, String(value))
  })
  const suffix = query.size ? `?${query}` : ''
  const proxy = microProps().microappProxy
  if (typeof proxy === 'function') return readPayload(await proxy(`/api/v1/traffic/${path}${suffix}`))
  return panelGet(`traffic/${path}`, params)
}
