import { boot } from 'quasar/wrappers';
import type { AxiosInstance } from 'axios';
import axios from 'axios';
import type { App } from 'vue';

declare module 'vue' {
  interface ComponentCustomProperties {
    $axios: AxiosInstance;
    $api: AxiosInstance;
  }
}

// Microservice base URLs
export const user_service = 'https://user-service-6ffj.onrender.com';
export const author_service = 'https://author-service-fjdq.onrender.com';
export const blog_service = 'https://blog-service-agq0.onrender.com';

const api = axios.create();

export default boot(({ app }: { app: App }) => {
  // for use inside Vue files (Options API) through this.$axios and this.$api
  app.config.globalProperties.$axios = axios;
  app.config.globalProperties.$api = api;
});

export { api };
