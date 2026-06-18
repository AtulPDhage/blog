<template>
  <q-page class="q-py-xl q-px-md bg-app-container">
    <div v-if="initLoading" class="flex flex-center initial-loading-container">
      <loading-spinner />
    </div>

    <div v-else class="new-blog-container mx-auto">
      <q-card flat class="blog-create-card q-pa-lg shadow-sm">
        <q-card-section class="q-pa-none q-mb-lg row items-center justify-between">
          <div>
            <h1 class="text-h4 text-bold font-brand text-main q-my-none">Edit Blog Post</h1>
            <p class="text-subtitle2 text-sub q-my-none">Update your story details and publish updates</p>
          </div>
          <span class="ai-badge">
            <q-icon name="auto_awesome" />
            <span>AI Assist Enabled</span>
          </span>
        </q-card-section>

        <q-card-section class="q-pa-none">
          <q-form @submit.prevent="handleSubmit" class="q-gutter-y-lg">
            <!-- Title -->
            <div>
              <div class="row justify-between items-center q-mb-xs">
                <label class="text-subtitle2 text-weight-bold text-main font-brand">Title</label>
                <q-btn
                  v-if="formData.title"
                  unevaluated
                  no-caps
                  class="ai-action-btn"
                  :loading="aiTitleLoading"
                  @click="aiImproveTitle"
                >
                  <div class="row items-center no-wrap q-gutter-x-xs">
                    <q-icon name="auto_awesome" size="14px" />
                    <span>AI Polish</span>
                  </div>
                </q-btn>
              </div>
              <q-input
                v-model="formData.title"
                outlined
                dense
                placeholder="Enter title"
                required
                class="title-input"
                :loading="aiTitleLoading"
              />
            </div>

            <!-- Description -->
            <div>
              <div class="row justify-between items-center q-mb-xs">
                <label class="text-subtitle2 text-weight-bold text-main font-brand">Short Description</label>
                <q-btn
                  v-if="formData.title"
                  unevaluated
                  no-caps
                  class="ai-action-btn"
                  :loading="aiDescLoading"
                  @click="aiImproveDescription"
                >
                  <div class="row items-center no-wrap q-gutter-x-xs">
                    <q-icon name="auto_awesome" size="14px" />
                    <span>AI Generate</span>
                  </div>
                </q-btn>
              </div>
              <q-input
                v-model="formData.description"
                outlined
                dense
                placeholder="Enter description"
                required
                class="desc-input"
                :loading="aiDescLoading"
              />
            </div>

            <!-- Category -->
            <div>
              <label class="text-subtitle2 text-weight-bold text-main font-brand q-mb-xs block">Category</label>
              <q-select
                v-model="formData.category"
                :options="categories"
                outlined
                dense
                placeholder="Select Category"
                required
                class="category-select"
              />
            </div>

            <!-- Cover Image Preview & Upload -->
            <div>
              <label class="text-subtitle2 text-weight-bold text-main font-brand q-mb-xs block">Cover Image</label>
              
              <div v-if="existingImage && !imageFile" class="q-mb-md existing-preview-container row items-end q-gutter-x-md">
                <div class="preview-img-wrapper shadow-sm rounded-borders">
                  <q-img
                    :src="existingImage"
                    class="cover-preview"
                    fit="cover"
                  />
                </div>
                <div class="column q-gutter-y-xs">
                  <span class="text-weight-medium text-body2 text-main">Current Cover Image</span>
                  <span class="text-caption text-sub">Will keep this unless you choose a new file below.</span>
                </div>
              </div>

              <q-file
                v-model="imageFile"
                outlined
                dense
                accept="image/*"
                placeholder="Click to upload new cover image..."
                class="image-file-input"
                @update:model-value="handleFileChange"
              >
                <template #prepend>
                  <q-icon name="cloud_upload" class="text-sub" />
                </template>
              </q-file>
            </div>

            <!-- Blog Content Jodit Editor -->
            <div>
              <div class="row items-center justify-between q-mb-xs">
                <label class="text-subtitle2 text-weight-bold text-main font-brand">Blog Story Content</label>
                <q-btn
                  unevaluated
                  no-caps
                  class="ai-action-btn"
                  :loading="aiBlogLoading"
                  @click="aiFixGrammar"
                >
                  <div class="row items-center no-wrap q-gutter-x-xs">
                    <q-icon name="auto_awesome" size="14px" />
                    <span>AI Fix Grammar</span>
                  </div>
                </q-btn>
              </div>
              <p class="text-caption text-sub q-mt-none q-mb-sm">
                Edit your story and format it with the dynamic editor toolbar as desired.
              </p>
              <div class="editor-wrapper">
                <textarea ref="editorRef"></textarea>
              </div>
            </div>

            <!-- Submit Button -->
            <q-btn
              type="submit"
              color="primary"
              label="Save Changes"
              no-caps
              unevaluated
              class="publish-btn full-width q-py-md text-weight-bold rounded-borders shadow-sm"
              :loading="submitLoading"
            />
          </q-form>
        </q-card-section>
      </q-card>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onBeforeUnmount, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useAppStore, blogCategories } from '@/stores/app';
import { blog_service, author_service } from '@/boot/axios';
import { Jodit } from 'jodit';
import 'jodit/es2021/jodit.min.css';
import axios from 'axios';
import Cookies from 'js-cookie';
import { Notify } from 'quasar';
import LoadingSpinner from '@/components/LoadingSpinner.vue';

const store = useAppStore();
const route = useRoute();
const router = useRouter();

