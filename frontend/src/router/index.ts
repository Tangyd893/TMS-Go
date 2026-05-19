import { createRouter, createWebHistory } from 'vue-router'
import { setupRouterGuards } from './guards'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/login/LoginView.vue'),
      meta: { title: '登录' },
    },
    {
      path: '/',
      component: () => import('@/layouts/BasicLayout.vue'),
      redirect: '/dashboard',
      children: [
        {
          path: 'dashboard',
          name: 'dashboard',
          component: () => import('@/views/dashboard/DashboardView.vue'),
          meta: { title: '工作台' },
        },
        // 基础资料
        {
          path: 'base/customers',
          name: 'customers',
          component: () => import('@/views/base/CustomerList.vue'),
          meta: { title: '客户管理' },
        },
        {
          path: 'base/carriers',
          name: 'carriers',
          component: () => import('@/views/base/CarrierList.vue'),
          meta: { title: '承运商管理' },
        },
        {
          path: 'base/vehicles',
          name: 'vehicles',
          component: () => import('@/views/base/VehicleList.vue'),
          meta: { title: '车辆管理' },
        },
        {
          path: 'base/drivers',
          name: 'drivers',
          component: () => import('@/views/base/DriverList.vue'),
          meta: { title: '司机管理' },
        },
        {
          path: 'base/routes',
          name: 'routes',
          component: () => import('@/views/base/RouteList.vue'),
          meta: { title: '线路管理' },
        },
        {
          path: 'base/stations',
          name: 'stations',
          component: () => import('@/views/base/StationList.vue'),
          meta: { title: '站点管理' },
        },
        // 订单管理
        {
          path: 'order/list',
          name: 'orderList',
          component: () => import('@/views/order/OrderList.vue'),
          meta: { title: '订单列表' },
        },
        {
          path: 'order/create',
          name: 'orderCreate',
          component: () => import('@/views/order/OrderCreate.vue'),
          meta: { title: '创建订单' },
        },
        {
          path: 'order/:id',
          name: 'orderDetail',
          component: () => import('@/views/order/OrderDetail.vue'),
          meta: { title: '订单详情' },
        },
        // 调度管理
        {
          path: 'dispatch/list',
          name: 'dispatchList',
          component: () => import('@/views/dispatch/DispatchList.vue'),
          meta: { title: '调度管理' },
        },
        {
          path: 'dispatch/plans',
          name: 'dispatchPlans',
          component: () => import('@/views/dispatch/DispatchPlanList.vue'),
          meta: { title: '调度计划' },
        },
        {
          path: 'dispatch/create',
          name: 'dispatchCreate',
          component: () => import('@/views/dispatch/DispatchPlanCreate.vue'),
          meta: { title: '创建调度' },
        },
        // 运输执行
        {
          path: 'transport/tasks',
          name: 'transportTasks',
          component: () => import('@/views/transport/TaskList.vue'),
          meta: { title: '运输任务' },
        },
        {
          path: 'transport/task/:id',
          name: 'transportTaskDetail',
          component: () => import('@/views/transport/TaskDetail.vue'),
          meta: { title: '任务详情' },
        },
        // 异常管理
        {
          path: 'exception/list',
          name: 'exceptionList',
          component: () => import('@/views/exception/ExceptionList.vue'),
          meta: { title: '异常管理' },
        },
        // 费用结算
        {
          path: 'finance/receivables',
          name: 'financeReceivables',
          component: () => import('@/views/finance/ReceivableList.vue'),
          meta: { title: '应收管理' },
        },
        {
          path: 'finance/payables',
          name: 'financePayables',
          component: () => import('@/views/finance/PayableList.vue'),
          meta: { title: '应付管理' },
        },
        {
          path: 'finance/statements',
          name: 'financeStatements',
          component: () => import('@/views/finance/StatementList.vue'),
          meta: { title: '对账管理' },
        },
        {
          path: 'finance/settlements',
          name: 'financeSettlements',
          component: () => import('@/views/finance/SettlementList.vue'),
          meta: { title: '结算管理' },
        },
        // 报表中心
        {
          path: 'report/order-stats',
          name: 'reportOrderStats',
          component: () => import('@/views/report/OrderStatsView.vue'),
          meta: { title: '订单统计' },
        },
        {
          path: 'report/transport-efficiency',
          name: 'reportTransportEfficiency',
          component: () => import('@/views/report/TransportEfficiencyView.vue'),
          meta: { title: '运输效率' },
        },
        // 系统管理
        {
          path: 'system/dict',
          name: 'dict',
          component: () => import('@/views/system/DictList.vue'),
          meta: { title: '字典管理' },
        },
      ],
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
})

setupRouterGuards(router)

export default router
