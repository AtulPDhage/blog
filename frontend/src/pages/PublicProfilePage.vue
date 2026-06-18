<template>
  <q-page class="flex flex-center q-pa-md bg-app-container">
    <div v-if="loading" class="initial-loading">
      <loading-spinner />
    </div>

    <div v-else-if="user" class="profile-container full-width">
      <q-card flat class="profile-card shadow-lg overflow-hidden">
        <!-- Dashboard Top Cover Banner (Gradient) -->
        <div class="profile-cover-banner"></div>

        <q-card-section class="column items-center profile-content-section relative-position q-px-lg q-pb-xl">
          <!-- Avatar -->
          <div class="avatar-wrapper shadow-md">
            <q-avatar size="120px" class="profile-avatar">
              <img :src="user.image || '/default-avatar.png'" alt="profile" />
            </q-avatar>
          </div>

          <!-- User Name -->
          <div class="text-h5 text-bold font-brand text-main text-center q-mt-md q-mb-xs">
            {{ user.name }}
          </div>

          <!-- User Bio -->
          <div v-if="user.bio" class="text-body2 text-muted text-center max-width-md q-mb-lg leading-relaxed">
            {{ user.bio }}
          </div>
          <div v-else class="text-body2 text-sub italic text-center max-width-md q-mb-lg">
            No bio details shared by this author.
          </div>

          <!-- Social Links -->
          <div class="row q-gutter-x-md">
            <a
              v-if="user.instagram"
              :href="user.instagram"
              target="_blank"
              class="social-icon-btn instagram-btn"
            >
              <q-icon name="camera_alt" size="xs" />
            </a>
            <a
              v-if="user.facebook"
              :href="user.facebook"
              target="_blank"
              class="social-icon-btn facebook-btn"
            >
              <q-icon name="facebook" size="xs" />
            </a>
            <a
              v-if="user.linkedin"
              :href="user.linkedin"
              target="_blank"
              class="social-icon-btn linkedin-btn"
            >
              <q-icon name="work" size="xs" />
            </a>
          </div>
        </q-card-section>
      </q-card>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { user_service } from '@/boot/axios';
import type { User } from '../stores/app';
import axios from 'axios';
import LoadingSpinner from '@/components/LoadingSpinner.vue';

const route = useRoute();
const user = ref<User | null>(null);
const loading = ref(true);

async function fetchUser() {
  const userId = route.params.id as string;
  if (!userId) return;

  loading.value = true;
  try {
    const { data } = await axios.get(`${user_service}/api/v1/user/${userId}`);
    user.value = data.user;
  } catch (error) {
    console.error('Error fetching public profile:', error);
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  void fetchUser();
});
</script>

<style scoped>
.profile-container {
  max-width: 600px;
}

.profile-card {
  border-radius: var(--radius-lg);
  background-color: var(--bg-card);
  border: 1px solid var(--border-color);
}

/* cover banner layout styling */
.profile-cover-banner {
  height: 140px;
  background: linear-gradient(135deg, #6366f1, #a855f7, #ec4899);
  width: 100%;
}

.profile-content-section {
  margin-top: -60px; /* pull avatar up overlap cover */
}

.avatar-wrapper {
  position: relative;
  border-radius: 50%;
  border: 5px solid var(--bg-card);
  overflow: hidden;
  background-color: var(--bg-card);
}

.profile-avatar {
  border: none;
}

/* Social Icon Buttons styling */
.social-icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 38px;
  height: 38px;
  border-radius: var(--radius-md);
  color: var(--text-muted);
  border: 1px solid var(--border-color);
  background-color: var(--bg-card);
  text-decoration: none;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.social-icon-btn:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-sm);
  color: white;
}

.instagram-btn:hover {
  background: radial-gradient(
    circle at 30% 107%,
    #fdf497 0%,
    #fdf497 5%,
    #fd5949 45%,
    #d6249f 60%,
    #285aeb 90%
  );
  border-color: transparent;
}

.facebook-btn:hover {
  background-color: #1877f2;
  border-color: #1877f2;
}

.linkedin-btn:hover {
  background-color: #0a66c2;
  border-color: #0a66c2;
}

.max-width-md {
  max-width: 400px;
}

.initial-loading {
  z-index: 2;
}
</style>
