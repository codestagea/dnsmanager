import { AxiosTransform, VAxios } from './vaxios'
import { ElMessage } from 'element-plus';
import { useUserStore } from '@/store/user'
import { useRouter } from 'vue-router'
import { reject } from 'lodash-es';

const ContentType = {
  urlencoded: 'application/x-www-form-urlencoded;charset=UTF-8',
  json: 'application/json',
  formData: 'multipart/form-data'
}
/**
 * @description: 数据处理，方便区分多种处理方式
 */
const transform: AxiosTransform = {
  /**
   * @description: 处理请求数据。如果数据不是预期格式，可直接抛出错误
   */
  transformRequestHook: (res: AxiosResponse<Result>) => {
    // const { isTransformResponse } = options;
    // 是否返回原生响应头 比如：需要获取响应头时使用该属性
    // if (isReturnNativeResponse) {
    //   return res;
    // }
    // 不进行任何处理，直接返回
    // // 用于页面代码可能需要直接获取code，data，message这些信息时开启
    // if (!isTransformResponse) {
    //   return res.data;
    // }
    // 错误的时候返回
    const { data } = res;
    const { code } = data;
    if (code === 0) {
      return data.data
    }

    // 未登录，重定向sso登录

    if (code === 3000) {
      ElMessage({
        message: data.message || '用户未登录',
        type: 'error'
      });
      const userStore = useUserStore()
      userStore.resetToken()

      const router = useRouter()
      router.push({ name: 'login' })
      window.location.href = data.message
    } else {
      ElMessage({
        message: data.message || '接口错误',
        type: 'error'
      });
      throw new Error(data.message || '接口错误')
    }
  },

  requestCatchHook: (error) => {
    console.log(error)
    const data = error.response?.data || {}// 请求data
    const status = error.response?.status // 请求状态吗

    if (status == 401) {
      ElMessage({
        message: '认证失败',
        type: 'error'
      });
      //token过期，或者非法
      // 重新登录
      // Modal.confirm({
      //   title: '确认注销',
      //   content: '您已经注销，您可以取消以留在此页面，也可以重新登录',
      //   okText: '重新登录',
      //   cancelText: '取消',
      //   onOk() {
      //     const userStore = useUserStore()
      //     userStore.resetToken()
      //     location.reload()
      //   },
      //   onCancel() { }
      // })
    } else {
      if (data.message) {
        ElMessage({
          message: data.message,
          type: 'error'
        });
      } else if (status === 403) {
        ElMessage({
          message: '无权限访问',
          type: 'error'
        });
      } else if (status == 404) {
        ElMessage({
          message: '接口不存在',
          type: 'error'
        });
      }
      else {
        let { message: msg } = error
        if (msg === 'Network Error') {
          msg = '连接异常'
        }
        if (msg.includes('timeout')) {
          msg = '请求超时'
        }
        if (msg.includes('Request failed with status code')) {
          const code = msg.substr(msg.length - 3)
          msg = '接口' + code + '异常'
        }
        ElMessage({
          message: msg,
          type: 'error'
        });
      }
    }
    return error
  },

  /**
   * @description: 请求拦截器处理
   */
  requestInterceptors: (config) => {
    // 请求之前处理config
    const userStore = useUserStore()
    // 请求前处理
    // if (token) {
    if (userStore.token) {
      // 设置headers key
      config.headers['Authorization'] = userStore.token
    }
    config.headers['Content-Type'] = ContentType[config.data instanceof FormData ? 'formData' : 'json']
    return config;
  },

  /**
   * @description: 响应拦截器处理
   */
  responseInterceptors: (res: AxiosResponse<any>) => {
    return res;
  },
};

const request = new VAxios(
  {
    baseURL: import.meta.env.VITE_API_URL,
    timeout: 10000,
    responseType: 'json',
    headers: {
      'X-Requested-With': 'XMLHttpRequest'
    }
  },
  transform
);
export default request;
