<template>
  <q-page class="q-py-xl q-px-md bg-grey-1">
    <div v-if="loading || !blog" class="flex flex-center">
      <loading-spinner />
    </div>

    <div v-else class="blog-detail-container mx-auto">
      <!-- Main Blog Card -->
      <q-card flat borderless class="blog-card shadow-15 q-pa-md q-mb-lg">
        <q-card-section>
          <!-- Blog Title -->
          <h1 class="text-h4 text-bold font-brand text-grey-9 q-mt-none q-mb-sm leading-tight">
            {{ blog.title }}
          </h1>

          <!-- Author Info & Save/Edit Actions -->
          <div class="row items-center justify-between q-mt-md q-col-gutter-y-sm">
            <div class="row items-center q-gutter-x-sm">
              <q-avatar
                size="32px"
                class="shadow-1 cursor-pointer"
                @click="router.push(`/profile/${author?._id}`)"
              >
                <img :src="author?.image || '/default-avatar.png'" alt="author" />
              </q-avatar>
              <span
                class="text-weight-bold text-grey-8 cursor-pointer hover-underline"
                @click="router.push(`/profile/${author?._id}`)"
              >
                {{ author?.name }}
              </span>

              <!-- Bookmark/Save Button -->
              <q-btn
                v-if="store.isAuth"
                flat
                round
                dense
                :color="saved ? 'positive' : 'grey-6'"
                :icon="saved ? 'bookmark' : 'bookmark_border'"
                class="q-ml-sm"
                :loading="saveLoading"
                @click="toggleSaveBlog"
              >
                <q-tooltip>{{ saved ? 'Saved' : 'Save Blog' }}</q-tooltip>
              </q-btn>
            </div>

            <!-- Author Owner Actions (Edit & Delete) -->
            <div v-if="blog.author === store.user?._id" class="row items-center q-gutter-x-sm">
              <q-btn
                unevaluated
                dense
                color="primary"
                icon="edit"
                label="Edit"
                no-caps
                class="q-px-sm rounded-borders"
                :to="`/blog/edit/${blog.id}`"
              />
              <q-btn
                outlined
                dense
                color="negative"
                icon="delete"
                label="Delete"
                no-caps
                class="q-px-sm rounded-borders"
                @click="handleDeleteBlog"
              />
            </div>
          </div>
        </q-card-section>

        <!-- Cover Image -->
        <q-card-section class="q-py-none">
          <q-img
            :src="blog.image"
            alt="cover"
            class="cover-img rounded-borders shadow-1 q-mb-md"
            style="max-height: 400px"
          />
        </q-card-section>

        <!-- Blog Description -->
        <q-card-section class="text-body1 text-grey-8 text-weight-medium q-py-sm">
          {{ blog.description }}
        </q-card-section>

        <q-separator class="q-my-md" />

        <!-- Sanitized HTML Blog Content -->
        <q-card-section class="blog-content text-grey-9">
          <div v-html="sanitizedBlogContent" class="prose max-width-none" />
        </q-card-section>
      </q-card>

      <!-- Write a Comment Card -->
      <q-card v-if="store.isAuth" flat borderless class="blog-card shadow-15 q-pa-md q-mb-lg">
        <q-card-section class="q-pb-none">
          <div class="text-subtitle1 text-bold font-brand text-grey-8">Leave a comment</div>
        </q-card-section>

        <q-card-section class="q-pt-sm">
          <q-input
            v-model="newComment"
            outlined
            dense
            type="textarea"
            rows="3"
            placeholder="Type your comment here..."
            class="q-mb-md comment-input"
          />
          <q-btn
            unevaluated
            color="primary"
            label="Post Comment"
            no-caps
            class="rounded-borders q-px-lg comment-btn"
            :loading="addCommentLoading"
            :disabled="!newComment.trim()"
            @click="handleAddComment"
          />
        </q-card-section>
      </q-card>

      <!-- Comments List Card -->
      <q-card flat borderless class="blog-card shadow-15 q-pa-md">
        <q-card-section class="q-pb-none">
          <div class="text-subtitle1 text-bold font-brand text-grey-8">All Comments</div>
        </q-card-section>

        <q-card-section>
          <div v-if="comments.length > 0" class="column q-gutter-y-sm">
            <div
              v-for="commentItem in comments"
              :key="commentItem.id"
              class="comment-card q-pa-md bg-grey-1 rounded-borders border-grey-3 row items-start no-wrap justify-between"
            >
              <div class="column q-gutter-y-xs full-width">
                <div class="row items-center q-gutter-x-xs text-weight-bold text-grey-8 text-body2">
                  <q-avatar size="20px" color="grey-3" class="text-grey-7">
                    <q-icon name="person" size="14px" />
                  </q-avatar>
                  <span>{{ commentItem.username }}</span>
                </div>
                <div class="text-body2 text-grey-9 wrap-text">
                  {{ commentItem.comment }}
                </div>
                <div class="text-caption text-grey-5">
                  {{ formatCommentDate(commentItem.created_at) }}
                </div>
              </div>

              <!-- Delete Comment Button -->
              <q-btn
                v-if="commentItem.userid === store.user?._id"
                flat
                round
                dense
                color="negative"
                icon="delete"
                size="sm"
                class="q-ml-sm"
                @click="handleDeleteComment(commentItem.id)"
              />
            </div>
          </div>
          <div v-else class="text-center text-grey-5 q-py-lg text-body1">No Comments yet</div>
        </q-card-section>
      </q-card>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useAppStore } from '../stores/app';
