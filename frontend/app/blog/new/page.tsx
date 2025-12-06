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
import { author_service, useAppData } from "@/context/AppContext";
import toast from "react-hot-toast";
import { useRouter } from "next/navigation";
const JoditEditor = dynamic(() => import("jodit-react"), { ssr: false });
export const blogCategories = [
  "Technology",
  "Health",
  "Entertainment",
  "Finance",
  "Travel",
  "Education",
  "Study",
];
const AddBlog = () => {
  const editor = useRef(null);
  const { fetchBlogs, isAuth } = useAppData();
  const [content, setContent] = useState("");
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState({
    title: "",
    description: "",
    category: "",
    image: "",
    blogcontent: "",
  });

  const router = useRouter();

  // Redirect if not authenticated
  useEffect(() => {
    if (!loading && !isAuth) {
      router.replace("/login");
    }
  }, [loading, isAuth]);

  const handleInputChange = (e: any) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleFileChnage = (e: any) => {
    const file = e.target.files[0];
    setFormData({ ...formData, image: file });
  };
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
        `${author_service}/api/v1/blog/new`,
        formDataToSend,
        {
          headers: { Authorization: `Bearer ${token}` },
        }
      );

      toast.success(data.message);

      setFormData({
        title: "",
        description: "",
        category: "",
        image: "",
        blogcontent: "",
      });
      setContent("");
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

  const [aiTitle, setAiTitle] = useState(false);

  const aiTitleResponse = async () => {
    try {
      setAiTitle(true);
      const { data }: any = await axios.post(
        `${author_service}/api/v1/ai/title`,
        {
          text: formData.title,
        }
      );
      setFormData({ ...formData, title: data });
    } catch (err) {
      toast.error("Problem while fetching ai");
      console.log(err);
    } finally {
      setAiTitle(false);
    }
  };

  const [aiDescription, setaiDescription] = useState(false);

  const aiDescriptionResponse = async () => {
    try {
      setaiDescription(true);
      const { data }: any = await axios.post(
        `${author_service}/api/v1/ai/description`,
        {
          title: formData.title,
          description: formData.description,
        }
      );
      setFormData({ ...formData, description: data });
    } catch (err) {
      toast.error("Problem while fetching ai");
      console.log(err);
    } finally {
      setaiDescription(false);
    }
  };
  const [aiBlogLoading, setAiBlogLoading] = useState(false);
  const aiBlogResponse = async () => {
    try {
      setAiBlogLoading(true);
      const { data }: any = await axios.post(
        `${author_service}/api/v1/ai/blog`,
        {
          blog: formData.blogcontent,
        }
      );
      setContent(data.html);
      setFormData({ ...formData, blogcontent: data.html });
    } catch (err) {
      toast.error("Problem while fetching ai");
      console.log(err);
    } finally {
      setAiBlogLoading(false);
    }
  };

  const config = useMemo(
    () => ({
      readonly: false, // all options from https://xdsoft.net/jodit/docs/,
      placeholder: "Start typings...",
    }),
    []
  );
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
                className={
                  aiTitle ? "animate-pulse placeholder:opacity-60" : ""
                }
                required
              />
              {formData.title === "" ? (
                ""
              ) : (
                <Button
                  type="button"
                  onClick={aiTitleResponse}
                  disabled={aiTitle}
                >
                  <RefreshCw className={aiTitle ? "animate-spin" : ""} />
                </Button>
              )}
            </div>
            <Label>Description</Label>
            <div className="flex justify-center items-center gap-2">
              <Input
                name="description"
                title={formData.description}
                value={formData.description}
                onChange={handleInputChange}
                placeholder="Enter Blog description"
                className={
                  aiDescription ? "animate-pulse placeholder:opacity-60" : ""
                }
                required
              />
              {formData.title === "" ? (
                ""
              ) : (
                <Button
                  type="button"
                  onClick={aiDescriptionResponse}
                  disabled={aiDescription}
                >
                  <RefreshCw className={aiDescription ? "animate-spin" : ""} />
                </Button>
              )}
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
              <Input type="file" accept="image/*" onChange={handleFileChnage} />
            </div>
            <div>
              <Label>Blog Content</Label>
              <div className="flex justify-between items-center mb-2">
                <p className="text-sm text-muted-foreground">
                  Paste your blog or type here. You can use rich text
                  formatting. Please add image after improving your grammer
                </p>
                <Button
                  type="button"
                  size={"sm"}
                  onClick={aiBlogResponse}
                  disabled={aiBlogLoading}
                >
                  <RefreshCw
                    size={16}
                    className={aiBlogLoading ? "animate-spin" : ""}
                  />
                  <span>Fix Grammer</span>
                </Button>
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

export default AddBlog;
