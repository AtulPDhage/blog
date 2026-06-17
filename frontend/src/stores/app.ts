import { defineStore, acceptHMRUpdate } from 'pinia';
import axios from 'axios';
import Cookies from 'js-cookie';
import { Notify } from 'quasar';
import { user_service, blog_service } from '@/boot/axios';

export interface User {
  _id: string;
  name: string;
  email: string;
  image: string;
  instagram: string;
  linkedin: string;
  facebook: string;
  bio: string;
}

export interface Blog {
  id: string;
  title: string;
  description: string;
  blogcontent: string;
  image: string;
  category: string;
  author: string;
  created_at: string;
}

export interface SavedBlogType {
  id: string;
  userid: string;
  blogid: string;
  created_at: string;
}

export const blogCategories = [
  'Technology',
  'Health',
  'Entertainment',
  'Finance',
  'Travel',
  'Education',
  'Study',
];

export const useAppStore = defineStore('app', {
  state: () => ({
    user: null as User | null,
    isAuth: false,
    loading: true,
    blogs: null as Blog[] | null,
    blogLoading: false,
    searchQuery: '',
    category: '',
    savedBlogs: null as SavedBlogType[] | null,
    leftDrawerOpen: false,
  }),

  actions: {
    async fetchUser() {
      this.loading = true;
      try {
        const token = Cookies.get('token');
        if (!token) {
          this.isAuth = false;
          this.user = null;
          return;
        }

        const { data } = await axios.get(`${user_service}/api/v1/me`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        this.user = data.user;
        this.isAuth = true;
      } catch (err) {
        console.error('Fetch user error:', err);
        this.isAuth = false;
        this.user = null;
      } finally {
        this.loading = false;
      }
    },

    async fetchBlogs() {
      this.blogLoading = true;
      try {
        const { data } = await axios.get(
          `${blog_service}/api/v1/blog/all?searchQuery=${this.searchQuery}&category=${this.category}`,
        );
        this.blogs = data?.blogs || [];
      } catch (err) {
        console.error('Fetch blogs error:', err);
      } finally {
        this.blogLoading = false;
      }
    },

    async getSavedBlogs() {
      const token = Cookies.get('token');
      if (!token) return;
      try {
        const { data } = await axios.get(`${blog_service}/api/v1/blogs/saved/all`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        this.savedBlogs = data;
      } catch (error) {
        console.error('Get saved blogs error:', error);
      }
    },

    logoutUser() {
      Cookies.remove('token');
      this.user = null;
      this.isAuth = false;
      Notify.create({
        type: 'positive',
        message: 'User logged out',
        position: 'top',
      });
    },

    setSearchQuery(query: string) {
      this.searchQuery = query;
      void this.fetchBlogs();
    },

    setCategory(category: string) {
      this.category = category;
      void this.fetchBlogs();
    },

    async initApp() {
      const token = Cookies.get('token');
      if (token) {
        await this.fetchUser();
      } else {
        this.loading = false;
      }
      await this.fetchBlogs();
    },
  },
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useAppStore, import.meta.hot));
}
