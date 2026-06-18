<template>
  <q-page class="q-py-xl q-px-md bg-app-container">
    <div class="new-blog-container mx-auto">
      <q-card flat class="blog-create-card q-pa-lg shadow-sm">
        <q-card-section class="q-pa-none q-mb-lg row items-center justify-between">
          <div>
            <h1 class="text-h4 text-bold font-brand text-main q-my-none">Create a New Post</h1>
            <p class="text-subtitle2 text-sub q-my-none">Write and publish your story on Postly</p>
          </div>
          <span class="ai-badge">
            <q-icon name="auto_awesome" />
            <span>AI Powered</span>
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
                  <q-tooltip class="bg-grey-9 text-white">AI Improve Title</q-tooltip>
                </q-btn>
              </div>
              <q-input
                v-model="formData.title"
                outlined
                dense
                placeholder="Enter a catchy title..."
                :rules="[(val) => !!val || 'Title is required']"
                lazy-rules
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
                  <q-tooltip class="bg-grey-9 text-white">AI Improve/Generate Description</q-tooltip>
                </q-btn>
              </div>
              <q-input
                v-model="formData.description"
                outlined
                dense
                placeholder="Write a brief overview of your blog post..."
                :rules="[(val) => !!val || 'Description is required']"
                lazy-rules
                class="desc-input"
                :loading="aiDescLoading"
              />
            </div>

            <!-- Category & Image Upload Section -->
            <div class="row q-col-gutter-md">
              <div class="col-12 col-sm-6">
                <label class="text-subtitle2 text-weight-bold text-main font-brand q-mb-xs block">Category</label>
                <q-select
                  v-model="formData.category"
                  :options="categories"
                  outlined
                  dense
                  placeholder="Select Category"
                  :rules="[(val) => !!val || 'Category is required']"
                  lazy-rules
                  class="category-select"
                />
              </div>

              <div class="col-12 col-sm-6">
                <label class="text-subtitle2 text-weight-bold text-main font-brand q-mb-xs block">Cover Image</label>
                <q-file
                  v-model="imageFile"
                  outlined
                  dense
                  accept="image/*"
                  placeholder="Choose cover image"
                  :rules="[(val) => !!val || 'Cover image is required']"
                  lazy-rules
                  class="image-file-input"
                  @update:model-value="handleFileChange"
                >
                  <template #prepend>
                    <q-icon name="cloud_upload" class="text-sub" />
                  </template>
                </q-file>
              </div>
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
                Write your thoughts. Use the formatting toolbar for customized rich formatting.
              </p>
              <div class="editor-wrapper">
                <textarea ref="editorRef"></textarea>
              </div>
            </div>

            <!-- Submit Button -->
            <q-btn
              type="submit"
              color="primary"
              label="Publish Post"
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
import { ref, reactive, onMounted, onBeforeUnmount } from 'vue';
import { useRouter } from 'vue-router';
import { useAppStore, blogCategories } from '@/stores/app';
import { author_service } from '@/boot/axios';
import { Jodit } from 'jodit';
import 'jodit/es2021/jodit.min.css';
import axios from 'axios';
import Cookies from 'js-cookie';
import { Notify } from 'quasar';

const store = useAppStore();
const router = useRouter();

const categories = blogCategories;

const submitLoading = ref(false);
const aiTitleLoading = ref(false);
const aiDescLoading = ref(false);
const aiBlogLoading = ref(false);

const imageFile = ref<File | null>(null);
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
    const { data } = await axios.post(`${author_service}/api/v1/blog/new`, fd, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    Notify.create({
      type: 'positive',
      message: data.message || 'Blog added successfully!',
      position: 'top',
    });

    // Reset Form
    formData.title = '';
    formData.description = '';
    formData.category = '';
    formData.image = null;
    formData.blogcontent = '';
    imageFile.value = null;
    if (joditInstance) {
      joditInstance.value = '';
    }

    void router.push('/blogs');
    setTimeout(() => {
      void store.fetchBlogs();
    }, 4000);
  } catch (error) {
    Notify.create({
      type: 'negative',
      message: 'Error while adding blog post',
      position: 'top',
    });
    console.error(error);
  } finally {
    submitLoading.value = false;
  }
}

onMounted(() => {
  if (!store.isAuth && !store.loading) {
    void router.replace('/login');
    return;
  }

  if (editorRef.value) {
    joditInstance = Jodit.make(editorRef.value, {
      readonly: false,
      placeholder: 'Start typing your story...',
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

.editor-wrapper {
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.publish-btn {
  font-size: 1rem;
}
</style>
