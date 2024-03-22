import { asyncRoutes, constantRoutes } from '@/router'

/**
 * Use meta.role to determine if the current user has permission
 * @param roles
 * @param route
 */
// 使用roles:Array<string> 和 route:{path,component,meta:{roles:[]}}
// 判断当前用户能不能访问当前路由规则，返回布尔值
function hasPermission(roles, route) {
  if (route.meta && route.meta.roles) {
    // 如果路由规则上有meta，并且meta上有roles这个自定义字段，说明这条路由是有权限的，要进行过滤权限。
    return roles.some(role => route.meta.roles.includes(role))
  } else {
    // 如果路由规则上没有meta这个属性，或者meta上没有roles这个自定义字段，任何用户都可以访问。
    return true
  }
}

/**
 * Filter asynchronous routing tables by recursion
 * @param routes asyncRoutes
 * @param roles
 */
// 生成当前用户可访问的路由规则。
export function filterAsyncRoutes(routes, roles) {
  const res = []

  routes.forEach(route => {
    // 深复制（深拷贝）仅对那些比较简单对象进行深拷贝
    const tmp = { ...route }
    if (hasPermission(roles, tmp)) {
      // 有没有嵌套视图，如果进行递归
      if (tmp.children) {
        tmp.children = filterAsyncRoutes(tmp.children, roles)
      }
      res.push(tmp)
    }
  })

  return res
}

const state = {
  routes: [],
  addRoutes: []
}

const mutations = {
  SET_ROUTES: (state, routes) => {
    state.addRoutes = routes
    state.routes = constantRoutes.concat(routes)
  }
}

const actions = {
  /*
   * 生成当前用户可访问的路由规则。
   * @param roles
   */
  generateRoutes({ commit }, roles) {
    return new Promise(resolve => {
      let accessedRoutes
      if (roles.includes('admin')) {
        accessedRoutes = asyncRoutes || []
      } else {
        accessedRoutes = filterAsyncRoutes(asyncRoutes, roles)
      }
      commit('SET_ROUTES', accessedRoutes)
      resolve(accessedRoutes)
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
