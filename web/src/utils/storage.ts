interface storageOptType {
  namespace?: string
  type: string
  default_cache_time?: number
}

const options = {
  namespace: 'dnsmanager-', // key prefix
  type: 'localStorage', // storage name session, local, memory
  default_cache_time: 60 * 60 * 24 * 7,
}

let hasSetStorage = false

export const storage = {
  getKey: (key: string) => {
    return options.namespace + key
  },
  setOptions: (opt: storageOptType) => {
    if (hasSetStorage) {
      console.error('Has set storage:', options)
      return
    }
    Object.assign(options, opt)
    hasSetStorage = true
  },
  set: (key: string, value, expire: number | null = options.default_cache_time) => {
    const stringData = JSON.stringify({
      value,
      expire: expire !== null ? new Date().getTime() + expire * 1000 : null
    })
    const storageKey = storage.getKey(key)
    console.log("storage set", storageKey, stringData)
    // localStorage.setItem(storageKey, stringData)
    window[options.type].setItem(storage.getKey(key), stringData)
  },
  setObj: (key: string, newVal, expire: number | null = options.default_cache_time) => {
    const oldVal = storage.get(key)
    if (!oldVal) {
      storage.set(key, newVal, expire)
    } else {
      Object.assign(oldVal, newVal)
      storage.set(key, oldVal, expire)
    }
  },
  /**
   * 读取缓存
   * @param {string} key 缓存键
   * @param {*=} def 默认值
   */
  get: (key: string) => {
    const storageKey = storage.getKey(key)
    let item = window[options.type].getItem(storageKey)
    // let item = localStorage.getItem(storageKey)
    // console.log("storage get", storageKey, item)
    if (item) {
      try {
        const data = JSON.parse(item)
        const { value, expire } = data
        // 在有效期内直接返回
        if (expire === null || expire >= Date.now()) {
          return value
        }
        storage.remove(storage.getKey(key))
      } catch (e) {
        console.error(e)
      }
    }
    return null
  },
  remove: (key: string) => {
    const storageKey = storage.getKey(key)
    console.log("storage remove", storageKey)
    window[options.type].removeItem(storageKey)
    // localStorage.removeItem(storageKey)
  },
  clear: (): void => {
    window[options.type].storage.clear()
    // localStorage.clear()
  },
}

export function getSidebarStatus() {
  return storage.get('sidebarStatus');
}

export function setSidebarStatus(sidebarStatus: string) {
  storage.set('sidebarStatus', sidebarStatus);
}
// 布局大小
export function getSize() {
  return storage.get('size');
}

export function setSize(size: string) {
  storage.set('size', size);
}

// 语言

export function getLanguage() {
  return storage.get('language');
}

export function setLanguage(language: string) {
  storage.set('language', language);
}


export default storage
