<template>
  <q-page class="flex flex-center q-pa-md bg-grey-1">
    <div v-if="loading">
      <loading-spinner />
    </div>

    <div v-else-if="user" class="profile-container full-width">
      <q-card flat borderless class="profile-card shadow-15 q-pa-lg">
        <q-card-section class="text-center">
          <div class="text-h5 text-bold font-brand text-grey-9 q-mb-md">Profile</div>
        </q-card-section>

        <q-card-section class="column items-center q-gutter-y-md">
          <!-- Avatar -->
          <q-avatar size="110px" class="profile-avatar shadow-4">
            <img :src="user.image || '/default-avatar.png'" alt="profile" />
          </q-avatar>

          <!-- User Name -->
          <div class="text-h6 text-weight-bold text-grey-9 text-center">
            {{ user.name }}
          </div>

          <!-- User Bio -->
          <div v-if="user.bio" class="text-body1 text-grey-7 text-center max-width-md">
            {{ user.bio }}
          </div>

          <!-- Social Links -->
          <div class="row q-gutter-x-md q-mt-sm">
            <a
              v-if="user.instagram"
              :href="user.instagram"
              target="_blank"
              class="social-link instagram"
            >
              <q-icon name="camera_alt" size="sm" />
            </a>
            <a
              v-if="user.facebook"
              :href="user.facebook"
              target="_blank"
              class="social-link facebook"
            >
              <q-icon name="facebook" size="sm" />
            </a>
            <a
              v-if="user.linkedin"
              :href="user.linkedin"
              target="_blank"
              class="social-link linkedin"
            >
              <q-icon name="work" size="sm" />
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
  max-width: 550px;
}

.profile-card {
  border-radius: 16px;
  background: white;
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.font-brand {
  font-family: 'Outfit', 'Inter', sans-serif;
  letter-spacing: -0.5px;
}

.profile-avatar {
  border: 4px solid #f3f4f6;
}

.social-link {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  color: white;
  text-decoration: none;
  transition: transform 0.2s ease;
}

.social-link:hover {
  transform: translateY(-2px);
}

.instagram {
  background: radial-gradient(
    circle at 30% 107%,
    #fdf497 0%,
    #fdf497 5%,
    #fd5949 45%,
    #d6249f 60%,
    #285aeb 90%
  );
}

.facebook {
  background: #1877f2;
}

.linkedin {
  background: #0a66c2;
}

.max-width-md {
  max-width: 320px;
}
</style>
