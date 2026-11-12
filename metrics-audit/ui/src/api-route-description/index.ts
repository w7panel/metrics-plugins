type Descriptions = Record<string, string>;

type DescribeOptions = {
  method?: string;
  route?: string;
  path?: string;
  fallback?: string;
  locale?: string;
};

const cache: Record<string, Descriptions> = {};
const pending: Record<string, Promise<Descriptions>> = {};

function currentLocale() {
  return localStorage.getItem('arco-locale') || 'zh-CN';
}

function normalizeLocale(locale?: string) {
  return locale === 'en-US' ? 'en-US' : 'zh-CN';
}

export async function loadApiRouteDescriptions(locale?: string) {
  const key = normalizeLocale(locale || currentLocale());
  if (cache[key]) return cache[key];
  if (pending[key]) return pending[key];
  pending[key] = (key === 'en-US'
    ? import('./locale/en-US')
    : import('./locale/zh-CN')
  ).then((module) => {
    cache[key] = expandDescriptions(module.default || {});
    return cache[key];
  }).finally(() => {
    delete pending[key];
  });
  return pending[key];
}

function expandDescriptions(descriptions: Descriptions) {
  const expanded: Descriptions = { ...descriptions };
  Object.entries(descriptions).forEach(([key, value]) => {
    const index = key.indexOf(' ');
    if (index <= 0) return;
    const method = key.slice(0, index);
    const route = key.slice(index + 1);
    const normalized = normalizeApiRoutePath(route);
    if (normalized && !expanded[`${method} ${normalized}`]) {
      expanded[`${method} ${normalized}`] = value;
    }
  });
  return expanded;
}

export function getLoadedApiRouteDescriptions(locale?: string) {
  return cache[normalizeLocale(locale || currentLocale())] || {};
}

export function normalizeApiRoutePath(route?: string) {
  return String(route || '').split('/').map((part) => {
    if (part.startsWith(':') || part.startsWith('*')) return '*';
    return part;
  }).join('/');
}

function routeCandidates(route?: string, path?: string) {
  const values = [route, path].filter(Boolean) as string[];
  const result: string[] = [];
  values.forEach((value) => {
    [value, normalizeApiRoutePath(value)].forEach((item) => {
      if (item && !result.includes(item)) result.push(item);
    });
  });
  return result;
}

export function fallbackApiRouteDescription(options: DescribeOptions = {}) {
  return options.fallback || '-';
}

export function resolveApiRouteDescription(
  descriptions: Descriptions,
  options: DescribeOptions = {},
) {
  const method = String(options.method || '').toUpperCase();
  const candidates = routeCandidates(options.route, options.path);
  for (const route of candidates) {
    const keys = [`${method} ${route}`, `* ${route}`];
    for (const key of keys) {
      if (descriptions[key]) return descriptions[key];
    }
  }
  if (!method && candidates.length === 0 && options.fallback) {
    return options.fallback;
  }
  return fallbackApiRouteDescription(options);
}

export async function describeApiRoute(options: DescribeOptions = {}) {
  const descriptions = await loadApiRouteDescriptions(options.locale);
  return resolveApiRouteDescription(descriptions, options);
}
