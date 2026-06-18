<template>
  <q-page class="home-page q-py-xl q-px-md">
    <!-- Initial Loading State with Skeleton Grids -->
    <div v-if="store.blogLoading && (!store.blogs || store.blogs.length === 0)" class="home-container max-width-xl">
      <div class="home-header row justify-between items-center q-mb-xl">
        <h1 class="text-h4 text-bold q-my-none font-brand text-main">
          {{ store.category ? `${store.category} Blogs` : 'Latest Blogs' }}
        </h1>
        <q-btn
          flat
          icon="filter_list"
          label="Filters"
          no-caps
          class="filter-btn text-weight-bold"
          @click="toggleFilter"
        />
      </div>
      <div class="blogs-grid">
        <skeleton-loader v-for="n in 8" :key="'init-skeleton-' + n" />
      </div>
    </div>

    <!-- Main Blogs Content Grid -->
    <div v-else class="home-container max-width-xl">
      <!-- Title & Filter Toggle Header -->
      <div class="home-header row justify-between items-center q-mb-xl">
        <h1 class="text-h4 text-bold q-my-none font-brand text-main">
          {{ store.category ? `${store.category} Blogs` : 'Latest Blogs' }}
        </h1>
        <q-btn
          flat
          icon="filter_list"
          label="Filters"
          no-caps
          class="filter-btn text-weight-bold"
          @click="toggleFilter"
        />
      </div>

      <!-- Empty Feed State -->
      <div
        v-if="!store.blogs || store.blogs.length === 0"
        class="empty-container column items-center justify-center text-center q-py-xl"
      >
        <div class="empty-icon-wrapper q-mb-lg flex flex-center">
          <q-icon name="article" size="40px" class="text-primary" />
        </div>
        <h2 class="text-h5 text-bold q-mt-none q-mb-xs font-brand text-main">No blogs found</h2>
        <p class="text-body1 text-muted max-width-xs q-mb-lg">
          We couldn't find any blogs matching your search query or selected category.
        </p>
        <q-btn
          v-if="store.isAuth"
          unevaluated
          color="primary"
          to="/blog/new"
          label="Write your first post"
          no-caps
          class="q-px-lg rounded-borders text-weight-bold"
        />
        <q-btn
          v-else
          outlined
          color="primary"
          label="Go back to all blogs"
          no-caps
          class="q-px-lg rounded-borders text-weight-bold"
          @click="resetFilters"
        />
      </div>

      <!-- Infinite Scroll Cards Feed Grid -->
      <q-infinite-scroll v-else @load="loadMore" :offset="250" ref="infiniteScrollRef">
        <div class="blogs-grid">
          <div v-for="blog in store.blogs" :key="blog.id" class="blog-grid-item">
            <blog-card
              :id="blog.id"
              :image="blog.image"
              :title="blog.title"
              :desc="blog.description"
              :time="blog.created_at"
              :category="blog.category"
              :views="blog.views"
            />
          </div>
        </div>

        <template v-slot:loading>
          <div class="blogs-grid q-mt-lg">
            <skeleton-loader v-for="n in 4" :key="'load-skeleton-' + n" />
          </div>
        </template>
      </q-infinite-scroll>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue';
import { useAppStore } from '@/stores/app';
import BlogCard from '@/components/BlogCard.vue';
import SkeletonLoader from '@/components/SkeletonLoader.vue';

const store = useAppStore();

const limit = 12;
const hasMore = ref(true);
const infiniteScrollRef = ref<{ resume: () => void; trigger: () => void } | null>(null);

function toggleFilter() {
  store.leftDrawerOpen = !store.leftDrawerOpen;
}

function resetFilters() {
  store.setCategory('');
  store.setSearchQuery('');
}

async function loadMore(index: number, done: (stop?: boolean) => void) {
  if (!hasMore.value) {
    done(true);
    return;
  }

  const offset = store.blogs ? store.blogs.length : 0;
  const fetched = await store.fetchBlogs(limit, offset, true);

  if (fetched.length < limit) {
    hasMore.value = false;
    done(true);
  } else {
    done();
  }
}

async function resetPagination() {
  hasMore.value = true;
  store.blogs = [];

  const fetched = await store.fetchBlogs(limit, 0, false);
  if (fetched.length < limit) {
    hasMore.value = false;
  }

  void nextTick(() => {
    if (infiniteScrollRef.value) {
      infiniteScrollRef.value.resume();
    }
  });
}

watch(
  () => [store.searchQuery, store.category],
  () => {
    void resetPagination();
  },
);
</script>

<style scoped>
.home-page {
  overflow-x: hidden;
}

.home-container {
  width: 100%;
  min-width: 0;
}

.max-width-xl {
  max-width: 1300px;
  margin: 0 auto;
}

.home-header {
  gap: 12px;
}

.filter-btn {
  border: 1px solid var(--border-color);
  background-color: var(--bg-card);
  color: var(--text-main);
  border-radius: var(--radius-md);
  padding: 8px 16px;
  transition: all 0.2s ease;
}

.filter-btn:hover {
  background-color: rgba(99, 102, 241, 0.05);
  border-color: var(--q-primary);
  color: var(--q-primary);
}

/* Responsive CSS Grid */
.blogs-grid {
  display: grid;
  grid-template-columns: repeat(1, minmax(0, 1fr));
  gap: 28px;
  width: 100%;
}

.blog-grid-item {
  min-width: 0;
}

/* Empty State Polish */
.empty-container {
  min-height: 40vh;
}

.empty-icon-wrapper {
  width: 80px;
  height: 80px;
  border-radius: var(--radius-xl);
  background-color: rgba(99, 102, 241, 0.08);
  border: 1px solid rgba(99, 102, 241, 0.15);
}

.max-width-xs {
  max-width: 340px;
}

@media (min-width: 600px) {
  .blogs-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (min-width: 960px) {
  .blogs-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@media (min-width: 1280px) {
  .blogs-grid {
    grid-template-columns: repeat(4, minmax(0, 1fr));
  }
}

@media (max-width: 599px) {
  .home-header {
    align-items: flex-start;
    flex-direction: column;
    margin-bottom: 24px;
  }

  .filter-btn {
    align-self: stretch;
    text-align: center;
  }
}
</style>
