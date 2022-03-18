function makePageEntry( part, id, name){

  return {
      path: "/" + part + "_" + id.toLowerCase(),
      component: ()=>import( './views/' + part+"/" + id),
      name: name,
  }

}

let routes = [
  {
    path: '/login',
    component: () => import('./views/Login'),
    name: '',
    hidden: true
  },
  {
    path: '/',
    component: () => import('./views/Home'),
    name: '进程管理',
    children: [
      makePageEntry('procmgr', 'ProcMgr', '进程控制'),
      makePageEntry('procmgr', 'TaskMgr', '任务管理'),
    ]
  },
  {
    path: '/',
    component: () => import('./views/Home'),
    name: '系统管理',
    children: [
      makePageEntry('sysmgr', 'ServerMgr', '服务器管理'),
    ]
  },
  {
    path: '/404',
    component: () => import('./views/404'),
    name: '',
    hidden: true
  },
  {
    path: '*',
    hidden: true,
    redirect: { path: '/404' }
  }
];

export default routes;
