<template>
  <q-page class="q-py-xl q-px-md bg-grey-1">
    <div class="new-blog-container mx-auto">
      <q-card flat borderless class="blog-card shadow-15 q-pa-lg">
        <q-card-section>
          <div class="text-h5 text-bold font-brand text-grey-9 q-mb-md">Add New Blog</div>
        </q-card-section>

        <q-card-section>
          <q-form @submit.prevent="handleSubmit" class="q-gutter-y-md">
            <!-- Title -->
            <div>
              <label class="text-subtitle2 text-grey-8 font-brand">Title</label>
              <div class="row items-center q-gutter-x-sm no-wrap q-mt-xs">
                <q-input
                  v-model="formData.title"
                  outlined
                  dense
                  placeholder="Enter Blog title"
                  :rules="[(val) => !!val || 'Title is required']"
                  lazy-rules
                  class="full-width title-input"
                  :loading="aiTitleLoading"
                />
                <q-btn
                  v-if="formData.title"
                  round
                  color="primary"
                  icon="auto_awesome"
                  :loading="aiTitleLoading"
                  @click="aiImproveTitle"
                >
                  <q-tooltip>AI Improve Title</q-tooltip>
                </q-btn>
              </div>
            </div>

            <!-- Description -->
            <div>
              <label class="text-subtitle2 text-grey-8 font-brand">Description</label>
              <div class="row items-center q-gutter-x-sm no-wrap q-mt-xs">
                <q-input
                  v-model="formData.description"
                  outlined
                  dense
                  placeholder="Enter Blog description"
                  :rules="[(val) => !!val || 'Description is required']"
                  lazy-rules
                  class="full-width desc-input"
                  :loading="aiDescLoading"
                />
                <q-btn
                  v-if="formData.title"
                  round
                  color="primary"
                  icon="auto_awesome"
                  :loading="aiDescLoading"
                  @click="aiImproveDescription"
                >
                  <q-tooltip>AI Improve Description</q-tooltip>
                </q-btn>
              </div>
            </div>

            <!-- Category -->
            <div>
              <label class="text-subtitle2 text-grey-8 font-brand">Category</label>
              <q-select
                v-model="formData.category"
                :options="categories"
                outlined
                dense
                placeholder="Select Category"
                :rules="[(val) => !!val || 'Category is required']"
                lazy-rules
                class="q-mt-xs"
              />
            </div>

            <!-- Image Upload -->
            <div>
              <label class="text-subtitle2 text-grey-8 font-brand">Upload Image</label>
              <q-file
                v-model="imageFile"
                outlined
                dense
                accept="image/*"
                placeholder="Click to upload cover image"
                :rules="[(val) => !!val || 'Cover image is required']"
                lazy-rules
                class="q-mt-xs"
                @update:model-value="handleFileChange"
              >
                <template #prepend>
                  <q-icon name="cloud_upload" />
                </template>
              </q-file>
            </div>

            <!-- Blog Content Jodit Editor -->
            <div>
              <div class="row items-center justify-between q-mb-xs">
                <label class="text-subtitle2 text-grey-8 font-brand">Blog Content</label>
                <q-btn
                  flat
                  dense
                  color="primary"
                  icon="auto_awesome"
                  label="Fix Grammar"
                  no-caps
                  :loading="aiBlogLoading"
                  @click="aiFixGrammar"
                />
              </div>
              <div class="text-caption text-grey-6 q-mb-sm">
                Paste your blog content or type here. Use the rich text formatting toolbar as
                desired.
              </div>
              <div class="editor-wrapper border-grey-4">
                <textarea ref="editorRef"></textarea>
              </div>
            </div>

            <!-- Submit Button -->
            <q-btn
              type="submit"
              color="primary"
              label="Submit Post"
              no-caps
              class="full-width q-py-md text-weight-bold rounded-borders q-mt-lg"
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
      placeholder: 'Start typings...',
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
  max-width: 800px;
  margin: 0 auto;
}

.blog-card {
  border-radius: 16px;
  background: white;
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.font-brand {
  font-family: 'Outfit', 'Inter', sans-serif;
  letter-spacing: -0.5px;
}

.editor-wrapper {
  border: 1px solid rgba(0, 0, 0, 0.12);
  border-radius: 8px;
  overflow: hidden;
}
</style>
