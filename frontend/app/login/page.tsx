"use client";
import React from "react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useAppData, user_service } from "@/context/AppContext";
import Cookies from "js-cookie";
import toast from "react-hot-toast";
import { useGoogleLogin } from "@react-oauth/google";
import axios from "axios";
import Loading from "@/components/loading";
import { redirect } from "next/navigation";
const LoginPage = () => {
  const { isAuth, user, setIsAuth, loading, setLoading, setUser } =
    useAppData();

  if (isAuth) return redirect("/");
  const responseGoogle = async (authResult: any) => {
    setLoading(true);
    try {
      const result: any = await axios.post(`${user_service}/api/v1/login`, {
        code: authResult["code"],
      });

      Cookies.set("token", result.data.token, {
        expires: 5,
        secure: true,
        path: "/",
      });

      toast.success(result.data.message);
      setIsAuth(true);
      setLoading(false);
      setUser(result.data.user);
    } catch (error) {
      toast.error("Login failed. Please try again.");
      console.error("Login error:", error);
    }
  };

  const googleLogin = useGoogleLogin({
    onSuccess: responseGoogle,
    onError: responseGoogle,
    flow: "auth-code",
  });

  return (
    <>
      {loading ? (
        <Loading />
      ) : (
        <div className="w-[350px] m-auto mt-[200px]">
          <Card className="w-[350px]">
            <CardHeader>
              <CardTitle>Login to Postly</CardTitle>
              <CardDescription>your go-to blogging platform</CardDescription>
            </CardHeader>
            <CardContent>
              <Button className="w-full" onClick={googleLogin}>
                <img
                  src={"/google.png"}
                  className="w-8 h-8"
                  alt="google icon"
                />
                Login with Google{" "}
              </Button>
            </CardContent>
          </Card>
        </div>
      )}
    </>
  );
};

export default LoginPage;
