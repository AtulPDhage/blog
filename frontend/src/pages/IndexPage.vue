<template>
  <q-page class="home-page q-py-lg q-px-md">
    <!-- Main Loader -->
    <!-- Blogs Content Grid -->
    <div v-if="store.blogLoading && (!store.blogs || store.blogs.length === 0)">
      <loading-spinner />
    </div>

    <div v-else class="home-container max-width-xl">
      <!-- Title & Filter Toggle Header -->
      <div class="home-header row justify-between items-center q-mb-lg">
        <h1 class="text-h4 text-bold q-my-none font-brand text-grey-9">
          {{ store.category ? `${store.category} Blogs` : 'Latest Blogs' }}
        </h1>
        <q-btn
          flat
          color="primary"
          icon="filter_list"
          label="Filter Blogs"
          no-caps
          class="q-px-md filter-btn rounded-borders"
          @click="toggleFilter"
        />
      </div>

      <!-- Empty Feed State -->
      <div
        v-if="!store.blogs || store.blogs.length === 0"
        class="column items-center q-py-xl text-grey-6 text-center"
      >
        <q-icon name="article" size="64px" class="q-mb-md" />
        <p class="text-h6 text-weight-regular">No blogs found yet!</p>
        <q-btn
          v-if="store.isAuth"
          color="primary"
          to="/blog/new"
          label="Write the first one"
          no-caps
          class="q-mt-sm rounded-borders"
        />
      </div>

      <!-- Infinite Scroll Cards Feed Grid -->
      <q-infinite-scroll v-else @load="loadMore" :offset="250" ref="infiniteScrollRef">
        <q-virtual-scroll
          type="list"
          :items="blogRows"
          :virtual-scroll-item-size="380"
          scroll-target="body"
        >
          <template v-slot="{ item: row, index }">
            <div :key="index" class="blogs-grid q-mb-md">
              <div v-for="blog in row" :key="blog.id" class="blog-grid-item">
                <blog-card
                  :id="blog.id"
                  :image="blog.image"
                  :title="blog.title"
                  :desc="blog.description"
                  :time="blog.created_at"
                />
              </div>
            </div>
          </template>
        </q-virtual-scroll>
        <template v-slot:loading>
          <div class="row justify-center q-my-md">
            <q-spinner-dots color="primary" size="40px" />
          </div>
        </template>
      </q-infinite-scroll>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue';
import { useQuasar } from 'quasar';
import { useAppStore } from '@/stores/app';
import BlogCard from '@/components/BlogCard.vue';
import LoadingSpinner from '@/components/LoadingSpinner.vue';

const store = useAppStore();
const $q = useQuasar();

const limit = 12;
const hasMore = ref(true);
const infiniteScrollRef = ref<{ resume: () => void; trigger: () => void } | null>(null);

const columnsCount = computed(() => {
  if ($q.screen.gt.md) return 4;
  if ($q.screen.gt.sm) return 3;
  if ($q.screen.gt.xs) return 2;
  return 1;
});

const blogRows = computed(() => {
  const blogs = store.blogs || [];
  const size = columnsCount.value;
  const chunked = [];
  for (let i = 0; i < blogs.length; i += size) {
    chunked.push(blogs.slice(i, i + size));
  }
  return chunked;
});

function toggleFilter() {
  store.leftDrawerOpen = !store.leftDrawerOpen;
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
  max-width: 1200px;
  margin: 0 auto;
}

.home-header {
  gap: 12px;
}

.font-brand {
  font-family: 'Outfit', 'Inter', sans-serif;
  letter-spacing: -0.5px;
}

.filter-btn {
  border: 1px solid rgba(25, 118, 210, 0.2);
  background: rgba(25, 118, 210, 0.04);
}

.blogs-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 16px;
  width: 100%;
  min-width: 0;
}

.blog-grid-item {
  min-width: 0;
}

@media (min-width: 1440px) {
  .blogs-grid {
    grid-template-columns: repeat(4, minmax(0, 1fr));
  }
}

@media (max-width: 599px) {
  .home-header {
    align-items: flex-start;
    flex-direction: column;
  }

  .filter-btn {
    align-self: flex-start;
  }
}
</style>
