import router from '@/router';
import { RouteRecordRaw } from 'vue-router';
import { useUserStoreHook } from '@/store/user';
import storage from '@/utils/storage';
import NProgress from 'nprogress';
import 'nprogress/nprogress.css';
NProgress.configure({ showSpinner: false }); // 进度条


// 白名单路由
const whiteList = ['/login'];

router.beforeEach(async (to, from, next) => {
  NProgress.start();
  const userStore = useUserStoreHook();
  const token = storage.get('user-token')
  console.log('usertoken', token)
  if (token) {
    const user = userStore.userInfo
    console.log('user', user)

    if (user && user.username) {
      if (to.path === '/login') {
        next({ path: '/' });
        NProgress.done();
      } else {
        next()
      }
    } else {
      try {
        const fetchedUser = await userStore.afterLoginAction()
        if (fetchedUser && fetchedUser.username) {
          next({ ...to, replace: true })
        } else {
          console.log("afterLoginAction fetch user: ", fetchedUser)
          throw new Error("无法获取用户信息")
        }
      } catch (error) {
        // 移除token 返回登录
        await userStore.resetToken()
        next(`/login?redirect=${to.path}`);
        NProgress.done()
      }
    }

  } else {
    // 未登录可以访问白名单页面
    if (whiteList.indexOf(to.path) !== -1) {
      next();
    } else {
      next(`/login?redirect=${to.path}`);
      NProgress.done();
    }
  }
});

router.afterEach(() => {
  NProgress.done();
});
