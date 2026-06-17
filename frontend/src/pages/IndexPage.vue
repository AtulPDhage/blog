<template>
  <q-page class="q-py-lg q-px-md">
    <!-- Main Loader -->
    <div v-if="store.loading">
      <loading-spinner />
    </div>

    <div v-else class="container mx-auto max-width-xl">
      <!-- Title & Filter Toggle Header -->
      <div class="row justify-between items-center q-mb-lg">
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

      <!-- Blogs Content Grid -->
      <div v-if="store.blogLoading">
        <loading-spinner />
      </div>

      <div v-else>
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

        <!-- Cards Feed Grid -->
        <div v-else class="row q-col-gutter-md">
          <div v-for="blog in store.blogs" :key="blog.id" class="col-12 col-sm-6 col-md-4 col-lg-3">
            <blog-card
              :id="blog.id"
              :image="blog.image"
              :title="blog.title"
              :desc="blog.description"
              :time="blog.created_at"
            />
          </div>
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useAppStore } from '@/stores/app';
import BlogCard from '@/components/BlogCard.vue';
import LoadingSpinner from '@/components/LoadingSpinner.vue';

const store = useAppStore();

function toggleFilter() {
  store.leftDrawerOpen = !store.leftDrawerOpen;
}

onMounted(() => {
  void store.fetchBlogs();
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

.filter-btn {
  border: 1px solid rgba(25, 118, 210, 0.2);
  background: rgba(25, 118, 210, 0.04);
}
</style>
