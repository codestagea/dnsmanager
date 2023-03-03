import request from "@/utils/http/requests";

export interface LoginReq {
  username: string;
  password: string;
}

export interface LoginResp {
  token: string,
}

export interface UserInfoResp {
  username: string;
  email: string;
  mobile: string;
  name: string;
  roles: string[];
  resources: string[];
}

/// 登录
export function apiLogin(data: LoginReq) {
  return request.post<LoginResp>('/v1/users/auth/login', { data: data })
}

export function apiCurrentUser() {
  return request.get<UserInfoResp>('/v1/users/me/info', {})
}


/// 注销
export function apiLogout() {
  return request.post('/v1/users/auth/logout', {})
}

// 用户分页
export const apiUserPage = (params, pagination) => {
  return request.request({
    url: `/v1/users/page`,
    method: 'get',
    params: { ...params, ...pagination }
  })
}

// 用户禁用/启用
export function apiUpdateUserStatus(username: string, status: string) {
  return request.post<UserInfoResp>(`/v1/users/${username}/status`, { data: { status } })
}

// 新增用户
export const apiUserSave = (data) => {
  return request.post(`/v1/users`, { data })
}
// 编辑用户
export const apiUserUpdate = (data) => {
  return request.post(`/v1/users/${data.username}`, { data })
}


// 用户详情
export const apiUserInfo = (username) => {
  return request.request({
    url: `/v1/users/${username}`,
    method: 'get',
  })
}

export const apiResetPassword = (data) => {
  return request.post('/v1/users/password/reset', { data })
}

export const apiAdminResetPassword = (username, password) => {
  return request.post(`/v1/users/${username}/password/reset`, { data: { password: password } })
}

export function apiDeleteUser(username: string) {
  return request.delete(`/v1/users/${username}`, {})
}

export function apiExportTemplate() {
  return Promise.resolve()
}

export function apiExportUser() {
  return Promise.resolve()
}

export function apiImportUser() {
  return Promise.resolve()
}
