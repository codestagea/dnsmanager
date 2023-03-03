import { defineStore } from 'pinia';
import { apiLogin, apiCurrentUser, apiLogout, type LoginReq, type LoginResp, type UserInfoResp } from '@/api/user'
import { passwdEnc } from '@/utils/encrypt'
import { store } from '@/store';
import { ref } from 'vue';
import storage from '@/utils/storage'
import type { RouteRecordRaw } from 'vue-router';
import { transformRouteToMenu } from '@/utils/helper/menuHelper';
import { filter } from '@/utils/helper/treeHelper';
import { cloneDeep } from 'lodash-es';
import router, { resetRouter, constantRoutes } from '@/router';


interface UserState {
  userInfo: Nullable<UserInfo>;
  token?: string;
  roles: string[];
  perms: string[];
  menus: RouteRecordRaw[];
}

export const useUserStore = defineStore({
  id: 'app-user',
  state: (): UserState => ({
    userInfo: null,
    token: storage.get('user-token'),
    roles: [],
    perms: [],
    menus: [],
  }),
  getters: {
  },
  actions: {
    setToken(info: string | undefined) {
      console.log(storage.get('user-token'))
      this.token = info ? info : '';
      storage.set('user-token', info);
    },
    resetToken() {
      storage.remove('user-token')
      this.token = undefined
    },
    async loadUserInfo(): Promise<UserInfoResp | null> {
      try {
        const data = await apiCurrentUser<UserInfoResp>();
        this.userInfo = data
        return data
      } catch (error) {
        return Promise.reject(error);
      }
    },
    async login(loginForm: LoginReq): Promise<LoginResp | null> {
      try {
        const req = cloneDeep(loginForm)
        req.username = loginForm.username.trim()
        req.password = passwdEnc.encrypt(loginForm.password.trim())
        const loginResp = await apiLogin(req)
        const { token } = loginResp
        this.setToken(token)
        this.afterLoginAction()
        return loginResp
      } catch (e) {
        return Promise.reject(e)
      }
    },
    async logout() {
      try {
        await apiLogout();
      } catch {
        console.log('注销Token失败');
      }
      this.resetToken()
    },


    async afterLoginAction(goHome?: boolean): Promise<UserInfo | null> {
      if (!this.token) return null;
      // get user info
      const userResp = await this.loadUserInfo();
      const { roles = [], resources = [] } = userResp
      const routes = await this.buildRoutes(roles, resources)

      resetRouter()
      routes.forEach((route) => {
        router.addRoute(route as unknown as RouteRecordRaw);
      });
      return this.userInfo;
    },

    async buildRoutes(roles: string[], resources: string[]): Promise<RouteRecordRaw[]> {
      const asyncRoutesClone = cloneDeep(constantRoutes)
      this.roles = roles
      this.perms = resources


      if (roles.some((role) => role === 'admin')) {
        this.menus = transformRouteToMenu(asyncRoutesClone, true);
        return asyncRoutesClone;
      } else {
        const routeFilter = (route) => {
          const { meta } = route;
          const { permission } = meta || {};
          if (!permission) return true;
          return resources.some((res) => permission.includes(res));
        }
        const routes = filter<Router>(asyncRoutesClone, routeFilter);
        this.menus = transformRouteToMenu(routes, true);
        return routes;
      }
    },
  }
})

// 非setup
export function useUserStoreHook() {
  return useUserStore(store);
}
