import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

/* Layout */
import Layout from '@/layout'

/* Router Modules */
// import componentsRouter from './modules/components'
// import chartsRouter from './modules/charts'
// import tableRouter from './modules/table'
// import nestedRouter from './modules/nested'

/**
 * Note: sub-menu only appear when route children.length >= 1
 * Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin','editor']    control the page roles (you can set multiple roles)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)
    icon: 'svg-name'/'el-icon-x' the icon show in the sidebar
    noCache: true                if set true, the page will no be cached(default is false)
    affix: true                  if set true, the tag will affix in the tags-view
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set
  }
 */

/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */

// constantRoutes是静态路由
// 本项目，路由分成了两大模块：静态路由 和 动态路由
// 静态路由：所有的用户可以访问，不需要权限
// 动态路由：需要权限，如果有权限，就可以访问，如果没有权限，就不能访问
// 路则规则：就是一个对象
//     path: '/redirect',  访问的url
//     component: Layout,  访问出口中放什么组件  在一级中币出口中放 Layout 组件
//     hidden: true,  隐藏   把侧边栏中不能看到声明式导航
//     children: [  配置二级路由
export const constantRoutes = [
  {
    path: '/redirect',
    component: Layout,
    hidden: true,
    children: [
      {
        path: '/redirect/:path(.*)',
        component: () => import('@/views/redirect/index')
      }
    ]
  },
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },
  {
    path: '/register',
    component: () => import('@/views/login/register'),
    hidden: true
  },
  {
    path: '/auth-redirect',
    component: () => import('@/views/login/auth-redirect'),
    hidden: true
  },
  {
    path: '/404',
    component: () => import('@/views/error-page/404'),
    hidden: true
  },
  {
    path: '/401',
    component: () => import('@/views/error-page/401'),
    hidden: true
  },
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        component: () => import('@/views/dashboard/index'),
        name: 'Dashboard',
        meta: { title: 'Dashboard', icon: 'dashboard', affix: true }
      }
    ]
  },
  // cluster-info
  // noCache 不保留在标签栏
  // affix: true 固定在标签栏
  {
    path: '/clusterinfo',
    component: Layout,
    // redirect: '/clusterinfo/index',
    children: [
      {
        path: 'index',
        component: () => import('@/views/cluster-resources/cluster-list.vue'),
        name: 'ClusterInfo',
        meta: { title: '集群信息', icon: 'dashboard', noCache: true }
      }
    ]
  },
  // cluster
  {
    path: '/cluster',
    component: Layout,
    redirect: '/cluster/node',
    alwaysShow: true, // will always show the root menu
    name: 'Cluster',
    meta: {
      title: '集群资源',
      icon: 'lock'
    },
    children: [
      {
        path: 'node',
        component: () => import('@/views/cluster-resources/node-list'),
        name: 'Nodes',
        meta: {
          title: '集群节点'
        }
      },
      {
        path: 'namespace',
        component: () => import('@/views/cluster-resources/namespace-list'),
        name: 'Namespaces',
        meta: {
          title: '命名空间'
        }
      }
    ]
  },
  // workload
  {
    path: '/workload',
    component: Layout,
    redirect: '/workload/deployment',
    alwaysShow: true, // will always show the root menu
    name: 'Workload',
    meta: {
      title: '工作负载',
      icon: 'lock'
    },
    children: [
      {
        path: 'deployment',
        component: () => import('@/views/workloads/deployment-list.vue'),
        name: 'Deployment',
        meta: {
          title: '无状态副本'
        },
        children: [
          {
            path: 'create',
            component: () => import('@/views/workloads/deployment-add.vue'),
            name: 'DeploymentCreate',
            meta: {
              title: '创建无状态副本'
            }
          }
        ]
      },
      {
        path: 'statefulset',
        component: () => import('@/views/workloads/statefulset-list.vue'),
        name: 'Statefulset',
        meta: {
          title: '有状态副本'
        }
      },
      {
        path: 'daemonset',
        component: () => import('@/views/workloads/daemonset-list.vue'),
        name: 'Daemonset',
        meta: {
          title: '守护进程'
        }
      },
      {
        path: 'pod',
        component: () => import('@/views/workloads/pod-list.vue'),
        name: 'Pod',
        meta: {
          title: '容器组'
        }
      }
    ]
  },
  // service
  {
    path: '/loadbalancing',
    component: Layout,
    redirect: '/loadbalancing/service',
    alwaysShow: true, // will always show the root menu
    name: 'Load-balancing',
    meta: {
      title: '负载均衡',
      icon: 'lock'
    },
    children: [
      {
        path: 'service',
        component: () => import('@/views/load-balancing/service-list.vue'),
        name: 'Service',
        meta: {
          title: '服务'
        }
      },
      {
        path: 'ingress',
        component: () => import('@/views/load-balancing/ingress-list.vue'),
        name: 'Ingress',
        meta: {
          title: '应用路由'
        }
      },
      {
        path: 'ingressclass',
        component: () => import('@/views/load-balancing/ingressclass-list.vue'),
        name: 'IngressClass',
        meta: {
          title: '应用路由类'
        }
      }
    ]
  },
  // profiles
  {
    path: '/profiles',
    component: Layout,
    redirect: '/profiles/configmap',
    alwaysShow: true, // will always show the root menu
    name: 'Profiles',
    meta: {
      title: '配置文件',
      icon: 'lock'
    },
    children: [
      {
        path: 'configmap',
        component: () => import('@/views/profiles/configmap-list.vue'),
        name: 'Configmap',
        meta: {
          title: '配置字典'
        }
      },
      {
        path: 'secret',
        component: () => import('@/views/profiles/secret-list.vue'),
        name: 'Secret',
        meta: {
          title: '保密字典'
        }
      }
    ]
  },
  // storage
  {
    path: '/storage',
    component: Layout,
    redirect: '/storage/persistenvolume',
    alwaysShow: true, // will always show the root menu
    name: 'storage',
    meta: {
      title: '存储资源',
      icon: 'lock'
    },
    children: [
      {
        path: 'persistenvolume',
        component: () => import('@/views/storage-resources/persistenvolume-list.vue'),
        name: 'PersistenVolume',
        meta: {
          title: '持久卷'
        }
      },
      {
        path: 'persistenvolumeclaim',
        component: () => import('@/views/storage-resources/persistenvolumeclaim-list.vue'),
        name: 'PersistenVolumeClaim',
        meta: {
          title: '持久卷声明'
        }
      },
      {
        path: 'storageclass',
        component: () => import('@/views/storage-resources/storageclass-list.vue'),
        name: 'StorageClass',
        meta: {
          title: '存储类'
        }
      }
    ]
  }
]

