const props = () => window.__POWERED_BY_WUJIE__ ? (window.$wujie?.props || {}) : (window.microData?.cd || {})
let localToken

function localHeaders() {
  if (!import.meta.env.DEV) return {}
  localToken ||= window.prompt('请输入 W7Panel Token（仅用于本地测试）')?.trim()
  if (!localToken) throw new Error('本地测试需要 W7Panel Token')
  return { 'X-W7Panel-Token': localToken }
}

export async function auditGet(path, options = {}) {
  const query = new URLSearchParams()
  Object.entries(options.params || {}).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') query.set(key, String(value))
  })
  const suffix = query.size ? `?${query}` : ''
  const { microappProxy: micro, panelProxy: panel } = props()
  let response
  if (typeof micro === 'function') {
    try { response = await micro(`/api/v1/audit/${path}${suffix}`) }
    catch (error) {
      if (error?.code !== 'INVALID_MICROAPP_PROXY_BASE_URL' || error.message !== 'proxy base URL is required') throw error
    }
  }
  if (!response && typeof panel === 'function') response = await panel(`audit/${path}${suffix}`)
  if (!response) response = await fetch(`/panel-api/v1/audit/${path}${suffix}`, { credentials: 'same-origin', headers: localHeaders() })
  const payload = await response.json().catch(() => null)
  if (!response.ok) throw new Error(payload?.msg || payload?.message || `请求失败 (${response.status})`)
  return { data: payload?.data ?? payload }
}
