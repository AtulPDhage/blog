<template>
  <q-page class="flex flex-center q-pa-md bg-grey-1">
    <div v-if="store.loading || localLoading">
      <loading-spinner />
    </div>

    <div v-else class="login-container">
      <q-card flat borderless class="login-card shadow-15">
        <q-card-section class="text-center q-py-lg">
          <div class="text-h5 text-bold font-brand text-grey-9 q-mb-xs">Login to Postly</div>
          <div class="text-subtitle2 text-grey-6">your go-to blogging platform</div>
        </q-card-section>

        <q-card-section class="q-px-lg q-pb-xl text-center">
          <q-btn
            unevaluated
            color="white"
            text-color="grey-9"
            class="google-btn full-width q-py-md text-weight-bold"
            no-caps
            @click="triggerGoogleLogin"
          >
            <div class="row items-center justify-center q-gutter-x-sm">
              <img src="/google.png" alt="Google icon" class="google-icon" />
              <span>Login with Google</span>
            </div>
          </q-btn>
        </q-card-section>
      </q-card>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAppStore } from '@/stores/app';
import { user_service } from '@/boot/axios';
import axios from 'axios';
import Cookies from 'js-cookie';
import { Notify } from 'quasar';
import LoadingSpinner from '@/components/LoadingSpinner.vue';

interface GoogleAuthClient {
  requestCode: () => void;
}

interface InitCodeClientConfig {
  client_id: string;
  scope: string;
  ux_mode: 'popup' | 'redirect';
  callback: (response: { code?: string }) => void;
}

declare const google: {
  accounts: {
    oauth2: {
      initCodeClient: (config: InitCodeClientConfig) => GoogleAuthClient;
    };
  };
};

const store = useAppStore();
const router = useRouter();
const localLoading = ref(false);

let googleClient: GoogleAuthClient | null = null;

function initializeGoogleClient() {
  if (typeof google === 'undefined') {
    // Retry loading in case the script tag hasn't loaded fully
    setTimeout(initializeGoogleClient, 200);
    return;
  }

  try {
    googleClient = google.accounts.oauth2.initCodeClient({
      client_id: import.meta.env.VITE_GOOGLE_CLIENT_ID,
      scope: 'openid email profile',
      ux_mode: 'popup',
      callback: (authResult) => {
        if (authResult && authResult.code) {
          void handleGoogleLogin(authResult.code);
        } else {
          Notify.create({
            type: 'negative',
            message: 'Google auth authorization code missing.',
            position: 'top',
          });
        }
      },
    });
  } catch (error) {
    console.error('Failed to initialize Google Auth Client', error);
  }
}

async function handleGoogleLogin(code: string) {
  localLoading.value = true;
  try {
    const { data } = await axios.post(`${user_service}/api/v1/login`, {
      code,
    });

    Cookies.set('token', data.token, {
      expires: 5,
      secure: true,
      path: '/',
    });

    Notify.create({
      type: 'positive',
      message: (data.message as string) || 'Login successful!',
      position: 'top',
    });

    store.isAuth = true;
    store.user = data.user;
    void router.push('/blogs');
  } catch (error) {
    Notify.create({
      type: 'negative',
      message: 'Login failed. Please try again.',
      position: 'top',
    });
    console.error('Google token exchange error:', error);
  } finally {
    localLoading.value = false;
  }
}

function triggerGoogleLogin() {
  if (!googleClient) {
    Notify.create({
      type: 'warning',
      message: 'Google login service is still loading, please try again in a moment.',
      position: 'top',
    });
    initializeGoogleClient();
    return;
  }
  googleClient.requestCode();
}

onMounted(() => {
  if (store.isAuth) {
    void router.replace('/blogs');
  } else {
    initializeGoogleClient();
  }
});
</script>

<style scoped>
.login-container {
  width: 100%;
  max-width: 380px;
}

.login-card {
  border-radius: 16px;
  background: white;
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.font-brand {
  font-family: 'Outfit', 'Inter', sans-serif;
  letter-spacing: -0.5px;
}

.google-btn {
  border: 1px solid #dadce0;
  border-radius: 8px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  transition:
    background-color 0.2s,
    border-color 0.2s;
}

.google-btn:hover {
  background-color: #f8f9fa !important;
  border-color: #c4c7c5;
}

.google-icon {
  width: 22px;
  height: 22px;
  display: block;
}
</style>
