const isWujie = () => Boolean(window.__POWERED_BY_WUJIE__)
let localToken

function localHeaders() {
  if (!import.meta.env.DEV) return {}
  localToken ||= window.prompt('请输入 W7Panel Token（仅用于本地测试）')?.trim()
  if (!localToken) throw new Error('本地测试需要 W7Panel Token')
  return { 'X-W7Panel-Token': localToken }
}

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
  return readPayload(await fetch(`/panel-api/v1/${path.replace(/^\//, '')}${suffix}`, { credentials: 'same-origin', headers: localHeaders() }))
}

export async function trafficGet(path, params = {}) {
  const query = new URLSearchParams()
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') query.set(key, String(value))
  })
  const suffix = query.size ? `?${query}` : ''
  const proxy = microProps().microappProxy
  if (typeof proxy === 'function') {
    try { return readPayload(await proxy(`/api/v1/traffic/${path}${suffix}`)) }
    catch (error) {
      if (error?.code !== 'INVALID_MICROAPP_PROXY_BASE_URL' || error.message !== 'proxy base URL is required') throw error
    }
  }
  return panelGet(`traffic/${path}`, params)
}
