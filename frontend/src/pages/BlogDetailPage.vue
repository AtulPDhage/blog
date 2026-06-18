<template>
  <q-page class="q-py-xl q-px-md bg-app-container">
    <div v-if="loading || !blog" class="flex flex-center initial-loading-container">
      <loading-spinner />
    </div>

    <div v-else class="blog-detail-container mx-auto">
      <!-- Back Navigation -->
      <div class="q-mb-md">
        <q-btn
          flat
          dense
          no-caps
          color="primary"
          icon="arrow_back"
          label="Back to blogs"
          class="back-btn text-weight-bold"
          @click="router.push('/blogs')"
        />
      </div>

      <!-- Main Blog Details Card -->
      <q-card flat class="blog-detail-card q-pa-lg q-mb-lg shadow-sm">
        <q-card-section class="q-pa-none">
          <!-- Category Badge if available -->
          <div v-if="blog.category" class="q-mb-sm">
            <span class="detail-category-badge font-brand">{{ blog.category }}</span>
          </div>

          <!-- Blog Title -->
          <h1 class="text-h3 text-bold font-brand text-main q-mt-none q-mb-md leading-tight">
            {{ blog.title }}
          </h1>

          <!-- Author Info & Save/Edit Actions Header -->
          <div class="row items-center justify-between q-mt-lg q-pb-md border-bottom-light q-col-gutter-y-sm">
            <div class="row items-center q-gutter-x-sm">
              <q-avatar
                size="40px"
                class="shadow-sm cursor-pointer author-avatar-wrapper"
                @click="router.push(`/profile/${author?._id}`)"
              >
                <img :src="author?.image || '/default-avatar.png'" alt="author" />
              </q-avatar>
              
              <div class="column">
                <span
                  class="text-weight-bold text-main cursor-pointer hover-underline text-body2"
                  @click="router.push(`/profile/${author?._id}`)"
                >
                  {{ author?.name }}
                </span>
                <div class="row items-center text-caption text-sub q-gutter-x-xs font-brand">
                  <span>{{ formatPublishDate(blog.created_at) }}</span>
                  <span>•</span>
                  <span>{{ readingTime }} min read</span>
                </div>
              </div>

              <!-- Bookmark/Save Button -->
              <q-btn
                v-if="store.isAuth"
                flat
                round
                dense
                :color="saved ? 'primary' : 'grey-5'"
                :icon="saved ? 'bookmark' : 'bookmark_border'"
                class="q-ml-sm bookmark-btn"
                :loading="saveLoading"
                @click="toggleSaveBlog"
              >
                <q-tooltip class="bg-grey-9 text-white">{{ saved ? 'Saved to bookmarks' : 'Bookmark this blog' }}</q-tooltip>
              </q-btn>
            </div>

            <!-- Author Owner Actions (Edit & Delete) -->
            <div v-if="blog.author === store.user?._id" class="row items-center q-gutter-x-sm">
              <q-btn
                unevaluated
                color="primary"
                icon="edit"
                label="Edit"
                no-caps
                class="q-px-md rounded-borders text-weight-bold shadow-sm"
                :to="`/blog/edit/${blog.id}`"
              />
              <q-btn
                outlined
                color="negative"
                icon="delete"
                label="Delete"
                no-caps
                class="q-px-md rounded-borders text-weight-bold"
                @click="handleDeleteBlog"
              />
            </div>
          </div>
        </q-card-section>

        <!-- Cover Image -->
        <q-card-section class="q-px-none q-py-lg">
          <div class="cover-image-wrapper shadow-sm rounded-borders">
            <q-img
              :src="blog.image || '/placeholder-cover.jpg'"
              alt="cover"
              class="cover-img"
              fit="cover"
            />
          </div>
        </q-card-section>

        <!-- Blog Description -->
        <q-card-section class="q-px-none q-pt-none q-pb-md">
          <div class="blog-subtitle text-body1 text-weight-medium text-muted leading-relaxed italic q-pl-md border-left-accent">
            {{ blog.description }}
          </div>
        </q-card-section>

        <q-separator class="q-my-md" />

        <!-- Sanitized HTML Blog Content with rich typography -->
        <q-card-section class="q-px-none q-pb-none blog-content-section">
          <div v-html="sanitizedBlogContent" class="prose max-width-none" />
        </q-card-section>
      </q-card>

      <!-- Write a Comment Card -->
      <q-card v-if="store.isAuth" flat class="comment-composer-card q-pa-lg q-mb-lg shadow-sm">
        <q-card-section class="q-pa-none">
          <div class="text-subtitle1 text-bold font-brand text-main q-mb-md">Leave a comment</div>
          
          <div class="row no-wrap items-start q-gutter-x-md">
            <q-avatar size="36px" class="shadow-sm">
              <img :src="store.user?.image || '/default-avatar.png'" alt="me" />
            </q-avatar>
            
            <div class="column flex-grow-1">
              <q-input
                v-model="newComment"
                outlined
                type="textarea"
                rows="3"
                placeholder="Share your thoughts on this article..."
                class="comment-input q-mb-md"
              />
              
              <q-btn
                unevaluated
                color="primary"
                label="Post Comment"
                no-caps
                class="rounded-borders q-px-lg comment-btn align-self-end text-weight-bold"
                :loading="addCommentLoading"
                :disabled="!newComment.trim()"
                @click="handleAddComment"
              />
            </div>
          </div>
        </q-card-section>
      </q-card>

      <!-- Comments List Card -->
      <q-card flat class="comments-card q-pa-lg shadow-sm">
        <q-card-section class="q-pa-none">
          <div class="text-subtitle1 text-bold font-brand text-main q-mb-lg">
            Discussion ({{ comments.length }})
          </div>

          <div v-if="comments.length > 0" class="comments-list-wrapper">
            <div
              v-for="(commentItem, index) in comments"
              :key="commentItem.id"
              class="comment-row-item q-py-md"
              :class="{ 'border-top-light': index > 0 }"
            >
              <div class="row no-wrap items-start q-gutter-x-md">
                <q-avatar size="36px" class="bg-indigo-1 text-primary shadow-sm font-brand text-weight-bold">
                  {{ commentItem.username ? (commentItem.username.charAt(0).toUpperCase() || 'U') : 'U' }}
                </q-avatar>
                
                <div class="column flex-grow-1 overflow-hidden">
                  <div class="row items-center justify-between no-wrap">
                    <span class="text-weight-bold text-main text-body2 ellipsis">{{ commentItem.username }}</span>
                    <span class="text-caption text-sub font-brand">{{ formatCommentDate(commentItem.created_at) }}</span>
                  </div>
                  <p class="text-body2 text-muted q-mt-xs q-mb-none wrap-text leading-relaxed">
                    {{ commentItem.comment }}
                  </p>
                </div>

                <!-- Delete Comment Button -->
                <q-btn
                  v-if="commentItem.userid === store.user?._id"
                  flat
                  round
                  dense
                  color="negative"
                  icon="delete_outline"
                  size="sm"
                  class="q-ml-sm self-start delete-btn-hover"
                  @click="handleDeleteComment(commentItem.id)"
                >
                  <q-tooltip class="bg-red text-white">Delete Comment</q-tooltip>
                </q-btn>
              </div>
            </div>
          </div>
          <div v-else class="text-center text-sub q-py-xl column items-center">
            <q-icon name="forum" size="48px" class="q-mb-md text-sub" />
            <div class="text-body1 text-weight-medium">No comments yet</div>
            <div class="text-caption text-sub">Be the first to share your thoughts!</div>
          </div>
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
  return store.savedBlogs.some((b) => String(b.id) === blogId.value);
});

