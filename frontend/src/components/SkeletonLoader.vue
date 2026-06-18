<template>
  <q-card flat class="skeleton-card">
    <div class="img-skeleton-wrapper">
      <q-skeleton height="200px" square animation="wave" />
    </div>

    <q-card-section class="q-pa-md">
      <!-- Date Skeleton -->
      <q-skeleton type="text" width="80px" class="q-mb-sm" animation="wave" />

      <!-- Title Skeleton -->
      <q-skeleton type="text" height="24px" class="q-mb-md" animation="wave" />

      <!-- Description Skeletons -->
      <q-skeleton type="text" class="q-mb-xs" animation="wave" />
      <q-skeleton type="text" width="70%" animation="wave" />
    </q-card-section>
  </q-card>
</template>

<script setup lang="ts">
// Reusable modern shimmer skeleton card loader
</script>

<style scoped>
.skeleton-card {
  border-radius: var(--radius-lg);
  background-color: var(--bg-card);
  border: 1px solid var(--border-color);
  overflow: hidden;
  height: 100%;
}
.img-skeleton-wrapper {
  background-color: var(--border-color);
}
</style>
