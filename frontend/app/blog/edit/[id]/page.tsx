"use client";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import React, { useEffect, useMemo, useRef, useState } from "react";
import { Car, RefreshCw } from "lucide-react";
import Cookies from "js-cookie";
import axios from "axios";
import {
  SelectContent,
  SelectTrigger,
  SelectValue,
  SelectItem,
  Select,
} from "@/components/ui/select";
import dynamic from "next/dynamic";
import {
  author_service,
  blog_service,
  useAppData,
  blogCategories,
  user_service,
} from "@/context/AppContext";
import toast from "react-hot-toast";
import { useParams, useRouter } from "next/navigation";
const JoditEditor = dynamic(() => import("jodit-react"), { ssr: false });

const EditBlogPage = () => {
  const router = useRouter();
  const { loading: contexLoader, isAuth, blogs, savedBlogs } = useAppData();

  // Redirect if not authenticated
  useEffect(() => {
    if (!contexLoader && !isAuth) {
      router.replace("/login");
    }
  }, [contexLoader, isAuth]);

  const editor = useRef(null);
  const { fetchBlogs } = useAppData();
  const { id } = useParams();
  const [content, setContent] = useState("");
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState({
    title: "",
    description: "",
    category: "",
    image: "",
    blogcontent: "",
  });

  const handleInputChange = (e: any) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleFileChnage = (e: any) => {
    const file = e.target.files[0];
    setFormData({ ...formData, image: file });
  };
  const config = useMemo(
    () => ({
      readonly: false, // all options from https://xdsoft.net/jodit/docs/,
      placeholder: "Start typings...",
    }),
    []
  );
  const [existingImage, setExistingImage] = useState(null);
  useEffect(() => {
    const fetchBlog = async () => {
      setLoading(true);
      try {
        const { data }: any = await axios.get(
          `${blog_service}/api/v1/blog/${id}`
        );
        const blog = data.blog.blog;
        setFormData({
          title: blog.title,
          description: blog.description,
          category: blog.category,
          image: "",
          blogcontent: blog.blogcontent,
        });

        setContent(blog.blogcontent);
        setExistingImage(blog.image);
      } catch (error) {
        console.log(error);
      } finally {
        setLoading(false);
      }
    };

    if (id) fetchBlog();
  }, [id]);

  const handleSubmit = async (e: any) => {
    e.preventDefault();
    setLoading(true);

    const formDataToSend = new FormData();
    formDataToSend.append("title", formData.title);
    formDataToSend.append("description", formData.description);
    formDataToSend.append("category", formData.category);
    formDataToSend.append("blogcontent", formData.blogcontent);

    if (formData.image) {
      formDataToSend.append("file", formData.image);
    }

    try {
      const token = Cookies.get("token");
      const { data }: any = await axios.post(
        `${author_service}/api/v1/blog/${id}`,
        formDataToSend,
        {
          headers: { Authorization: `Bearer ${token}` },
        }
      );

      toast.success(data.message);
      setTimeout(() => {
        fetchBlogs();
      }, 4000);
    } catch (error) {
      console.log(error);
      toast.error("Error while adding blog");
    } finally {
      setLoading(false);
    }
  };
  return (
    <div className="max-w-4xl mx-auto p-6">
      <Card>
        <CardHeader>
          <h2 className="text-2xl font-bold">Add New Blog</h2>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <Label>Title</Label>
            <div className="flex justify-center items-center gap-2">
              <Input
                name="title"
                value={formData.title}
                onChange={handleInputChange}
                placeholder="Enter Blog title"
                required
              />
            </div>
            <Label>Description</Label>
            <div className="flex justify-center items-center gap-2">
              <Input
                name="description"
                title={formData.description}
                value={formData.description}
                onChange={handleInputChange}
                placeholder="Enter Blog description"
                required
              />
            </div>
            <Label>Category</Label>
            <Select
              value={formData.category}
              onValueChange={(value: any) =>
                setFormData({ ...formData, category: value })
              }
            >
              <SelectTrigger>
                <SelectValue placeholder={"Select Category"}></SelectValue>
                <SelectContent>
                  {blogCategories?.map((e, i) => (
                    <SelectItem key={i} value={e}>
                      {e}
                    </SelectItem>
                  ))}
                </SelectContent>
              </SelectTrigger>
            </Select>
            <div>
              <Label>Upload Image</Label>
              {existingImage && !formData.image && (
                <img
                  src={existingImage}
                  className="w-40 h-40 object-cover rounded mb-2 mt-2"
                />
              )}
              <Input type="file" accept="image/*" onChange={handleFileChnage} />
            </div>
            <div>
              <Label>Blog Content</Label>
              <div className="flex justify-between items-center mb-2">
                <p className="text-sm text-muted-foreground">
                  Paste your blog or type here. You can use rich text
                  formatting. Please add image after improving your grammer
                </p>
              </div>
              <JoditEditor
                ref={editor}
                value={content}
                config={config}
                tabIndex={1}
                onBlur={(newContent) => {
                  setContent(newContent);
                  setFormData({ ...formData, blogcontent: newContent });
                }}
              />
            </div>
            <Button type="submit" className="w-full" disabled={loading}>
              {loading ? "Submitting" : "Submit"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
};

export default EditBlogPage;
