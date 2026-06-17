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
export const user_service = import.meta.env.VITE_USER_SERVICE;
export const author_service = import.meta.env.VITE_AUTHOR_SERVICE;
export const blog_service = import.meta.env.VITE_BLOG_SERVICE;

const api = axios.create();

export default boot(({ app }: { app: App }) => {
  // for use inside Vue files (Options API) through this.$axios and this.$api
  app.config.globalProperties.$axios = axios;
  app.config.globalProperties.$api = api;
});

export { api };
