<template>
  <q-page class="flex flex-center q-pa-md bg-grey-1">
    <div v-if="store.loading || localLoading">
      <loading-spinner />
    </div>

    <div v-else-if="store.user" class="profile-container full-width">
      <q-card flat borderless class="profile-card shadow-15 q-pa-lg">
        <q-card-section class="text-center">
          <div class="text-h5 text-bold font-brand text-grey-9 q-mb-md">Profile</div>
        </q-card-section>

        <q-card-section class="column items-center q-gutter-y-md">
          <!-- Clickable Avatar to Upload Pic -->
          <div class="avatar-wrapper relative-position cursor-pointer" @click="triggerFileInput">
            <q-avatar size="110px" class="profile-avatar shadow-4">
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

          <!-- User Name -->
          <div class="text-h6 text-weight-bold text-grey-9 text-center">
            {{ store.user.name }}
          </div>

          <!-- User Bio -->
          <div v-if="store.user.bio" class="text-body1 text-grey-7 text-center max-width-md">
            {{ store.user.bio }}
          </div>

          <!-- Social Links -->
          <div class="row q-gutter-x-md q-mt-sm">
            <a
              v-if="store.user.instagram"
              :href="store.user.instagram"
              target="_blank"
              class="social-link instagram"
            >
              <q-icon name="camera_alt" size="sm" />
            </a>
            <a
              v-if="store.user.facebook"
              :href="store.user.facebook"
              target="_blank"
              class="social-link facebook"
            >
              <q-icon name="facebook" size="sm" />
            </a>
            <a
              v-if="store.user.linkedin"
              :href="store.user.linkedin"
              target="_blank"
              class="social-link linkedin"
            >
              <q-icon name="work" size="sm" />
            </a>
          </div>

          <!-- Action Buttons -->
          <div class="row q-gutter-md q-mt-lg">
            <q-btn
              unevaluated
              color="primary"
              label="Add Blog"
              to="/blog/new"
              no-caps
              class="rounded-borders q-px-lg"
            />
            <q-btn
              outlined
              color="grey-7"
              label="Edit Profile"
              no-caps
              class="rounded-borders q-px-lg"
              @click="openEditDialog"
            />
          </div>
        </q-card-section>
      </q-card>
    </div>

    <!-- Edit Profile Dialog -->
    <q-dialog v-model="editDialogOpen" backdrop-filter="blur(4px)">
      <q-card style="width: 450px; max-width: 90vw" class="rounded-borders q-pa-md shadow-24">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6 text-bold font-brand text-grey-9">Edit Profile</div>
          <q-space />
          <q-btn flat round dense icon="close" v-close-popup />
        </q-card-section>

        <q-card-section class="q-gutter-y-sm">
          <div>
            <label class="text-caption text-weight-bold text-grey-7">Name</label>
            <q-input outlined dense v-model="formData.name" placeholder="Name" class="q-mt-xs" />
          </div>

          <div>
            <label class="text-caption text-weight-bold text-grey-7">Bio</label>
            <q-input outlined dense v-model="formData.bio" placeholder="Bio" class="q-mt-xs" />
          </div>

          <div>
            <label class="text-caption text-weight-bold text-grey-7">Instagram URL</label>
            <q-input
              outlined
              dense
              v-model="formData.instagram"
              placeholder="https://instagram.com/username"
              class="q-mt-xs"
            />
          </div>

          <div>
            <label class="text-caption text-weight-bold text-grey-7">Facebook URL</label>
            <q-input
              outlined
              dense
              v-model="formData.facebook"
              placeholder="https://facebook.com/username"
              class="q-mt-xs"
            />
          </div>

          <div>
            <label class="text-caption text-weight-bold text-grey-7">LinkedIn URL</label>
            <q-input
              outlined
              dense
              v-model="formData.linkedin"
              placeholder="https://linkedin.com/in/username"
              class="q-mt-xs"
            />
          </div>
        </q-card-section>

        <q-card-actions align="right" class="q-px-md">
          <q-btn flat label="Cancel" color="grey-6" v-close-popup no-caps />
          <q-btn
            unevaluated
            label="Save Changes"
            color="primary"
            no-caps
            class="q-px-md"
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
  // Wait if app is loading
  if (store.loading) {
    await store.fetchUser();
  }
  syncFormData();
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

.avatar-wrapper {
  position: relative;
  border-radius: 50%;
  overflow: hidden;
}

.profile-avatar {
  border: 4px solid #f3f4f6;
}

.avatar-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.4);
  opacity: 0;
  transition: opacity 0.2s ease;
  border-radius: 50%;
}

.avatar-wrapper:hover .avatar-overlay {
  opacity: 1;
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
