import router from './router'
import store from './store'
import { Message } from 'element-ui'
import NProgress from 'nprogress' // progress bar
import 'nprogress/nprogress.css' // progress bar style
import { getToken } from '@/utils/auth' // get token from cookie
import getPageTitle from '@/utils/get-page-title'

NProgress.configure({ showSpinner: false }) // NProgress Configuration

const whiteList = ['/login', '/register', '/auth-redirect'] // no redirect whitelist

router.beforeEach(async (to, from, next) => {
  // start progress bar
  NProgress.start()

  // set page title
  document.title = getPageTitle(to.meta.title)

  // determine whether the user has logged in
  const hasToken = getToken()

  if (hasToken) {
    if (to.path === '/login') {
      // if is logged in, redirect to the home page
      next({ path: '/' })
      NProgress.done() // hack: https://github.com/PanJiaChen/vue-element-admin/pull/2939
    } else {
      // determine whether the user has obtained his permission roles through getInfo
      // 已经登录了，去其它页面
      // store.getters.roles得到vuex中的角色
      // 如果登录了，我们会调用一个接口，去拿用户信息，在用户信息中，有当前用户的角色
      // 点击登录，先发一个登录请求，服务器响应一个token，前端把token存储到cookie
      // 紧接着发第二个请求，是用来获取用户信息的，前端把用户信息存储到了vuex中，用户信息中有一个角色
      // 也就是说，在vuex中是可以获取角色的  通过store.getters.roles
      // store.getters.roles.length > 0 表示vuex是有角色
      const hasRoles = store.getters.roles && store.getters.roles.length > 0
      if (hasRoles) {
        // 从vuex中获取角色，直接放行
        next()
      } else {
        // else表示vuex中没有角色  当你又去刷新浏览器时，vuex中就没有角色，vuex中的数据也是存储在内存
        try {
          // get user info
          // note: roles must be a object array! such as: ['admin'] or ,['developer','editor']
          // store.dispatch('user/getInfo')  重新获取用户信息   肯定是在vuex中发送ajax请求
          // roles 表示用户信息，用户信息中包含用户角色
          const { roles } = await store.dispatch('user/getInfo')

          // generate accessible routes map based on roles
          // dispatch('permission/generateRoutes', roles)  根据用户角色，生成路由规则
          const accessRoutes = await store.dispatch('permission/generateRoutes', roles)

          // dynamically add accessible routes
          // 一个路由器中，可以有很多的规则，计算了当前用户角色的个规则
          // 利用addRoutes，把规则，添加到路由器
          router.addRoutes(accessRoutes)

          // hack method to ensure that addRoutes is complete
          // set the replace: true, so the navigation will not leave a history record
          next({ ...to, replace: true })
        } catch (error) {
          // remove token and go to login page to re-login
          // 如果在生成规则时，出问题了
          // store.dispatch('user/resetToken')  清除token
          await store.dispatch('user/resetToken')
          Message.error(error || 'Has Error')
          // 重新回到登录页面
          next(`/login?redirect=${to.path}`)
          NProgress.done()
        }
      }
    }
  } else {
    /* has no token*/
    // 没有token
    // 没有token看一下，你访问的路径有没有在白名单中
    if (whiteList.indexOf(to.path) !== -1) {
      // in the free login whitelist, go directly
      // 如果在白名单中，就放行
      next()
    } else {
      // other pages that do not have permission to access are redirected to the login page.
      // 如果没有在白名单中，表示你访问的路由规则，需要登录
      // 需要登录，放行到登录页面
      next(`/login?redirect=${to.path}`)
      NProgress.done()
    }
  }
})

router.afterEach(() => {
  // finish progress bar
  NProgress.done()
})
