<template>
  <q-card
    flat
    class="blog-card cursor-pointer"
    @click="router.push(`/blog/${id}`)"
  >
    <div class="img-wrapper overflow-hidden relative-position shadow-sm">
      <q-img :src="image || '/placeholder-cover.jpg'" :ratio="16 / 10" fit="cover" class="blog-img" />
      <div v-if="category" class="category-badge-overlay">
        <span class="category-badge font-brand">{{ category }}</span>
      </div>
    </div>

    <q-card-section class="q-pa-md card-content">
      <div class="row items-center text-sub text-caption q-mb-xs q-gutter-x-xs font-brand text-medium">
        <q-icon name="schedule" size="14px" />
        <span>{{ formattedDate }}</span>
      </div>

      <h3 class="text-weight-bold text-subtitle1 q-mt-none q-mb-sm title-ellipsis text-main">
        {{ title }}
      </h3>

      <div class="text-body2 text-muted desc-ellipsis leading-relaxed">
        {{ truncatedDesc }}
      </div>
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRouter } from 'vue-router';

const props = defineProps<{
  image: string;
  title: string;
  desc: string;
  id: string;
  time: string;
  category?: string;
}>();

const router = useRouter();

const formattedDate = computed(() => {
  if (!props.time) return '';
  const date = new Date(props.time);
  const day = String(date.getDate()).padStart(2, '0');
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const year = date.getFullYear();
  return `${day}-${month}-${year}`;
});

const truncatedDesc = computed(() => {
  if (!props.desc) return '';
  return props.desc.length > 80 ? `${props.desc.slice(0, 80)}...` : props.desc;
});
</script>

<style scoped>
.blog-card {
  height: 100%;
  display: flex;
  flex-direction: column;
  border-radius: var(--radius-lg);
  background-color: var(--bg-card);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-sm);
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
}

.blog-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-md);
  border-color: rgba(99, 102, 241, 0.2);
}

.img-wrapper {
  overflow: hidden;
  border-bottom: 1px solid var(--border-color);
}

.blog-img {
  transition: transform 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

.blog-card:hover .blog-img {
  transform: scale(1.04);
}

.category-badge-overlay {
  position: absolute;
  top: 12px;
  left: 12px;
  z-index: 2;
}

.category-badge {
  background: rgba(15, 23, 42, 0.75);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  color: #ffffff;
  font-size: 0.75rem;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: var(--radius-sm);
  border: 1px solid rgba(255, 255, 255, 0.15);
  text-transform: capitalize;
}

.card-content {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
}

.title-ellipsis {
  font-family: var(--font-brand);
  line-height: 1.35;
  margin-top: 4px;
  font-weight: 700;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: normal;
  height: 2.7em; /* Consistent spacing */
}

.desc-ellipsis {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.5;
  color: var(--text-muted);
}
</style>
