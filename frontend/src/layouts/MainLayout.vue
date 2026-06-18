<template>
  <q-layout view="lHh Lpr lFf">
    <!-- Header with Glassmorphism -->
    <q-header
      borderless
      class="navbar-blur"
      :class="isDarkMode ? 'navbar-blur--dark text-white' : 'navbar-blur--light text-grey-9'"
    >
      <q-toolbar class="main-toolbar q-py-sm">
        <q-btn
          v-if="isBlogsRoute"
          flat
          dense
          round
          icon="menu"
          aria-label="Menu"
          class="filter-menu-btn q-mr-sm gt-xs"
          @click="toggleLeftDrawer"
        />

        <div
          class="app-title row items-center cursor-pointer q-gutter-x-xs"
          @click="router.push('/blogs')"
        >
          <q-avatar size="32px" class="bg-indigo-1 text-primary shadow-sm brand-logo-container">
            <q-icon name="auto_awesome" size="18px" class="brand-logo-icon" />
          </q-avatar>
          <span class="text-bold text-h6 font-brand logo-text">Postly</span>
        </div>

        <!-- Desktop Navigation Links -->
        <div class="nav-links gt-xs row items-center no-wrap q-gutter-sm">
          <q-btn
            flat
            no-caps
            label="Home"
            to="/blogs"
            class="nav-link-btn"
            :class="{ 'nav-link-btn--active': route.path === '/blogs' }"
          />

          <q-btn
            v-if="store.isAuth"
            flat
            no-caps
            label="Saved Blogs"
            to="/blog/saved"
            class="nav-link-btn"
            :class="{ 'nav-link-btn--active': route.path === '/blog/saved' }"
          />
          <q-btn
            v-if="store.isAuth"
            flat
            no-caps
            label="Create Blog"
            to="/blog/new"
            class="nav-link-btn"
            :class="{ 'nav-link-btn--active': route.path === '/blog/new' }"
          />
        </div>

        <q-space />

        <!-- Desktop Right Actions -->
        <div class="nav-actions gt-xs row items-center no-wrap q-gutter-md">
          <q-btn
            flat
            round
            dense
            :icon="isDarkMode ? 'light_mode' : 'dark_mode'"
            class="theme-toggle-btn"
            :color="isDarkMode ? 'amber-4' : 'grey-7'"
            @click="toggleDarkMode"
          >
            <q-tooltip class="bg-grey-9 text-white">{{ isDarkMode ? 'Light mode' : 'Dark mode' }}</q-tooltip>
          </q-btn>

          <!-- Profile Dropdown Trigger -->
          <q-btn
            v-if="store.isAuth && store.user"
            flat
            round
            dense
            class="profile-avatar-btn"
          >
            <q-avatar size="36px" class="avatar-ring shadow-sm">
              <img :src="store.user.image || '/default-avatar.png'" alt="profile" />
            </q-avatar>

            <q-menu
              anchor="bottom right"
              self="top right"
              class="profile-menu shadow-md"
              style="border-radius: var(--radius-md); border: 1px solid var(--border-color); background-color: var(--bg-card);"
            >
              <div class="dropdown-profile-header q-px-md q-py-md row items-center no-wrap q-gutter-x-sm">
                <q-avatar size="36px" class="shadow-sm">
                  <img :src="store.user.image || '/default-avatar.png'" alt="profile" />
                </q-avatar>
                <div class="column overflow-hidden">
                  <span class="text-weight-bold text-caption ellipsis text-main">{{ store.user.name }}</span>
                  <span class="text-caption text-sub ellipsis" style="font-size: 11px">{{ store.user.email }}</span>
                </div>
              </div>

              <q-separator />

              <q-list style="min-width: 200px" class="q-py-xs dropdown-list">
                <q-item clickable v-close-popup to="/profile" class="dropdown-item">
                  <q-item-section avatar>
                    <q-icon name="person" size="xs" />
                  </q-item-section>
                  <q-item-section>My Profile</q-item-section>
                </q-item>
                
                <q-item clickable v-close-popup to="/blog/saved" class="dropdown-item">
                  <q-item-section avatar>
                    <q-icon name="bookmarks" size="xs" />
                  </q-item-section>
                  <q-item-section>Saved Blogs</q-item-section>
                </q-item>

                <q-separator class="q-my-xs" />

                <q-item clickable v-close-popup @click="handleLogout" class="dropdown-item text-negative-item">
                  <q-item-section avatar>
                    <q-icon name="logout" size="xs" color="negative" />
                  </q-item-section>
                  <q-item-section class="text-negative">Logout</q-item-section>
                </q-item>
              </q-list>
            </q-menu>
          </q-btn>

          <q-btn
            v-else-if="!store.loading"
            unevaluated
            color="primary"
            to="/login"
            label="Get Started"
            no-caps
            class="q-px-md rounded-borders shadow-sm text-weight-bold"
          />
        </div>

        <!-- Mobile Navigation Menu -->
        <div class="lt-sm">
          <q-btn
            flat
            round
            dense
            icon="menu"
            aria-label="Open navigation menu"
            class="mobile-menu-btn"
            :class="isDarkMode ? 'text-white' : 'text-grey-8'"
          >
            <q-menu auto-close class="mobile-menu shadow-24" anchor="bottom right" self="top right">
              <!-- Mobile Profile Block if authenticated -->
              <div v-if="store.isAuth && store.user" class="q-pa-md row items-center q-gutter-x-sm bg-grey-2 border-bottom-light">
                <q-avatar size="36px" class="shadow-sm">
                  <img :src="store.user.image || '/default-avatar.png'" alt="profile" />
                </q-avatar>
                <div class="column overflow-hidden">
                  <span class="text-weight-bold text-caption ellipsis text-main">{{ store.user.name }}</span>
                  <span class="text-caption text-sub ellipsis" style="font-size: 11px">{{ store.user.email }}</span>
                </div>
              </div>
              <q-separator v-if="store.isAuth" />

              <q-list style="min-width: 220px" class="q-py-xs">
                <q-item v-if="isBlogsRoute" clickable @click="toggleLeftDrawer">
                  <q-item-section avatar>
                    <q-icon name="filter_list" size="xs" />
                  </q-item-section>
                  <q-item-section>Filters</q-item-section>
                </q-item>
                
                <q-item clickable @click="toggleDarkMode">
                  <q-item-section avatar>
                    <q-icon :name="isDarkMode ? 'light_mode' : 'dark_mode'" size="xs" />
                  </q-item-section>
                  <q-item-section>{{ isDarkMode ? 'Light mode' : 'Dark mode' }}</q-item-section>
                </q-item>
                
                <q-separator class="q-my-xs" />

                <q-item clickable to="/blogs">
                  <q-item-section avatar>
                    <q-icon name="home" size="xs" />
                  </q-item-section>
                  <q-item-section>Home</q-item-section>
                </q-item>
                <q-item v-if="store.isAuth" clickable to="/blog/saved">
                  <q-item-section avatar>
                    <q-icon name="bookmarks" size="xs" />
                  </q-item-section>
                  <q-item-section>Saved Blogs</q-item-section>
                </q-item>
                <q-item v-if="store.isAuth" clickable to="/blog/new">
                  <q-item-section avatar>
                    <q-icon name="edit_note" size="xs" />
                  </q-item-section>
                  <q-item-section>Create Blog</q-item-section>
                </q-item>
                
                <q-separator v-if="store.isAuth" class="q-my-xs" />
                
                <q-item v-if="store.isAuth" clickable to="/profile">
                  <q-item-section avatar>
                    <q-icon name="person" size="xs" />
                  </q-item-section>
                  <q-item-section>My Profile</q-item-section>
                </q-item>
                <q-item v-if="store.isAuth" clickable @click="handleLogout" class="text-negative">
                  <q-item-section avatar>
                    <q-icon name="logout" size="xs" color="negative" />
                  </q-item-section>
                  <q-item-section>Logout</q-item-section>
                </q-item>
                <q-item v-else clickable to="/login">
                  <q-item-section avatar>
                    <q-icon name="login" size="xs" color="primary" />
                  </q-item-section>
                  <q-item-section>Login</q-item-section>
                </q-item>
              </q-list>
            </q-menu>
          </q-btn>
        </div>
      </q-toolbar>
    </q-header>

    <!-- Sidebar / Drawer Filters for Blogs Feed -->
    <q-drawer
      v-if="isBlogsRoute"
      v-model="store.leftDrawerOpen"
      show-if-above
      bordered
      :width="280"
      class="sidebar-drawer"
    >
      <div class="q-pa-md q-gutter-y-lg">
        <!-- Search Section -->
        <div>
          <div class="sidebar-section-title text-caption text-bold text-uppercase text-sub tracking-wider q-mb-sm">Search</div>
          <q-input
            v-model="searchVal"
            outlined
            dense
            placeholder="Search blogs..."
            class="sidebar-search-input"
            @update:model-value="updateSearch"
          >
            <template #prepend>
              <q-icon name="search" size="xs" class="text-sub" />
            </template>
            <template v-if="searchVal" #append>
              <q-icon name="close" size="xs" class="cursor-pointer text-sub" @click="clearSearch" />
            </template>
          </q-input>
        </div>

        <q-separator />

        <!-- Categories Section -->
        <div>
          <div class="sidebar-section-title text-caption text-bold text-uppercase text-sub tracking-wider q-mb-sm">Categories</div>
          <q-list class="q-gutter-y-xs categories-list">
            <q-item
              clickable
              v-ripple
              class="sidebar-category-item"
              :active="store.category === ''"
              active-class="sidebar-category-item--active"
              @click="selectCategory('')"
            >
              <q-item-section avatar class="category-icon-section">
                <q-icon name="grid_view" size="xs" />
              </q-item-section>
              <q-item-section>All Categories</q-item-section>
            </q-item>

            <q-item
              v-for="cat in categories"
              :key="cat"
              clickable
              v-ripple
              class="sidebar-category-item"
              :active="store.category === cat"
              active-class="sidebar-category-item--active"
              @click="selectCategory(cat)"
            >
              <q-item-section avatar class="category-icon-section">
                <q-icon name="tag" size="xs" />
              </q-item-section>
              <q-item-section>{{ cat }}</q-item-section>
            </q-item>
          </q-list>
        </div>
      </div>
    </q-drawer>

    <!-- Main Container -->
    <q-page-container class="bg-app-container">
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useQuasar } from 'quasar';
import { useRoute, useRouter } from 'vue-router';
import { useAppStore, blogCategories } from '@/stores/app';

