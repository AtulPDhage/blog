<template>
  <q-page class="q-py-lg q-px-md">
    <!-- Loading State -->
    <div v-if="store.loading || !store.blogs || !store.savedBlogs" class="flex flex-center">
      <loading-spinner />
    </div>

    <div v-else class="container mx-auto max-width-xl">
      <!-- Title -->
      <h1 class="text-h4 text-bold q-my-md font-brand text-grey-9">Saved Blogs</h1>

      <!-- Filtered Saved Blogs Grid -->
      <div v-if="filteredBlogs.length > 0" class="row q-col-gutter-md">
        <div v-for="blog in filteredBlogs" :key="blog.id" class="col-12 col-sm-6 col-md-4 col-lg-3">
          <blog-card
            :id="blog.id"
            :image="blog.image"
            :title="blog.title"
            :desc="blog.description"
            :time="blog.created_at"
          />
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="column items-center q-py-xl text-grey-6 text-center">
        <q-icon name="bookmark_border" size="64px" class="q-mb-md" />
        <p class="text-h6 text-weight-regular">No saved blogs yet!</p>
        <q-btn
          color="primary"
          to="/blogs"
          label="Browse Blogs"
          no-caps
          class="q-mt-sm rounded-borders"
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
  if (!store.blogs || !store.savedBlogs) return [];
  return store.blogs.filter((blog) =>
    store.savedBlogs!.some((saved) => saved.blogid === String(blog.id)),
  );
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
  await store.fetchBlogs();
  await store.getSavedBlogs();
});
</script>

<style scoped>
.max-width-xl {
  max-width: 1200px;
  margin: 0 auto;
}

.font-brand {
  font-family: 'Outfit', 'Inter', sans-serif;
  letter-spacing: -0.5px;
}
</style>
