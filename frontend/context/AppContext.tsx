"use client";

import React, { ReactNode, useContext, useEffect, useState } from "react";
import { createContext } from "react";
import Cookies from "js-cookie";
import axios from "axios";
import toast, { Toaster } from "react-hot-toast";
import { GoogleOAuthProvider } from "@react-oauth/google";
import { usePathname } from "next/navigation";

export const user_service = "http://localhost:5000";
export const author_service = "http://localhost:5001";
export const blog_service = "http://localhost:5002";

export interface User {
  _id: string;
  name: string;
  email: string;
  image: string;
  instagram: string;
  linkedin: string;
  facebook: string;
  bio: string;
}

export interface Blog {
  id: string;
  title: string;
  description: string;
  blogcontent: string;
  image: string;
  category: string;
  author: string;
  created_at: string;
}

interface SavedBlogType {
  id: string;
  userid: string;
  blogid: string;
  created_at: string;
}
interface AppContextType {
  user: User | null;
  loading: boolean;
  isAuth: boolean;
  setLoading: React.Dispatch<React.SetStateAction<boolean>>;
  setIsAuth: React.Dispatch<React.SetStateAction<boolean>>;
  setUser: React.Dispatch<React.SetStateAction<User | null>>;
  logoutUser: () => Promise<void>;
  blogs: Blog[] | null;
  blogLoading: boolean;
  setBlogLoading: React.Dispatch<React.SetStateAction<boolean>>;
  setSearchQuery: React.Dispatch<React.SetStateAction<string>>;
  searchQuery: string;
  setCategory: React.Dispatch<React.SetStateAction<string>>;
  category: string;
  fetchBlogs: () => Promise<void>;
  savedBlogs: SavedBlogType[] | null;
  getSavedBlogs: () => Promise<void>;
}

const AppContext = createContext<AppContextType | undefined>(undefined);

interface AppProviderProps {
  children: ReactNode;
}

export const AppProvider: React.FC<AppProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isAuth, setIsAuth] = useState(false);
  const [loading, setLoading] = useState(true);
  const [category, setCategory] = useState("");
  const [searchQuery, setSearchQuery] = useState("");

  async function fetchUser() {
    setLoading(true);
    try {
      const token = Cookies.get("token");

      const { data }: any = await axios.get(`${user_service}/api/v1/me`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      setUser(data.user);
      setIsAuth(true);
      setLoading(false);
    } catch (err) {
      console.log(err);
      setLoading(false);
    }
  }
  const [blogLoading, setBlogLoading] = useState(false);
  const [blogs, setBlogs] = useState<Blog[] | null>(null);
  async function fetchBlogs() {
    setBlogLoading(true);
    try {
      const { data }: any = await axios.get(
        `${blog_service}/api/v1/blog/all?searchQuery=${searchQuery}&category=${category}`
      );
      setBlogs(data?.blogs);
    } catch (err) {
      console.log(err);
    } finally {
      setBlogLoading(false);
    }
  }

  const [savedBlogs, setSavedBlogs] = useState<SavedBlogType[] | null>(null);
  async function getSavedBlogs() {
    const token = Cookies.get("token");
    try {
      const { data }: any = await axios.get(
        `${blog_service}/api/v1/blogs/saved/all`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      setSavedBlogs(data);
    } catch (error) {
      console.log(error);
    }
  }

  const pathname = usePathname();
  useEffect(() => {
    if (pathname === "/blog/saved") {
      getSavedBlogs();
    }
  }, [pathname]);

  async function logoutUser() {
    Cookies.remove("token");
    setUser(null);
    setIsAuth(false);
    toast.success("user logged out");
  }

  useEffect(() => {
    const token = Cookies.get("token");

    if (token) {
      fetchUser().finally(() => setLoading(false));
    } else {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchBlogs();
  }, [searchQuery, category]);

  return (
    <AppContext.Provider
      value={{
        user,
        isAuth,
        loading,
        setIsAuth,
        setLoading,
        setUser,
        logoutUser,
        setBlogLoading,
        blogLoading,
        blogs,
        setSearchQuery,
        searchQuery,
        setCategory,
        category,
        fetchBlogs,
        savedBlogs,
        getSavedBlogs,
      }}
    >
      <GoogleOAuthProvider clientId="720531130636-e9r983didu9ske1smpbol7l6egj1e4im.apps.googleusercontent.com">
        {children}
        <Toaster />
      </GoogleOAuthProvider>
    </AppContext.Provider>
  );
};

export const useAppData = (): AppContextType => {
  const context = useContext(AppContext);
  if (!context) {
    throw new Error("useAppData must be used within an AppProvider");
  }
  return context;
};
