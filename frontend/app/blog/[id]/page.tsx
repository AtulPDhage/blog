"use client";
import Loading from "@/components/loading";
import {
  author_service,
  Blog,
  blog_service,
  useAppData,
  User,
} from "@/context/AppContext";
import { useParams, useRouter } from "next/navigation";
import React, { useEffect, useState } from "react";
import axios from "axios";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import {
  Bookmark,
  BookmarkCheck,
  Edit,
  Trash2,
  User2,
  Users,
} from "lucide-react";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import Cookies from "js-cookie";
import toast from "react-hot-toast";

const BlogPage = () => {
  const router = useRouter();
  const { isAuth, user, fetchBlogs, savedBlogs, getSavedBlogs } = useAppData();
  const { id } = useParams();
  const [blog, setBlog] = useState<Blog | null>(null);
  const [author, setAuthor] = useState<User | null>(null);
  const [loading, setLoading] = useState(false);

  interface Comment {
    id: string;
    userid: string;
    comment: string;
    created_at: string;
    username: string;
  }
  async function fetchSingleBlog() {
    try {
      setLoading(true);
      const { data }: any = await axios.get(
        `${blog_service}/api/v1/blog/${id}`
      );
      setBlog(data.blog.blog);
      setAuthor(data.blog.author.user);
    } catch (error) {
      console.log(error);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    fetchSingleBlog();
  }, [id]);

  const deleteComment = async (id: string) => {
    if (confirm("Are you sure , you want to delete the comment")) {
      try {
        setLoading(true);
        const token = Cookies.get("token");
        const { data }: any = await axios.delete(
          `${blog_service}/api/v1/comment/${id}`,
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );
        toast.success(data.message);
        fetchComments();
      } catch (error) {
        toast.error("Problem while deleting comment");
      } finally {
        setLoading(false);
      }
    }
  };
  const [comments, setComments] = useState<Comment[]>([]);
  async function fetchComments() {
    try {
      setLoading(true);
      const { data }: any = await axios.get(
        `${blog_service}/api/v1/comment/${id}`
      );
      setComments(data);
    } catch (error) {
      toast.error("Problem while adding comment");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    fetchComments();
  }, [id]);

  const [comment, setComment] = useState("");
  async function addComment() {
    try {
      setLoading(true);
      const token = Cookies.get("token");
      const { data }: any = await axios.post(
        `${blog_service}/api/v1/comment/${id}`,
        { comment },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      toast.success(data.message);
      setComment("");
      fetchComments();
    } catch (error) {
      toast.error("Problem while adding comment");
    } finally {
      setLoading(false);
    }
  }
  async function deleteBlog() {
    if (confirm("Are you sure , you want to delete the blog")) {
      try {
        setLoading(true);
        const token = Cookies.get("token");
        const { data }: any = await axios.delete(
          `${author_service}/api/v1/blog/${id}`,
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );
        toast.success(data.message);
        router.push("/blogs");
        setTimeout(() => {
          fetchBlogs();
        }, 4000);
      } catch (error) {
        toast.error("Problem while deleting comment");
      } finally {
        setLoading(false);
      }
    }
  }

  const [saved, setSaved] = useState(false);
  useEffect(() => {
    if (savedBlogs && savedBlogs.some((b) => b.blogid === id)) {
      setSaved(true);
    } else {
      setSaved(false);
    }
  }, [savedBlogs, id]);

  async function saveBlog() {
    const token = Cookies.get("token");
    try {
      setLoading(true);
      const { data }: any = await axios.post(
        `${blog_service}/api/v1/save/${id}`,
        {},
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      toast.success(data.message);
      setSaved(!saved);
    } catch (error) {
      toast.error("Problem while saving blog");
    } finally {
      setLoading(false);
    }
  }
  if (!blog) {
    return <Loading />;
  }
  return (
    <div className="max-w-4xl mx-auto p-6 space-y-6">
      <Card>
        <CardHeader>
          <h1 className="text-3xl font-bold text-gray-900">{blog.title}</h1>
          <p className="text-gray-600 mt-2 flex items-center">
            <Link
              href={`/profile/${author?._id}`}
              className="flex items-center gap-2"
            >
              <img
                src={author?.image}
                className="w-8 h-8 rounded-full"
                alt=""
              />
              {author?.name}
            </Link>
            {isAuth && (
              <Button
                variant={"ghost"}
                className="mx-3"
                size={"lg"}
                disabled={loading}
                onClick={saveBlog}
              >
                {saved ? <BookmarkCheck color="#40d317" /> : <Bookmark />}
              </Button>
            )}
            {blog.author === user?._id && (
              <>
                <Button
                  size={"sm"}
                  onClick={() => router.push(`/blog/edit/${id}`)}
                >
                  <Edit />
                </Button>
                <Button
                  variant={"destructive"}
                  className="mx-2"
                  size={"sm"}
                  onClick={deleteBlog}
                  disabled={loading}
                >
                  <Trash2 />
                </Button>
              </>
            )}
          </p>
        </CardHeader>
        <CardContent>
          <img
            src={blog.image}
            alt=""
            className="w-full h-64 object-cover rounded-lg mb-4"
          />
          <p className="text-lg text-gray-700 mb-4">{blog.description}</p>
          <div
            className="prose max-w-none"
            dangerouslySetInnerHTML={{ __html: blog.blogcontent }}
          />
        </CardContent>
      </Card>

      {isAuth && (
        <Card>
          <CardHeader>
            <h3 className="text-xl font-semibold">Leave a comment</h3>
          </CardHeader>
          <CardContent>
            <Label htmlFor="comment">Your Comment</Label>
            <Input
              id="comment"
              placeholder="Type your comment here..."
              className="my-2"
              value={comment}
              onChange={(e) => setComment(e.target.value)}
            />
            <Button onClick={addComment} disabled={loading}>
              {loading ? "Adding Comment..." : "Post Comment"}
            </Button>
          </CardContent>
        </Card>
      )}
      <Card>
        <CardHeader>
          <h3 className="text-lg font-medium">All Comments</h3>
        </CardHeader>
        <CardContent>
          {comments && comments.length > 0 ? (
            comments.map((e, i) => {
              return (
                <div
                  key={i}
                  className="border-b mb-4 p-3  rounded-2xl border-black flex items-center "
                >
                  <div>
                    <p className="font-semibold flex items-center gap-1">
                      <span className="user border border-gray-400 rounded-full p-1">
                        <User2 />
                      </span>
                      {e?.username}
                    </p>
                    <p className="break-all">{e.comment}</p>
                    <p className="text-xs text-gray-500">
                      {new Date(e.created_at).toLocaleString()}
                    </p>
                  </div>
                  {e.userid === user?._id && (
                    <Button
                      variant={"destructive"}
                      onClick={() => deleteComment(e.id)}
                      disabled={loading}
                    >
                      <Trash2 />
                    </Button>
                  )}
                </div>
              );
            })
          ) : (
            <p>No Comments yet</p>
          )}
        </CardContent>
      </Card>
    </div>
  );
};

export default BlogPage;