const store = useAppStore();
const route = useRoute();
const router = useRouter();
const $q = useQuasar();

const searchVal = ref(store.searchQuery);

const isBlogsRoute = computed(() => route.path === '/blogs');
const isDarkMode = computed(() => $q.dark.isActive);
const categories = blogCategories;

function toggleLeftDrawer() {
  store.leftDrawerOpen = !store.leftDrawerOpen;
}

function updateSearch(val: string | number | null) {
  store.setSearchQuery(val ? String(val) : '');
}

function clearSearch() {
  searchVal.value = '';
  store.setSearchQuery('');
}

function selectCategory(cat: string) {
  store.setCategory(cat);
}

function handleLogout() {
  store.logoutUser();
  void router.push('/blogs');
}

function toggleDarkMode() {
  const nextValue = !$q.dark.isActive;
  $q.dark.set(nextValue);
  localStorage.setItem('postly-theme', nextValue ? 'dark' : 'light');
}

onMounted(() => {
  const savedTheme = localStorage.getItem('postly-theme');
  if (savedTheme) {
    $q.dark.set(savedTheme === 'dark');
  }
  void store.initApp();
});
</script>

<style scoped>
/* Glassmorphism Header */
.navbar-blur {
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border-bottom: 1px solid var(--border-color) !important;
  transition: all 0.3s ease;
}