/**
 * asyncRoutes
 * the routes that need to be dynamically loaded based on user roles
 */
// asyncRoutes是动态路由
// 页面级（路由级）权限：
//     不同用户，登录到系统，看到的侧边栏是不一样，也就是有不同的页面
//     同一个页面，有的用户可以访问，有的用户不能访问
//  并不是说，你在下面配置完就OK，背后还有很多代码
export const asyncRoutes = [
  // user
  {
    path: '/user',
    component: Layout,
    redirect: '/user/list',
    alwaysShow: true, // will always show the root menu
    name: 'User',
    meta: {
      title: '用户管理',
      icon: 'lock',
      // roles表示什么样的用户可以访问permission
      // 不同的用户有不同的角色
      // 本项目就两个角色：admin   editor
      // roles: ['admin', 'editor']  表示amdin可以访问persmisson  editor也可以访问persmisson
      roles: ['admin'] // you can set roles in root nav
    },
    children: [
      {
        path: 'list',
        component: () => import('@/views/user-management/user-list.vue'),
        name: 'List',
        meta: {
          title: '用户列表',
          roles: ['admin'] // or you can only set roles in sub nav
        }
      },
      {
        path: 'permission',
        component: () => import('@/views/user-management/user-permission.vue'),
        name: 'Permission',
        meta: {
          title: '权限管理',
          roles: ['admin']
        }
      }
    ]
  },
  // system
  {
    path: '/system',
    component: Layout,
    redirect: '/system/email',
    alwaysShow: true, // will always show the root menu
    name: 'System',
    meta: {
      title: '系统设置',
      icon: 'lock',
      // roles表示什么样的用户可以访问permission
      // 不同的用户有不同的角色
      // 本项目就两个角色：admin   editor
      // roles: ['admin', 'editor']  表示amdin可以访问persmisson  editor也可以访问persmisson
      roles: ['admin'] // you can set roles in root nav
    },
    children: [
      {
        path: 'email',
        component: () => import('@/views/system-management/system-mail.vue'),
        name: 'Email',
        meta: {
          title: '邮箱设置',
          roles: ['admin']// or you can only set roles in sub nav
        }
      }
    ]
  },
  // 404 page must be placed at the end !!!
  { path: '*', redirect: '/404', hidden: true }
]

const createRouter = () => new Router({
  // mode: 'history', // require service support
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRoutes
})

const router = createRouter()

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // reset router
}

export default router
