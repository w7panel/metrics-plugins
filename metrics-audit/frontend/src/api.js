const props = () => window.__POWERED_BY_WUJIE__ ? (window.$wujie?.props || {}) : (window.microData?.cd || {})

export async function auditGet(path, options = {}) {
  const query = new URLSearchParams()
  Object.entries(options.params || {}).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') query.set(key, String(value))
  })
  const suffix = query.size ? `?${query}` : ''
  const micro = props().microappProxy
  const panel = props().panelProxy
  let response
  if (typeof micro === 'function') response = await micro(`/api/v1/audit/${path}${suffix}`)
  else if (typeof panel === 'function') response = await panel(`audit/${path}${suffix}`)
  else response = await fetch(`/panel-api/v1/audit/${path}${suffix}`, { credentials: 'same-origin' })
  const payload = await response.json().catch(() => null)
  if (!response.ok) throw new Error(payload?.msg || payload?.message || `请求失败 (${response.status})`)
  return { data: payload?.data ?? payload }
}
