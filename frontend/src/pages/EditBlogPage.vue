<template>
  <q-page class="q-py-xl q-px-md bg-grey-1">
    <div v-if="initLoading" class="flex flex-center">
      <loading-spinner />
    </div>

    <div v-else class="new-blog-container mx-auto">
      <q-card flat borderless class="blog-card shadow-15 q-pa-lg">
        <q-card-section>
          <div class="text-h5 text-bold font-brand text-grey-9 q-mb-md">Edit Blog Post</div>
        </q-card-section>

        <q-card-section>
          <q-form @submit.prevent="handleSubmit" class="q-gutter-y-md">
            <!-- Title -->
            <div>
              <label class="text-subtitle2 text-grey-8 font-brand">Title</label>
              <q-input
                v-model="formData.title"
                outlined
                dense
                placeholder="Enter Blog title"
                required
                class="q-mt-xs full-width"
              />
            </div>

            <!-- Description -->
            <div>
              <label class="text-subtitle2 text-grey-8 font-brand">Description</label>
              <q-input
                v-model="formData.description"
                outlined
                dense
                placeholder="Enter Blog description"
                required
                class="q-mt-xs full-width"
              />
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
                required
                class="q-mt-xs"
              />
            </div>

            <!-- Cover Image Preview & Upload -->
            <div>
              <label class="text-subtitle2 text-grey-8 font-brand">Cover Image</label>
              <div v-if="existingImage && !imageFile" class="q-mt-xs q-mb-sm">
                <q-img
                  :src="existingImage"
                  class="cover-preview rounded-borders shadow-1"
                  fit="cover"
                />
                <div class="text-caption text-grey-5 q-mt-xs">Current cover image</div>
              </div>
              <q-file
                v-model="imageFile"
                outlined
                dense
                accept="image/*"
                placeholder="Click to replace cover image"
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
              <label class="text-subtitle2 text-grey-8 font-brand">Blog Content</label>
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
              label="Save Changes"
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

.cover-preview {
  max-width: 260px;
  max-height: 160px;
  border-radius: 8px;
}

.editor-wrapper {
  border: 1px solid rgba(0, 0, 0, 0.12);
  border-radius: 8px;
  overflow: hidden;
}
</style>
