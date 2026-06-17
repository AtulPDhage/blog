import { boot } from 'quasar/wrappers';
import DOMPurify from 'dompurify';
import type { App } from 'vue';

export const sanitizeHtml = (html: string): string => {
  return DOMPurify.sanitize(html);
};

export default boot(({ app }: { app: App }) => {
  app.config.globalProperties.$sanitizeHtml = sanitizeHtml;
});

declare module 'vue' {
  interface ComponentCustomProperties {
    $sanitizeHtml: (html: string) => string;
  }
}