import type { Blog, User } from '../stores/app';
import { blog_service, author_service } from '@/boot/axios';
import { sanitizeHtml } from '@/boot/dompurify';
import axios from 'axios';
import Cookies from 'js-cookie';
import { Notify, Dialog } from 'quasar';
import LoadingSpinner from '@/components/LoadingSpinner.vue';

interface Comment {
  id: string;
  userid: string;
  comment: string;
  created_at: string;
  username: string;
}

const route = useRoute();
const router = useRouter();
const store = useAppStore();

const blog = ref<Blog | null>(null);
const author = ref<User | null>(null);
const loading = ref(true);

const comments = ref<Comment[]>([]);
const newComment = ref('');
const addCommentLoading = ref(false);
const saveLoading = ref(false);

const blogId = computed(() => route.params.id as string);

const sanitizedBlogContent = computed(() => {
  if (!blog.value?.blogcontent) return '';
  return sanitizeHtml(blog.value.blogcontent);
});

const saved = computed(() => {
  if (!store.savedBlogs) return false;
  return store.savedBlogs.some((b) => b.blogid === blogId.value);
});

async function fetchBlogDetails() {
  loading.value = true;
  try {
    const { data } = await axios.get(`${blog_service}/api/v1/blog/${blogId.value}`);
    blog.value = data.blog;
    author.value = data.author?.user || null;
  } catch (error) {
    console.error('Failed to load blog:', error);
  } finally {
    loading.value = false;
  }
}

async function fetchComments() {
  try {
    const { data } = await axios.get(`${blog_service}/api/v1/comment/${blogId.value}`);
    comments.value = data || [];
  } catch (error) {
    console.error('Failed to load comments:', error);
  }
}

async function handleAddComment() {
  if (!newComment.value.trim()) return;

  addCommentLoading.value = true;
  try {
    const token = Cookies.get('token');
    const { data } = await axios.post(
      `${blog_service}/api/v1/comment/${blogId.value}`,
      { comment: newComment.value },
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      },
    );
    Notify.create({
      type: 'positive',
      message: data.message || 'Comment added!',
      position: 'top',
    });
    newComment.value = '';
    await fetchComments();
  } catch (error) {
    Notify.create({
      type: 'negative',
      message: 'Problem while adding comment',
      position: 'top',
    });
    console.error(error);
  } finally {
    addCommentLoading.value = false;
  }
}

function handleDeleteComment(commentId: string) {
  Dialog.create({
    title: 'Confirm',
    message: 'Are you sure you want to delete this comment?',
    cancel: true,
    persistent: true,
  }).onOk(() => {
    const token = Cookies.get('token');
    axios
      .delete(`${blog_service}/api/v1/comment/${commentId}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      .then(({ data }) => {
        Notify.create({
          type: 'positive',
          message: data.message || 'Comment deleted',
          position: 'top',
        });
        void fetchComments();
      })
      .catch((error) => {
        Notify.create({
          type: 'negative',
          message: 'Problem while deleting comment',
          position: 'top',
        });
        console.error(error);
      });
  });
}

function handleDeleteBlog() {
  Dialog.create({
    title: 'Confirm Delete',
    message: 'Are you sure you want to delete this blog post?',
    cancel: true,
    persistent: true,
  }).onOk(() => {
    const token = Cookies.get('token');
    axios
      .delete(`${author_service}/api/v1/blog/${blogId.value}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      .then(({ data }) => {
        Notify.create({
          type: 'positive',
          message: data.message || 'Blog deleted successfully',
          position: 'top',
        });
        void router.push('/blogs');
        setTimeout(() => {
          void store.fetchBlogs();
        }, 4000);
      })
      .catch((error) => {
        Notify.create({
          type: 'negative',
          message: 'Problem while deleting blog',
          position: 'top',
        });
        console.error(error);
      });
  });
}

async function toggleSaveBlog() {
  saveLoading.value = true;
  try {
    const token = Cookies.get('token');
    const { data } = await axios.post(
      `${blog_service}/api/v1/save/${blogId.value}`,
      {},
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      },
    );
    Notify.create({
      type: 'positive',
      message: data.message || 'Bookmark updated',
      position: 'top',
    });
    await store.getSavedBlogs();
  } catch (error) {
    Notify.create({
      type: 'negative',
      message: 'Problem while bookmarking blog',
      position: 'top',
    });
    console.error(error);
  } finally {
    saveLoading.value = false;
  }
}

function formatCommentDate(time: string) {
  if (!time) return '';
  return new Date(time).toLocaleString();
}

onMounted(async () => {
  await fetchBlogDetails();
  await fetchComments();
  if (store.isAuth) {
    await store.getSavedBlogs();
  }
});
</script>

<style scoped>
.blog-detail-container {
  max-width: 800px;
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

.hover-underline:hover {
  text-decoration: underline;
}

.cover-img {
  border-radius: 12px;
  max-height: 420px;
}

.blog-content {
  font-size: 1.1rem;
  line-height: 1.75;
}

.comment-input :deep(.q-field__control) {
  border-radius: 8px;
}

.comment-card {
  border: 1px solid rgba(0, 0, 0, 0.06);
  border-radius: 12px;
  transition:
    box-shadow 0.2s,
    transform 0.2s;
}

.comment-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
  transform: translateY(-1px);
  background: #fdfdfd !important;
}

.wrap-text {
  word-break: break-word;
  white-space: pre-line;
}
</style>
