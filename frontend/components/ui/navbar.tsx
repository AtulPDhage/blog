"use client";

import { CircleUserRound, Ghost, LogIn, Menu, X } from "lucide-react";
import React, { useState } from "react";
import Link from "next/link";
import { Button } from "./button";
import { cn } from "@/lib/utils";
import { useAppData } from "@/context/AppContext";

const Navbar = () => {
  const [isOpen, setIsOpen] = useState(false);
  const { loading, isAuth } = useAppData();

  return (
    <nav className="bg-white shadow-md p-4 z-50">
      <div className="container mx-auto flex justify-between items-center">
        <Link href="/" className="text-xl font-bold text-gray-900">
          Postly
        </Link>

        {/* Mobile Menu Toggle */}
        <div className="md:hidden">
          <Button variant={"ghost"} onClick={() => setIsOpen(!isOpen)}>
            {isOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
          </Button>
        </div>
        <ul className="hidden md:flex justify-center items-center space-x-6 text-gray-700">
          <li>
            <Link href="/" className="hover:text-blue-500">
              Home
            </Link>
          </li>
          <li>
            <Link href="/blogs/saved" className="hover:text-blue-500">
              Saved Blogs
            </Link>
          </li>
          {loading ? (
            ""
          ) : (
            <li>
              {isAuth ? (
                <Link href="/profile" className="hover:text-blue-500">
                  <CircleUserRound />
                </Link>
              ) : (
                <Link href="/login" className="hover:text-blue-500">
                  <LogIn />
                </Link>
              )}
            </li>
          )}
        </ul>
      </div>
      {/* Mobile Menu */}
      <div
        className={cn(
          "md:hidden overflow-hidden transition-all duration-300 ease-in-out",
          isOpen ? "max-h-40 opacity-100" : "max-h-0 opacity-0"
        )}
      >
        <ul className="flex flex-col justify-center items-center space-y-4 p-4  text-gray-700 bg-white shadow-md">
          <li>
            <Link href="/" className="hover:text-blue-500">
              Home
            </Link>
          </li>
          <li>
            <Link href="/blogs/saved" className="hover:text-blue-500">
              Saved Blogs
            </Link>
          </li>
          <li>
            <Link href="/login" className="hover:text-blue-500">
              <LogIn />
            </Link>
          </li>
        </ul>
      </div>
    </nav>
  );
};

export default Navbar;
