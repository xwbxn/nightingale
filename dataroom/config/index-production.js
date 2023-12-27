window.ENV = 'production'
var productionConfig = {
  baseUrl: '/dataroom/bigScreenServer',
  fileUrlPrefix: '/dataroom/bigScreenServer/static'

}
// 必须的
window.CONFIG = configDeepMerge(window.CONFIG, productionConfig)
