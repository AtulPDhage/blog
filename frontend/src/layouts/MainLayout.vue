<template>
  <q-layout view="lHh Lpr lFf">
    <!-- Header with Glassmorphism -->
    <q-header
      borderless
      class="navbar-blur"
      :class="isDarkMode ? 'navbar-blur--dark text-white' : 'text-grey-9'"
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

        <q-toolbar-title
          class="app-title text-bold text-h6 font-brand cursor-pointer"
          @click="router.push('/blogs')"
        >
          Postly
        </q-toolbar-title>

        <!-- Desktop Navigation Links -->
        <div class="nav-links gt-xs row items-center no-wrap q-gutter-sm">
          <q-btn
            flat
            no-caps
            label="Home"
            to="/blogs"
            class="nav-link"
            active-class="nav-link-active"
          />
          <q-btn
            v-if="store.isAuth"
            flat
            no-caps
            label="Saved Blogs"
            to="/blog/saved"
            class="nav-link"
            active-class="nav-link-active"
          />
          <q-btn
            v-if="store.isAuth"
            flat
            no-caps
            label="Create Blog"
            to="/blog/new"
            class="nav-link"
            active-class="nav-link-active"
          />
        </div>

        <q-space />

        <!-- Desktop Right Actions -->
        <div class="nav-actions gt-xs row items-center no-wrap q-gutter-sm">
          <q-btn
            flat
            round
            dense
            :icon="isDarkMode ? 'light_mode' : 'dark_mode'"
            :color="isDarkMode ? 'amber-4' : 'grey-7'"
            @click="toggleDarkMode"
          >
            <q-tooltip>{{ isDarkMode ? 'Light mode' : 'Dark mode' }}</q-tooltip>
          </q-btn>

          <q-btn-dropdown
            v-if="store.isAuth && store.user"
            flat
            round
            dense
            no-caps
            class="profile-dropdown"
          >
            <template #label>
              <q-avatar size="36px" class="shadow-1">
                <img :src="store.user.image" alt="profile" />
              </q-avatar>
            </template>

            <q-list style="min-width: 180px" class="q-py-xs">
              <q-item clickable v-close-popup to="/profile">
                <q-item-section avatar>
                  <q-icon name="person" size="xs" />
                </q-item-section>
                <q-item-section>Profile</q-item-section>
              </q-item>
              <q-item clickable v-close-popup @click="handleLogout">
                <q-item-section avatar>
                  <q-icon name="logout" size="xs" color="negative" />
                </q-item-section>
                <q-item-section class="text-negative">Logout</q-item-section>
              </q-item>
            </q-list>
          </q-btn-dropdown>

          <q-btn
            v-else-if="!store.loading"
            flat
            round
            dense
            icon="login"
            to="/login"
            color="primary"
            class="q-ml-sm"
          >
            <q-tooltip>Login</q-tooltip>
          </q-btn>
        </div>

        <!-- Mobile Navigation Menu -->
        <div class="lt-sm">
          <q-btn
            flat
            round
            dense
            icon="menu"
            aria-label="Open navigation menu"
            :class="isDarkMode ? 'text-white' : 'text-grey-7'"
          >
            <q-menu auto-close class="mobile-menu shadow-15" anchor="bottom right" self="top right">
              <q-list style="min-width: 220px">
                <q-item v-if="isBlogsRoute" clickable @click="toggleLeftDrawer">
                  <q-item-section avatar>
                    <q-icon name="filter_list" size="xs" />
                  </q-item-section>
                  <q-item-section>Filter Blogs</q-item-section>
                </q-item>
                <q-separator v-if="isBlogsRoute" />
                <q-item clickable @click="toggleDarkMode">
                  <q-item-section avatar>
                    <q-icon :name="isDarkMode ? 'light_mode' : 'dark_mode'" size="xs" />
                  </q-item-section>
                  <q-item-section>{{ isDarkMode ? 'Light mode' : 'Dark mode' }}</q-item-section>
                </q-item>
                <q-separator />
                <q-item clickable to="/blogs">
                  <q-item-section>Home</q-item-section>
                </q-item>
                <q-item v-if="store.isAuth" clickable to="/blog/saved">
                  <q-item-section>Saved Blogs</q-item-section>
                </q-item>
                <q-item v-if="store.isAuth" clickable to="/blog/new">
                  <q-item-section>Create Blog</q-item-section>
                </q-item>
                <q-separator v-if="store.isAuth" />
                <q-item v-if="store.isAuth" clickable to="/profile">
                  <q-item-section>Profile</q-item-section>
                </q-item>
                <q-item v-if="store.isAuth" clickable @click="handleLogout" class="text-negative">
                  <q-item-section>Logout</q-item-section>
                </q-item>
                <q-item v-else clickable to="/login">
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
      class="bg-grey-1"
    >
      <div class="q-pa-md q-gutter-y-md">
        <div class="text-subtitle1 text-bold font-brand text-grey-8">Search</div>
        <q-input
          v-model="searchVal"
          outlined
          dense
          placeholder="Search blogs..."
          bg-color="white"
          class="search-input"
          @update:model-value="updateSearch"
        >
          <template #append>
            <q-icon name="search" size="xs" />
          </template>
        </q-input>

        <q-separator class="q-my-md" />

        <div class="text-subtitle1 text-bold font-brand text-grey-8">Categories</div>
        <q-list class="q-gutter-y-xs">
          <q-item
            clickable
            v-ripple
            class="category-item"
            :active="store.category === ''"
            active-class="category-active"
            @click="selectCategory('')"
          >
            <q-item-section avatar>
              <q-icon name="grid_view" size="xs" />
            </q-item-section>
            <q-item-section>All Categories</q-item-section>
          </q-item>

          <q-item
            v-for="cat in categories"
            :key="cat"
            clickable
            v-ripple
            class="category-item"
            :active="store.category === cat"
            active-class="category-active"
            @click="selectCategory(cat)"
          >
            <q-item-section avatar>
              <q-icon name="tag" size="xs" />
            </q-item-section>
            <q-item-section>{{ cat }}</q-item-section>
          </q-item>
        </q-list>
      </div>
    </q-drawer>

    <!-- Main Container -->
    <q-page-container class="bg-grey-1">
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
.navbar-blur {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
}

