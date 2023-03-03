import { UserConfig, ConfigEnv, loadEnv } from 'vite';
import vue from '@vitejs/plugin-vue';
import { createSvgIconsPlugin } from 'vite-plugin-svg-icons';
import path from 'path';
import WindiCSS from 'vite-plugin-windicss'
import { createStyleImportPlugin, ElementPlusResolve } from 'vite-plugin-style-import'

export default ({ mode }: ConfigEnv): UserConfig => {
  // 获取 .env 环境配置文件
  const env = loadEnv(mode, process.cwd());

  return {
    plugins: [
      vue(),
      WindiCSS(),
      createSvgIconsPlugin({
        // 指定需要缓存的图标文件夹
        iconDirs: [path.resolve(process.cwd(), 'src/assets/icons')],
        // 指定symbolId格式
        symbolId: 'icon-[dir]-[name]'
      }),
      createStyleImportPlugin({
        resolves: [ElementPlusResolve()],
        libs: [{
          libraryName: 'element-plus',
          esModule: true,
          resolveStyle: (name) => {
            return `element-plus/es/components/${name.substring(3)}/style/css`
          }
        }]
      }),
    ],
    css: {
      preprocessorOptions: {
        //define global scss variable
        scss: {
          javascriptEnabled: true,
          additionalData: `
            @use "@/styles/variables.scss" as *;
          `
        }
      }
    },
    // 本地反向代理解决浏览器跨域限制
    server: {
      host: '0.0.0.0',
      port: 3000,
      open: true, // 运行自动打开浏览器
      proxy: {
        [env.VITE_API_URL]: {
          // 本地API地址
          target: 'http://localhost:8000',
          changeOrigin: true,
        }
      }
    },
    resolve: {
      // Vite路径别名配置
      alias: {
        '@': path.resolve('./src')
      }
    }
  };
};
