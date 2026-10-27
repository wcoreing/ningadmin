export const ADMIN_BASE = '/admin'
export const CMS_BASE = '/cms'

export function adminPath(path = '') {
  if (!path) return ADMIN_BASE
  if (path === ADMIN_BASE || path.startsWith(`${ADMIN_BASE}/`)) return path
  return ADMIN_BASE + (path.startsWith('/') ? path : `/${path}`)
}
