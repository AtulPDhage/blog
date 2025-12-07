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
export const dynamic = "force-dynamic";

const BlogPage = () => {
  const router = useRouter();
  const { isAuth, user, fetchBlogs, savedBlogs, getSavedBlogs } = useAppData();
  const { id } = useParams();
  const [blog, setBlog] = useState<Blog | null>(null);
  const [author, setAuthor] = useState<User | null>(null);
  const [loading, setLoading] = useState(false);
  const [addCommentLoading, setAddCommentLoading] = useState(false);

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
      setAddCommentLoading(true);
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
      setAddCommentLoading(false);
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
            <Button
              onClick={addComment}
              disabled={addCommentLoading}
              className={`transition-all duration-300 ease-in-out transform 
              ${
                comment.length === 0
                  ? "opacity-0 scale-95"
                  : "opacity-100 scale-100"
              }`}
            >
              {addCommentLoading ? "Adding Comment..." : "Post Comment"}
            </Button>
          </CardContent>
        </Card>
      )}
      <Card className="">
        <CardHeader className="px-6 pt-6">
          <h3 className="text-xl font-semibold text-gray-800">All Comments</h3>
        </CardHeader>

        <CardContent className="px-6 pb-6 flex flex-col gap-4">
          {comments && comments.length > 0 ? (
            comments.map((e, i) => (
              <div
                key={i}
                className="flex items-start justify-between gap-4 p-4 bg-white rounded-2xl shadow-sm border border-gray-200 
             transition-all duration-200 hover:shadow-md hover:-translate-y-1 hover:bg-gray-100"
              >
                <div className="flex flex-col gap-1 w-full">
                  <p className="font-semibold flex items-center gap-2 text-gray-700">
                    <span className="user border border-gray-300 rounded-full p-1 flex items-center justify-center">
                      <User2 className="w-4 h-4 text-gray-600" />
                    </span>
                    {e?.username}
                  </p>
                  <p className="wrap-word-break text-gray-800">{e.comment}</p>
                  <p className="text-xs text-gray-400">
                    {new Date(e.created_at).toLocaleString()}
                  </p>
                </div>

                {e.userid === user?._id && (
                  <Button
                    variant={"destructive"}
                    onClick={() => deleteComment(e.id)}
                    disabled={loading}
                    className="ml-2 h-8 w-8 p-0 flex items-center justify-center"
                  >
                    <Trash2 className="w-4 h-4 hover:bg-amber-300" />
                  </Button>
                )}
              </div>
            ))
          ) : (
            <p className="text-gray-500 text-center py-4">No Comments yet</p>
          )}
        </CardContent>
      </Card>
    </div>
  );
};

export default BlogPage;
