<template>
  <q-page class="flex flex-center q-pa-md bg-app-container">
    <div v-if="store.loading || localLoading" class="initial-loading">
      <loading-spinner />
    </div>

    <div v-else-if="store.user" class="profile-container full-width">
      <q-card flat class="profile-card shadow-lg overflow-hidden">
        <!-- Dashboard Top Cover Banner (Gradient) -->
        <div class="profile-cover-banner"></div>

        <q-card-section class="column items-center profile-content-section relative-position q-px-lg q-pb-xl">
          <!-- Clickable Overlapping Avatar with Camera Hover Edit Icon -->
          <div class="avatar-wrapper shadow-md cursor-pointer" @click="triggerFileInput">
            <q-avatar size="120px" class="profile-avatar">
              <img :src="store.user.image || '/default-avatar.png'" alt="profile" />
            </q-avatar>
            <div class="avatar-overlay flex flex-center">
              <q-icon name="photo_camera" color="white" size="sm" />
            </div>
            <input
              ref="fileInput"
              type="file"
              accept="image/*"
              class="hidden"
              @change="handleImageUpload"
            />
          </div>

          <!-- User Info Details -->
          <div class="text-h5 text-bold font-brand text-main text-center q-mt-md q-mb-xs">
            {{ store.user.name }}
          </div>
          <div class="text-caption text-sub text-center q-mb-md font-brand">
            {{ store.user.email }}
          </div>

          <!-- User Bio -->
          <div v-if="store.user.bio" class="text-body2 text-muted text-center max-width-md q-mb-lg leading-relaxed">
            {{ store.user.bio }}
          </div>
          <div v-else class="text-body2 text-sub italic text-center max-width-md q-mb-lg">
            No bio written yet. Add one in your profile settings!
          </div>

          <!-- Social Links (Minimalist slate icons with custom color transitions) -->
          <div class="row q-gutter-x-md q-mb-lg">
            <a
              v-if="store.user.instagram"
              :href="store.user.instagram"
              target="_blank"
              class="social-icon-btn instagram-btn"
            >
              <q-icon name="camera_alt" size="xs" />
            </a>
            <a
              v-if="store.user.facebook"
              :href="store.user.facebook"
              target="_blank"
              class="social-icon-btn facebook-btn"
            >
              <q-icon name="facebook" size="xs" />
            </a>
            <a
              v-if="store.user.linkedin"
              :href="store.user.linkedin"
              target="_blank"
              class="social-icon-btn linkedin-btn"
            >
              <q-icon name="work" size="xs" />
            </a>
          </div>

          <!-- Dashboard Action Controls -->
          <div class="row q-gutter-md">
            <q-btn
              unevaluated
              color="primary"
              label="Create Post"
              icon="add"
              to="/blog/new"
              no-caps
              class="rounded-borders q-px-lg text-weight-bold shadow-sm"
            />
            <q-btn
              outlined
              label="Edit Profile"
              icon="settings"
              no-caps
              class="rounded-borders q-px-lg profile-settings-btn text-weight-bold text-muted"
              @click="openEditDialog"
            />
          </div>
        </q-card-section>
      </q-card>
    </div>

    <!-- Edit Profile Dialog -->
    <q-dialog v-model="editDialogOpen" backdrop-filter="blur(4px)">
      <q-card style="width: 460px; max-width: 90vw" class="rounded-borders q-pa-md shadow-24 edit-profile-modal">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6 text-bold font-brand text-main">Edit Profile Settings</div>
          <q-space />
          <q-btn flat round dense icon="close" class="text-sub" v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-y-md q-pt-md">
          <div>
            <label class="text-caption text-weight-bold text-main block q-mb-xs">Display Name</label>
            <q-input outlined dense v-model="formData.name" placeholder="E.g. Jane Doe" />
          </div>

          <div>
            <label class="text-caption text-weight-bold text-main block q-mb-xs">About / Bio</label>
            <q-input outlined dense type="textarea" rows="2" v-model="formData.bio" placeholder="Write a short summary about yourself..." />
          </div>

          <div>
            <label class="text-caption text-weight-bold text-main block q-mb-xs">Instagram Handle/URL</label>
            <q-input
              outlined
              dense
              v-model="formData.instagram"
              placeholder="https://instagram.com/username"
            />
          </div>

          <div>
            <label class="text-caption text-weight-bold text-main block q-mb-xs">Facebook Profile URL</label>
            <q-input
              outlined
              dense
              v-model="formData.facebook"
              placeholder="https://facebook.com/username"
            />
          </div>

          <div>
            <label class="text-caption text-weight-bold text-main block q-mb-xs">LinkedIn Profile URL</label>
            <q-input
              outlined
              dense
              v-model="formData.linkedin"
              placeholder="https://linkedin.com/in/username"
            />
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-px-md q-pt-md">
          <q-btn flat label="Cancel" color="grey-6" v-close-popup no-caps class="text-weight-bold" />
          <q-btn
            unevaluated
            label="Save Settings"
            color="primary"
            no-caps
            class="q-px-lg text-weight-bold rounded-borders shadow-sm"
            @click="handleSave"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAppStore } from '@/stores/app';
