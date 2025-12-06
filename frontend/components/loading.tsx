import React from "react";
import { Loader2 } from "lucide-react";

const Loading = () => {
  return (
    <div className="flex flex-col items-center justify-center h-screen bg-white animate-fadeIn">
      {/* Lucide Spinner */}
      <Loader2 className="w-14 h-14 text-black animate-spin" />
    </div>
  );
};

export default Loading;
