<template>
  <q-page class="flex flex-center q-pa-md bg-app-container relative-position overflow-hidden">
    <!-- Glowing background elements for visual premium depth -->
    <div class="glow-orb glow-orb-primary"></div>
    <div class="glow-orb glow-orb-secondary"></div>

    <div v-if="store.loading || localLoading" class="initial-loading">
      <loading-spinner />
    </div>

    <div v-else class="login-container z-top">
      <q-card flat class="login-card shadow-xl q-pa-lg">
        <q-card-section class="text-center q-py-lg">
          <q-avatar size="44px" class="bg-indigo-1 text-primary q-mb-md brand-logo-icon shadow-sm">
            <q-icon name="auto_awesome" size="24px" class="brand-sparkle" />
          </q-avatar>
          <div class="text-h5 text-bold font-brand text-main q-mb-xs">Welcome to Postly</div>
          <div class="text-body2 text-sub">Sign in to start writing and reading stories</div>
        </q-card-section>

        <q-card-section class="q-px-md q-pb-lg">
          <q-btn
            unevaluated
            no-caps
            class="google-btn full-width q-py-md text-weight-bold"
            @click="triggerGoogleLogin"
          >
            <div class="row items-center justify-center q-gutter-x-md">
              <img src="/google.png" alt="Google icon" class="google-icon" />
              <span class="text-body2 text-weight-bold">Continue with Google</span>
            </div>
          </q-btn>
        </q-card-section>

        <q-card-section class="q-pa-none text-center q-pb-md">
          <p class="text-caption text-sub q-px-md">
            By signing in, you agree to our Terms of Service and Privacy Policy.
          </p>
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
  max-width: 400px;
}

.login-card {
  border-radius: var(--radius-lg);
  background-color: var(--bg-card);
  border: 1px solid var(--border-color);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
}

.brand-logo-icon {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.1), rgba(99, 102, 241, 0.2)) !important;
  border: 1px solid rgba(99, 102, 241, 0.15);
}

.brand-sparkle {
  background: linear-gradient(135deg, #6366f1, #a855f7);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.google-btn {
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md) !important;
  background-color: var(--bg-card) !important;
  box-shadow: var(--shadow-sm);
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  color: var(--text-main) !important;
}

.google-btn:hover {
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
  border-color: var(--q-primary);
  background-color: rgba(99, 102, 241, 0.02) !important;
}

.google-icon {
  width: 20px;
  height: 20px;
  display: block;
}

/* Background aesthetic glow orbs */
.glow-orb {
  position: absolute;
  width: 350px;
  height: 350px;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.12;
  z-index: 1;
}

.glow-orb-primary {
  background-color: #6366f1;
  top: -100px;
  left: -100px;
}

.glow-orb-secondary {
  background-color: #ec4899;
  bottom: -150px;
  right: -100px;
}

body.body--dark .glow-orb {
  opacity: 0.08;
}

.initial-loading {
  z-index: 2;
}
</style>