import { user_service } from '@/boot/axios';
import axios from 'axios';
import Cookies from 'js-cookie';
import { Notify } from 'quasar';
import LoadingSpinner from '@/components/LoadingSpinner.vue';

const store = useAppStore();
const router = useRouter();

const localLoading = ref(false);
const editDialogOpen = ref(false);
const fileInput = ref<HTMLInputElement | null>(null);

const formData = reactive({
  name: '',
  bio: '',
  instagram: '',
  facebook: '',
  linkedin: '',
});

function syncFormData() {
  if (store.user) {
    formData.name = store.user.name || '';
    formData.bio = store.user.bio || '';
    formData.instagram = store.user.instagram || '';
    formData.facebook = store.user.facebook || '';
    formData.linkedin = store.user.linkedin || '';
  }
}

function triggerFileInput() {
  fileInput.value?.click();
}

async function handleImageUpload(e: Event) {
  const target = e.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file) return;

  const fd = new FormData();
  fd.append('file', file);

  localLoading.value = true;
  try {
    const token = Cookies.get('token');
    const { data } = await axios.post(`${user_service}/api/v1/user/update/pic`, fd, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    Notify.create({
      type: 'positive',
      message: data.message || 'Image uploaded successfully!',
      position: 'top',
    });

    Cookies.set('token', data.token);
    store.user = data.user;
  } catch (error) {
    Notify.create({
      type: 'negative',
      message: 'Image upload failed',
      position: 'top',
    });
    console.error(error);
  } finally {
    localLoading.value = false;
  }
}

function openEditDialog() {
  syncFormData();
  editDialogOpen.value = true;
}

async function handleSave() {
  localLoading.value = true;
  try {
    const token = Cookies.get('token');
    const { data } = await axios.post(`${user_service}/api/v1/user/update`, formData, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    Notify.create({
      type: 'positive',
      message: 'Profile updated successfully!',
      position: 'top',
    });

    Cookies.set('token', data.token);
    store.user = data.user;
    editDialogOpen.value = false;
  } catch (error) {
    Notify.create({
      type: 'negative',
      message: 'Update failed',
      position: 'top',
    });
    console.error(error);
  } finally {
    localLoading.value = false;
  }
}

onMounted(async () => {
  if (!store.isAuth && !store.loading) {
    void router.replace('/login');
    return;
  }
  if (store.loading) {
    await store.fetchUser();
  }
  syncFormData();
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
  transition: transform 0.2s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  background-color: var(--bg-card);
}

.avatar-wrapper:hover {
  transform: scale(1.03);
}

.avatar-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.45);
  opacity: 0;
  transition: opacity 0.25s ease;
  border-radius: 50%;
}

.avatar-wrapper:hover .avatar-overlay {
  opacity: 1;
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

.profile-settings-btn {
  border-color: var(--border-color) !important;
}

.profile-settings-btn:hover {
  background-color: rgba(0, 0, 0, 0.03) !important;
}

body.body--dark .profile-settings-btn:hover {
  background-color: rgba(255, 255, 255, 0.03) !important;
}

.edit-profile-modal {
  background-color: var(--bg-card) !important;
  border: 1px solid var(--border-color);
}

.max-width-md {
  max-width: 400px;
}

.initial-loading {
  z-index: 2;
}
</style>