.navbar-blur--dark {
  background: rgba(15, 23, 42, 0.94);
  border-bottom-color: rgba(148, 163, 184, 0.18);
  box-shadow: 0 1px 18px rgba(0, 0, 0, 0.22);
}

.font-brand {
  font-family: 'Outfit', 'Inter', sans-serif;
  letter-spacing: -0.5px;
}

.main-toolbar {
  width: 100%;
  min-width: 0;
  padding-left: 16px;
  padding-right: 16px;
}

.app-title {
  flex: 0 0 auto;
  min-width: max-content;
}

.nav-links {
  flex: 0 1 auto;
  min-width: 0;
  margin-left: 18px;
}

.nav-actions {
  flex: 0 0 auto;
  min-width: max-content;
}

.filter-menu-btn {
  flex: 0 0 auto;
}

.nav-link {
  color: #4a5568;
  font-weight: 500;
  border-radius: 8px;
  transition: all 0.2s ease;
  min-width: 0;
}

.nav-link:hover {
  background: rgba(0, 0, 0, 0.04);
  color: #1a202c;
}

.navbar-blur--dark .nav-link {
  color: #cbd5e1;
}

.navbar-blur--dark .nav-link:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #f8fafc;
}

.nav-link-active {
  color: var(--q-primary) !important;
  background: rgba(25, 118, 210, 0.08);
}

.navbar-blur--dark .nav-link-active {
  color: #60a5fa !important;
  background: rgba(96, 165, 250, 0.16);
}

.category-item {
  border-radius: 8px;
  color: #4a5568;
  transition: all 0.2s ease;
}

.category-item:hover {
  background: rgba(0, 0, 0, 0.03);
}

.category-active {
  background: rgba(25, 118, 210, 0.08) !important;
  color: #1976d2 !important;
  font-weight: 600;
}

.search-input :deep(.q-field__control) {
  border-radius: 8px;
}

.profile-dropdown :deep(.q-btn__content) {
  min-height: auto;
}

.mobile-menu {
  border-radius: 10px;
}

@media (max-width: 599px) {
  .main-toolbar {
    padding-left: 12px;
    padding-right: 12px;
  }

  .app-title {
    flex: 1 1 auto;
  }
}
</style>
