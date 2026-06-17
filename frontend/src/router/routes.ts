import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    children: [
      {
        path: '',
        redirect: '/blogs',
      },
      {
        path: 'blogs',
        component: () => import('@/pages/IndexPage.vue'),
      },
      {
        path: 'login',
        component: () => import('@/pages/LoginPage.vue'),
        meta: { guestOnly: true },
      },
      {
        path: 'profile',
        component: () => import('@/pages/ProfilePage.vue'),
        meta: { requiresAuth: true },
      },
      {
        path: 'profile/:id',
        component: () => import('@/pages/PublicProfilePage.vue'),
      },
      {
        path: 'blog/new',
        component: () => import('@/pages/NewBlogPage.vue'),
        meta: { requiresAuth: true },
      },
      {
        path: 'blog/:id',
        component: () => import('@/pages/BlogDetailPage.vue'),
      },
      {
        path: 'blog/edit/:id',
        component: () => import('@/pages/EditBlogPage.vue'),
        meta: { requiresAuth: true },
      },
      {
        path: 'blog/saved',
        component: () => import('@/pages/SavedBlogsPage.vue'),
        meta: { requiresAuth: true },
      },
    ],
  },

  // Always leave this as last one
  {
    path: '/:catchAll(.*)*',
    component: () => import('@/pages/ErrorNotFound.vue'),
  },
];

export default routes;
