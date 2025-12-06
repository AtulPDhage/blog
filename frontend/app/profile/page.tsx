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
import { useRouter } from "next/navigation";
import Loading from "@/components/loading";

import { useAppData, user_service } from "@/context/AppContext";

const ProfilePage = () => {
  const { user, setUser, isAuth, loading } = useAppData();
  const router = useRouter();

  // Redirect if not authenticated
  useEffect(() => {
    if (!loading && !isAuth) {
      router.replace("/login");
    }
  }, [loading, isAuth]);

  const [formData, setFormData] = useState({
    name: "",
    bio: "",
    instagram: "",
    facebook: "",
    linkedin: "",
  });

  const [open, setOpen] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);
  const [localLoading, setLocalLoading] = useState(false);

  // 🟢 Sync formData automatically once the user loads
  useEffect(() => {
    if (user) {
      setFormData({
        name: user.name || "",
        bio: user.bio || "",
        instagram: user.instagram || "",
        facebook: user.facebook || "",
        linkedin: user.linkedin || "",
      });
    }
  }, [user]);

  const handleImageUpload = async (e: any) => {
    const file = e.target.files[0];

    if (!file) return;

    const fd = new FormData();
    fd.append("file", file);

    try {
      setLocalLoading(true);

      const token = Cookies.get("token");
      const { data }: any = await axios.post(
        `${user_service}/api/v1/user/update/pic`,
        fd,
        {
          headers: { Authorization: `Bearer ${token}` },
        }
      );

      toast.success(data.message);
      Cookies.set("token", data.token);
      setUser(data.user);
    } catch (err) {
      toast.error("Image upload failed");
    } finally {
      setLocalLoading(false);
    }
  };

  const handleSave = async () => {
    try {
      setLocalLoading(true);

      const token = Cookies.get("token");
      const { data }: any = await axios.post(
        `${user_service}/api/v1/user/update`,
        formData,
        {
          headers: { Authorization: `Bearer ${token}` },
        }
      );

      toast.success("Profile updated");
      Cookies.set("token", data.token);
      setUser(data.user);
      setOpen(false);
    } catch (err) {
      toast.error("Update failed");
    } finally {
      setLocalLoading(false);
    }
  };

  if (loading || localLoading || !user) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <Loading />
      </div>
    );
  }

  return (
    <div className="flex items-center justify-center min-h-screen p-4">
      <Card className="w-full max-w-xl shadow-lg border rounded-2xl p-6">
        <CardHeader className="text-center">
          <CardTitle className="text-2xl font-semibold">Profile</CardTitle>
        </CardHeader>

        <CardContent className="flex flex-col items-center space-y-4">
          {/* Avatar */}
          <Avatar
            className="w-28 h-28 border-4 border-gray-200 shadow-md cursor-pointer"
            onClick={() => inputRef.current?.click()}
          >
            <AvatarImage src={user.image} alt="profile" />
            <input
              type="file"
              accept="image/*"
              className="hidden"
              ref={inputRef}
              onChange={handleImageUpload}
            />
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

          {/* Buttons */}
          <div className="flex flex-col sm:flex-row gap-2 mt-6">
            <Button onClick={() => router.push("/blog/new")}>Add Blog</Button>

            {/* Edit Modal */}
            <Dialog open={open} onOpenChange={setOpen}>
              <DialogTrigger asChild>
                <Button variant="outline">Edit</Button>
              </DialogTrigger>

              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Edit Profile</DialogTitle>
                </DialogHeader>

                <div className="space-y-3 mt-2">
                  <div>
                    <Label>Name</Label>
                    <Input
                      value={formData.name}
                      onChange={(e) =>
                        setFormData({ ...formData, name: e.target.value })
                      }
                    />
                  </div>

                  <div>
                    <Label>Bio</Label>
                    <Input
                      value={formData.bio}
                      onChange={(e) =>
                        setFormData({ ...formData, bio: e.target.value })
                      }
                    />
                  </div>

                  <div>
                    <Label>Instagram</Label>
                    <Input
                      value={formData.instagram}
                      onChange={(e) =>
                        setFormData({ ...formData, instagram: e.target.value })
                      }
                    />
                  </div>

                  <div>
                    <Label>Facebook</Label>
                    <Input
                      value={formData.facebook}
                      onChange={(e) =>
                        setFormData({ ...formData, facebook: e.target.value })
                      }
                    />
                  </div>

                  <div>
                    <Label>LinkedIn</Label>
                    <Input
                      value={formData.linkedin}
                      onChange={(e) =>
                        setFormData({ ...formData, linkedin: e.target.value })
                      }
                    />
                  </div>

                  <Button className="w-full mt-4" onClick={handleSave}>
                    Save Changes
                  </Button>
                </div>
              </DialogContent>
            </Dialog>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default ProfilePage;
