import type { AxiosRequestConfig, AxiosInstance, AxiosResponse } from 'axios';
import axios from 'axios';
import JSONBig from 'json-bigint'

export abstract class AxiosTransform {
  /**
   * @description: Request successfully processed
   */
  transformRequestHook?: (res: AxiosResponse<Result>) => any;

  /**
 * @description: 请求失败处理
 */
  requestCatchHook?: (error: Error) => Promise<any>;

  /**
* @description: 请求之前的拦截器
*/
  requestInterceptors?: (config: AxiosRequestConfig) => AxiosRequestConfig;

  /**
   * @description: 请求之后的拦截器
   */
  responseInterceptors?: (res: AxiosResponse<any>) => AxiosResponse<any>;
}

/**
 * @description:  axios module
 */
export class VAxios {
  private axiosInstance: AxiosInstance;
  private transform: AxiosTransform;

  constructor(config?: AxiosRequestConfig, transfrom?: AxiosTransform) {
    this.axiosInstance = axios.create(config);
    this.axiosInstance.defaults.transformResponse = [(data: any) => {
      return JSONBig.parse(data)
    }]
    this.transform = transfrom || {}
    this.setupInterceptors(transfrom);
  }

  getAxios(): AxiosInstance {
    return this.axiosInstance;
  }

  /**
   * @description: Interceptor configuration
   */
  private setupInterceptors(transform?: AxiosTransform) {
    const { requestInterceptors, responseInterceptors } = transform || {};
    // Request interceptor configuration processing
    if (requestInterceptors) {
      this.axiosInstance.interceptors.request.use((config: AxiosRequestConfig) => {
        config = requestInterceptors(config);
        return config;
      }, undefined);
    }

    // Response result interceptor processing
    if (responseInterceptors) {
      this.axiosInstance.interceptors.response.use((res: AxiosResponse<any>) => {
        res = responseInterceptors(res);
        return res
      }, undefined);
    }
  }

  /**
   * @description:  File Upload
   */
  uploadFile<T = any>(config: AxiosRequestConfig, params: UploadFileParams) {
    const formData = new window.FormData();
    const customFilename = params.name || 'file';

    if (params.filename) {
      formData.append(customFilename, params.file, params.filename);
    } else {
      formData.append(customFilename, params.file);
    }

    if (params.data) {
      Object.keys(params.data).forEach((key) => {
        const value = params.data![key];
        if (Array.isArray(value)) {
          value.forEach((item) => {
            formData.append(`${key}[]`, item);
          });
          return;
        }

        formData.append(key, params.data![key]);
      });
    }

    return this.axiosInstance.request<T>({
      ...config,
      method: 'POST',
      data: formData,
      headers: { 'Content-type': 'multipart/form-data;charset=UTF-8' },
    });
  }

  get<T = any>(url: string, config: AxiosRequestConfig): Promise<T> {
    return this.request({ ...config, method: 'GET', url });
  }

  post<T = any>(url: string, config: AxiosRequestConfig): Promise<T> {
    return this.request({ ...config, method: 'POST', url });
  }

  put<T = any>(url: string, config: AxiosRequestConfig): Promise<T> {
    return this.request({ ...config, method: 'PUT', url });
  }

  delete<T = any>(url: string, config: AxiosRequestConfig): Promise<T> {
    return this.request({ ...config, method: 'DELETE', url });
  }

  request<T = any>(config: AxiosRequestConfig): Promise<T> {
    const { requestCatchHook, transformRequestHook } = this.transform || {};
    return new Promise((resolve, reject) => {
      this.axiosInstance
        .request<any, AxiosResponse<R>>(config)
        .then((res: AxiosResponse<R>) => {
          if (transformRequestHook) {
            try {
              const ret = transformRequestHook(res);
              resolve(ret);
            } catch (err) {
              reject(err || new Error('request error!'));
            }
            return;
          }
          resolve(res as unknown as Promise<T>);
        })
        .catch((e: Error) => {
          if (requestCatchHook) {
            reject(requestCatchHook(e));
          } else {
            reject(e)
          }
        });
    });
  }
}
