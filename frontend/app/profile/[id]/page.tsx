"use client";

import { Avatar, AvatarImage } from "@/components/ui/avatar";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Facebook, Instagram, Linkedin } from "lucide-react";

import React, { useState, useRef, useEffect } from "react";
import axios from "axios";
import Cookies from "js-cookie";
import toast from "react-hot-toast";
import { useParams, useRouter } from "next/navigation";
import Loading from "@/components/loading";
import { useAppData, User, user_service } from "@/context/AppContext";

const UserProfilePage = () => {
  const [user, setUser] = useState<User | null>(null);

  const { id } = useParams();
  async function fetchuser() {
    try {
      const { data }: any = await axios.get(
        `${user_service}/api/v1/user/${id}`
      );
      setUser(data.user);
    } catch (error) {
      console.log(error);
    }
  }

  useEffect(() => {
    fetchuser();
  }, [id]);
  if (!user) {
    return <Loading />;
  }
  return (
    <div className="flex items-center justify-center min-h-screen p-4">
      <Card className="w-full max-w-xl shadow-lg border rounded-2xl p-6">
        <CardHeader className="text-center">
          <CardTitle className="text-2xl font-semibold">Profile</CardTitle>
        </CardHeader>

        <CardContent className="flex flex-col items-center space-y-4">
          {/* Avatar */}
          <Avatar className="w-28 h-28 border-4 border-gray-200 shadow-md cursor-pointer">
            <AvatarImage src={user.image} alt="profile" />
          </Avatar>

          {/* Name */}
          <div className="text-center">
            <p className="font-medium text-lg">{user.name}</p>
          </div>

          {/* Bio */}
          {user.bio && (
            <div className="text-center">
              <p>{user.bio}</p>
            </div>
          )}

          {/* Social Links */}
          <div className="flex gap-4 mt-2">
            {user.instagram && (
              <a href={user.instagram} target="_blank">
                <Instagram className="text-pink-500" />
              </a>
            )}
            {user.facebook && (
              <a href={user.facebook} target="_blank">
                <Facebook className="text-blue-500" />
              </a>
            )}
            {user.linkedin && (
              <a href={user.linkedin} target="_blank">
                <Linkedin className="text-blue-700" />
              </a>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default UserProfilePage;