// Dynamic reading time calculator
const readingTime = computed(() => {
  if (!blog.value?.blogcontent) return 0;
  const words = blog.value.blogcontent.replace(/<[^>]*>/g, '').split(/\s+/).length;
  return Math.ceil(words / 200) || 1; // 200 words per minute
});

async function fetchBlogDetails() {
  loading.value = true;
  try {
    const { data } = await axios.get(`${blog_service}/api/v1/blog/${blogId.value}`);
    blog.value = data.blog;
    author.value = data.author?.user || null;
  } catch (error) {
    console.error('Failed to load blog:', error);
    Notify.create({
      type: 'negative',
      message: 'Failed to load blog post details',
      position: 'top',
    });
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
    title: 'Delete Comment',
    message: 'Are you sure you want to permanently delete this comment?',
    cancel: {
      flat: true,
      color: 'grey-6',
      label: 'Cancel'
    },
    ok: {
      unevaluated: true,
      color: 'negative',
      label: 'Delete'
    },
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
    title: 'Delete Post',
    message: 'Are you sure you want to delete this blog post? This action cannot be undone.',
    cancel: {
      flat: true,
      color: 'grey-6',
      label: 'Cancel'
    },
    ok: {
      unevaluated: true,
      color: 'negative',
      label: 'Delete Post'
    },
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

function formatPublishDate(time: string) {
  if (!time) return '';
  const date = new Date(time);
  return date.toLocaleDateString('en-US', {
    day: 'numeric',
    month: 'short',
    year: 'numeric'
  });
}

function formatCommentDate(time: string) {
  if (!time) return '';
  const date = new Date(time);
  return date.toLocaleDateString('en-US', {
    day: 'numeric',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit'
  });
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
.initial-loading-container {
  min-height: 70vh;
}

.blog-detail-container {
  max-width: 820px;
  margin: 0 auto;
}

.back-btn {
  border-radius: var(--radius-sm);
  color: var(--text-muted) !important;
}

.back-btn:hover {
  background-color: rgba(99, 102, 241, 0.05) !important;
  color: var(--q-primary) !important;
}

.blog-detail-card {
  border-radius: var(--radius-lg);
  background-color: var(--bg-card);
  border: 1px solid var(--border-color);
}

.detail-category-badge {
  background-color: rgba(99, 102, 241, 0.08);
  color: var(--q-primary);
  border: 1px solid rgba(99, 102, 241, 0.15);
  font-size: 0.75rem;
  font-weight: 600;
  padding: 4px 12px;
  border-radius: 9999px;
  text-transform: capitalize;
  display: inline-block;
}

.border-bottom-light {
  border-bottom: 1px solid var(--border-color);
}

.border-top-light {
  border-top: 1px solid var(--border-color);
}

.author-avatar-wrapper {
  transition: transform 0.2s ease;
  border: 1px solid var(--border-color);
}

.author-avatar-wrapper:hover {
  transform: scale(1.05);
}

.hover-underline:hover {
  text-decoration: underline;
  text-decoration-thickness: 1px;
}

.bookmark-btn {
  transition: transform 0.2s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

.bookmark-btn:active {
  transform: scale(0.85);
}

.cover-image-wrapper {
  overflow: hidden;
  max-height: 480px;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color);
}

.cover-img {
  max-height: 480px;
}

.blog-subtitle {
  font-size: 1.15rem;
}

.border-left-accent {
  border-left: 3px solid var(--q-primary);
}

/* Comments section styling */
.comment-composer-card,
.comments-card {
  border-radius: var(--radius-lg);
  background-color: var(--bg-card);
  border: 1px solid var(--border-color);
}

.comment-input :deep(.q-field__control) {
  border-radius: var(--radius-md) !important;
}

.comment-btn {
  align-self: flex-end;
}

.comment-row-item {
  transition: background-color 0.2s ease;
}

.delete-btn-hover {
  opacity: 0.6;
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.delete-btn-hover:hover {
  opacity: 1;
  transform: scale(1.08);
}

.wrap-text {
  word-break: break-word;
  white-space: pre-line;
}
</style>
