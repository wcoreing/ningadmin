let app = null
let ready = false

export function setApp(instance) {
  app = instance
}

export async function ensureModuleUi() {
  if (ready || !app) return
  const ElementPlus = (await import('element-plus')).default
  const icons = await import('@element-plus/icons-vue')
  await import('element-plus/dist/index.css')
  Object.entries(icons).forEach(([key, component]) => {
    if (key !== 'default' && component) app.component(key, component)
  })
  app.use(ElementPlus)
  ready = true
}
