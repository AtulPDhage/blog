"use client";
import BlogCard from "@/components/BlogCard";
import Loading from "@/components/loading";
import { useAppData } from "@/context/AppContext";
import { useRouter } from "next/navigation";
import React, { useEffect } from "react";

const SavedBlogs = () => {
  const { loading, isAuth, blogs, savedBlogs } = useAppData();

  const router = useRouter();

  // Redirect if not authenticated
  useEffect(() => {
    if (!loading && !isAuth) {
      router.replace("/login");
    }
  }, [loading, isAuth]);
  if (!blogs || !savedBlogs) {
    return <Loading />;
  }

  const filteredBlogs = blogs.filter((blog) =>
    savedBlogs.some((saved) => saved.blogid === blog.id.toString())
  );

  return (
    <div className="container mx-auto px-4">
      <h1 className="text-3xl font-bold mt-2 ">Saved Blogs</h1>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3">
        {filteredBlogs.length > 0 ? (
          filteredBlogs.map((e, i) => {
            return (
              <BlogCard
                key={i}
                image={e.image}
                title={e.title}
                desc={e.description}
                id={e.id}
                time={e.created_at}
              />
            );
          })
        ) : (
          <p> No saved blogs yet!</p>
        )}
      </div>
    </div>
  );
};

export default SavedBlogs;