const categories = blogCategories;
const blogId = computed(() => route.params.id as string);

const initLoading = ref(true);
const submitLoading = ref(false);
const aiTitleLoading = ref(false);
const aiDescLoading = ref(false);
const aiBlogLoading = ref(false);

const imageFile = ref<File | null>(null);
const existingImage = ref<string | null>(null);
const editorRef = ref<HTMLTextAreaElement | null>(null);
let joditInstance: Jodit | null = null;

const formData = reactive({
  title: '',
  description: '',
  category: '',
  image: null as File | null,
  blogcontent: '',
});

function handleFileChange(file: File | null) {
  formData.image = file;
}

async function aiImproveTitle() {
  aiTitleLoading.value = true;
  try {
    const { data } = await axios.post(`${author_service}/api/v1/ai/title`, {
      text: formData.title,
    });
    formData.title = data as string;
  } catch (error) {
    Notify.create({
      type: 'negative',
      message: 'Problem while generating title via AI',
      position: 'top',
    });
    console.error(error);
  } finally {
    aiTitleLoading.value = false;
  }
}

async function aiImproveDescription() {
  aiDescLoading.value = true;
  try {
    const { data } = await axios.post(`${author_service}/api/v1/ai/description`, {
      title: formData.title,
      description: formData.description,
    });
    formData.description = data as string;
  } catch (error) {
    Notify.create({
      type: 'negative',
      message: 'Problem while generating description via AI',
      position: 'top',
    });
    console.error(error);
  } finally {
    aiDescLoading.value = false;
  }
}

async function aiFixGrammar() {
  aiBlogLoading.value = true;
  try {
    const { data } = await axios.post(`${author_service}/api/v1/ai/blog`, {
      blog: formData.blogcontent,
    });
    formData.blogcontent = data.html as string;
    if (joditInstance) {
      joditInstance.value = data.html as string;
    }
  } catch (error) {
    Notify.create({
      type: 'negative',
      message: 'Problem while fixing grammar via AI',
      position: 'top',
    });
    console.error(error);
  } finally {
    aiBlogLoading.value = false;
  }
}

async function fetchBlogDetails() {
  initLoading.value = true;
  try {
    const { data } = await axios.get(`${blog_service}/api/v1/blog/${blogId.value}`);
    const blog = data.blog;
    formData.title = blog.title;
    formData.description = blog.description;
    formData.category = blog.category;
    formData.blogcontent = blog.blogcontent;
    existingImage.value = blog.image;
  } catch (error) {
    Notify.create({
      type: 'negative',
      message: 'Failed to fetch blog post details',
      position: 'top',
    });
    console.error(error);
  } finally {
    initLoading.value = false;
  }
}

async function handleSubmit() {
  if (!formData.blogcontent.trim()) {
    Notify.create({
      type: 'warning',
      message: 'Please write some content for the blog post.',
      position: 'top',
    });
    return;
  }

  submitLoading.value = true;
  const fd = new FormData();
  fd.append('title', formData.title);
  fd.append('description', formData.description);
  fd.append('category', formData.category);
  fd.append('blogcontent', formData.blogcontent);

  if (formData.image) {
    fd.append('file', formData.image);
  }

  try {
    const token = Cookies.get('token');
    const { data } = await axios.post(`${author_service}/api/v1/blog/${blogId.value}`, fd, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    Notify.create({
      type: 'positive',
      message: data.message || 'Blog updated successfully!',
      position: 'top',
    });

    void router.push(`/blog/${blogId.value}`);
    setTimeout(() => {
      void store.fetchBlogs();
    }, 4000);
  } catch (error) {
    Notify.create({
      type: 'negative',
      message: 'Error while saving changes',
      position: 'top',
    });
    console.error(error);
  } finally {
    submitLoading.value = false;
  }
}

onMounted(async () => {
  if (!store.isAuth && !store.loading) {
    void router.replace('/login');
    return;
  }

  await fetchBlogDetails();

  if (editorRef.value) {
    joditInstance = Jodit.make(editorRef.value, {
      readonly: false,
      placeholder: 'Start typing...',
      height: 380,
    });

    joditInstance.value = formData.blogcontent;
    joditInstance.events.on('change', (newValue: string) => {
      formData.blogcontent = newValue;
    });
  }
});

onBeforeUnmount(() => {
  if (joditInstance) {
    joditInstance.destruct();
  }
});
</script>

<style scoped>
.initial-loading-container {
  min-height: 70vh;
}

.new-blog-container {
  max-width: 820px;
  margin: 0 auto;
}

.blog-create-card {
  border-radius: var(--radius-lg);
  background-color: var(--bg-card);
  border: 1px solid var(--border-color);
}

.ai-action-btn {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.08), rgba(236, 72, 153, 0.08));
  border: 1px solid rgba(99, 102, 241, 0.15);
  color: var(--q-primary);
  border-radius: var(--radius-sm);
  font-size: 0.8rem;
  padding: 4px 10px;
  min-height: auto;
  font-weight: 600;
  transition: all 0.25s ease;
}

.ai-action-btn:hover {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.15), rgba(236, 72, 153, 0.15));
  border-color: var(--q-primary);
  box-shadow: 0 0 10px rgba(99, 102, 241, 0.1);
}

.preview-img-wrapper {
  max-width: 200px;
  max-height: 125px;
  overflow: hidden;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
}

.cover-preview {
  max-width: 200px;
  max-height: 125px;
}

.editor-wrapper {
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.publish-btn {
  font-size: 1rem;
}
</style>
