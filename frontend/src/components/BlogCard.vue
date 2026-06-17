<template>
  <q-card flat class="blog-card cursor-pointer bg-transparent" @click="router.push(`/blog/${id}`)">
    <div class="img-wrapper overflow-hidden rounded-borders shadow-1">
      <q-img :src="image" :ratio="16 / 10" fit="cover" class="blog-img" />
    </div>

    <q-card-section class="q-px-none q-pt-sm">
      <div class="flex justify-center items-center text-grey-6 text-caption q-gutter-x-xs">
        <q-icon name="calendar_today" size="14px" />
        <span>{{ formattedDate }}</span>
      </div>

      <div class="text-weight-bold text-h6 text-center text-grey-9 q-mt-xs q-mb-xs title-ellipsis">
        {{ title }}
      </div>

      <div class="text-body2 text-center text-grey-7 desc-ellipsis">
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
  return props.desc.length > 30 ? `${props.desc.slice(0, 30)}...` : props.desc;
});
</script>

<style scoped>
.blog-card {
  transition:
    transform 0.25s ease,
    box-shadow 0.25s ease;
  border-radius: 12px;
}

.img-wrapper {
  overflow: hidden;
  border-radius: 12px;
}

.blog-img {
  transition: transform 0.4s ease;
}

.blog-card:hover .blog-img {
  transform: scale(1.06);
}

.title-ellipsis {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  display: block;
}

.desc-ellipsis {
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}
</style>