.navbar-blur--light {
  background: rgba(255, 255, 255, 0.82) !important;
}

.navbar-blur--dark {
  background: rgba(9, 15, 32, 0.85) !important;
  box-shadow: 0 4px 30px rgba(0, 0, 0, 0.2);
}

.main-toolbar {
  max-width: 1300px;
  margin: 0 auto;
  width: 100%;
  padding-left: 20px;
  padding-right: 20px;
}

/* App title and Logo */
.app-title {
  flex: 0 0 auto;
  transition: opacity 0.2s ease;
}

.app-title:hover {
  opacity: 0.85;
}

.brand-logo-container {
  border: 1px solid rgba(99, 102, 241, 0.15);
  box-shadow: var(--shadow-sm);
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.1), rgba(99, 102, 241, 0.2)) !important;
}

.brand-logo-icon {
  background: linear-gradient(135deg, #6366f1, #a855f7);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.logo-text {
  font-weight: 800;
  letter-spacing: -0.04em;
  background: linear-gradient(135deg, var(--text-main) 30%, var(--q-primary) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

/* Navigation Links (Pill layout) */
.nav-links {
  margin-left: 28px;
}

.nav-link-btn {
  color: var(--text-muted);
  font-weight: 500;
  border-radius: var(--radius-md);
  padding: 4px 14px;
  min-height: auto;
  font-size: 0.9rem;
}

.nav-link-btn:hover {
  background-color: rgba(99, 102, 241, 0.06) !important;
  color: var(--text-main);
}

.nav-link-btn--active {
  background-color: rgba(99, 102, 241, 0.08) !important;
  color: var(--q-primary) !important;
  font-weight: 600;
}

/* User Menu & Dropdown list */
.avatar-ring {
  background-color: var(--bg-card);
  transition: transform 0.2s ease;
}

.avatar-ring:hover {
  transform: rotate(6deg) scale(1.05);
}

.dropdown-profile-header {
  min-width: 220px;
}

.dropdown-list {
  border: none;
}

.dropdown-item {
  margin: 4px 8px;
  border-radius: var(--radius-sm);
  font-size: 0.875rem;
}

.text-negative-item:hover {
  background-color: rgba(244, 63, 94, 0.08) !important;
}

/* Drawer Filter Layout styling */
.sidebar-drawer {
  background-color: var(--bg-drawer) !important;
}

.sidebar-section-title {
  font-size: 0.72rem;
  letter-spacing: 0.08em;
  font-weight: 700;
}

.sidebar-search-input :deep(.q-field__control) {
  border-radius: var(--radius-md);
  background-color: var(--bg-card) !important;
}

.sidebar-category-item {
  border-radius: var(--radius-md);
  margin: 2px 0;
  font-size: 0.9rem;
  color: var(--text-muted);
  padding: 8px 12px;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.sidebar-category-item:hover {
  background-color: rgba(99, 102, 241, 0.05) !important;
  color: var(--text-main);
}

.sidebar-category-item--active {
  background-color: rgba(99, 102, 241, 0.08) !important;
  color: var(--q-primary) !important;
  font-weight: 600;
}

.category-icon-section {
  min-width: auto;
  padding-right: 12px;
}

.theme-toggle-btn {
  transition: transform 0.3s ease;
}

.theme-toggle-btn:hover {
  transform: rotate(30deg);
}

.mobile-menu-btn {
  border-radius: var(--radius-md);
}

/* Container Background helper */
.bg-app-container {
  background-color: var(--bg-app);
  min-height: 100vh;
}

@media (max-width: 599px) {
  .main-toolbar {
    padding-left: 14px;
    padding-right: 14px;
  }

  .logo-text {
    font-size: 1.15rem;
  }
}
</style>
