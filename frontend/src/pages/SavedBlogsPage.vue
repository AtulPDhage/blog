<template>
  <q-page class="q-py-xl q-px-md bg-app-container">
    <!-- Loading State -->
    <div v-if="store.loading || !store.blogs || !store.savedBlogs" class="flex flex-center initial-loading-container">
      <loading-spinner />
    </div>

    <div v-else class="container mx-auto max-width-xl">
      <!-- Title -->
      <h1 class="text-h4 text-bold q-my-none q-mb-xl font-brand text-main">Saved Bookmarks</h1>

      <!-- Filtered Saved Blogs Grid -->
      <div v-if="filteredBlogs.length > 0" class="blogs-grid">
        <div v-for="blog in filteredBlogs" :key="blog.id" class="blog-grid-item">
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

      <!-- Empty State -->
      <div v-else class="empty-container column items-center justify-center text-center q-py-xl">
        <div class="empty-icon-wrapper q-mb-lg flex flex-center">
          <q-icon name="bookmark_border" size="40px" class="text-primary" />
        </div>
        <h2 class="text-h5 text-bold q-mt-none q-mb-xs font-brand text-main">No bookmarks saved</h2>
        <p class="text-body1 text-muted max-width-xs q-mb-lg">
          Articles you bookmark will be displayed here for quick reading later.
        </p>
        <q-btn
          unevaluated
          color="primary"
          to="/blogs"
          label="Browse Articles"
          no-caps
          class="q-px-lg rounded-borders text-weight-bold shadow-sm"
        />
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAppStore } from '@/stores/app';
import BlogCard from '@/components/BlogCard.vue';
import LoadingSpinner from '@/components/LoadingSpinner.vue';

const store = useAppStore();
const router = useRouter();

const filteredBlogs = computed(() => {
  return store.savedBlogs || [];
});

onMounted(async () => {
  if (!store.isAuth && !store.loading) {
    void router.replace('/login');
    return;
  }

  // Hydrate user and blogs if needed
  if (store.loading) {
    await store.fetchUser();
  }
  await store.getSavedBlogs();
});
</script>

<style scoped>
.initial-loading-container {
  min-height: 70vh;
}

.max-width-xl {
  max-width: 1300px;
  margin: 0 auto;
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
</style>
