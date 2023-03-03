import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';

export const Layout = () => import('@/layout/index.vue');

// 静态路由
export const constantRoutes: RouteRecordRaw[] = [
  {
    path: '/redirect',
    component: Layout,
    meta: { hidden: true },
    children: [
      {
        path: '/redirect/:path(.*)',
        component: () => import('@/views/redirect/index.vue')
      }
    ]
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/login/index.vue'),
    meta: { hidden: true }
  },
  {
    path: '/404',
    component: () => import('@/views/error-page/404.vue'),
    meta: { hidden: true }
  },
  {
    path: '/',
    component: Layout,
    redirect: 'zone',
    meta: { title: 'dns' },
    children: [
      {
        path: 'zone',
        component: () => import('@/views/zone/index.vue'),
        name: 'dns_zone',
        meta: { title: 'dns_zone', affix: true }
      },
      {
        path: 'view',
        component: () => import('@/views/dns-view/index.vue'),
        name: 'dns_view',
        meta: { title: 'dns_view', affix: true }
      },
      {
        path: 'operation/log',
        name: 'operation_log',
        component: () => import('@/views/operation-log/index.vue'),
        meta: { title: 'operation_log', affix: true },
      },
      {
        path: 'zone/:zoneId/records',
        name: 'dns_record',
        component: () => import('@/views/records/index.vue'),
        meta: { title: 'dns_record', icon: 'bx-analyse', hidden: true },
      },
      {
        path: 'operation/log/:targetType/:targetId',
        name: 'target_operation_log',
        component: () => import('@/views/operation-log/target-log.vue'),
        meta: { title: '操作记录', icon: 'bx-analyse', hidden: true },
      }
    ]
  },
  {
    path: '/sys',
    component: Layout,
    redirect: '/user',
    meta: { title: 'system' },
    children: [
      {
        path: 'user',
        component: () => import('@/views/system/user/index.vue'),
        name: 'sys-user',
        meta: { title: 'sys_user', icon: 'homepage', affix: true }
      },
    ]
  }

];

// 创建路由
const router = createRouter({
  history: createWebHistory(import.meta.env.VITE_PUBLIC_PATH),
  routes: constantRoutes as RouteRecordRaw[],
  // 刷新时，滚动条位置还原
  scrollBehavior: () => ({ left: 0, top: 0 })
});

// 重置路由
export function resetRouter() {
  const newRouter = createRouter({
    history: createWebHistory(import.meta.env.VITE_PUBLIC_PATH),
    routes: constantRoutes as unknown as RouteRecordRaw[],
  });
  (router as any).matcher = (newRouter as any).matcher // reset router
}

export default router;
