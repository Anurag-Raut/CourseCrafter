"use client";
import React from "react";
import CheckIcon from "./checkIcon";
import { Badge } from "@/components/ui/badge";
import { ClockIcon } from "lucide-react";
interface CoursecardProps {
  topic: string;
  status: boolean;
  username: string;
}
const Coursecard = ({ topic, status, username }: CoursecardProps) => {
  const getDaysAgo = (date: Date): number => {
    const today = new Date();
    const diffTime = Math.abs(today.getTime() - date.getTime());
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    return diffDays;
  };
  return (
    <div className="border p-4 flex items-center justify-between rounded-md shadow-md">
      <div className="w-16 h-16 m-4 bg-zinc-200 flex items-center justify-center">
        <img src="https://placehold.co/64" alt="Course Thumbnail" />
      </div>
      <div>
        <h3 className="text-lg font-semibold">{topic}</h3>
        <p className="text-sm text-gray-500">PowerPoint, Excel, Tableau</p>
        <p className="text-xs text-gray-400 mt-1">{username} - 18d ago</p>
        <div className="flex items-center mt-2">
          <Badge className="text-black" variant="secondary">
            {status == true ? "Converted" : "Pending"}
          </Badge>

          {status == true ? (
            <CheckIcon className="ml-2 h-4 w-4 text-green-500" />
          ) : (
            <ClockIcon className="ml-2 h-4 w-4 text-yellow-500" />
          )}
        </div>
      </div>
      <div className="ml-auto bg-gray-100 py-1 px-3 rounded-full">7/12</div>
    </div>
  );
};

export default Coursecard;