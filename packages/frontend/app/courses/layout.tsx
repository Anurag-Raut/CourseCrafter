import { ReactNode } from "react";
import Sidebar from "./_components/sidebar";
import Navbar from "./_components/navbar";
import { get } from "@/api";
import { redirect } from "next/navigation";
import { headers } from "next/headers";
import { getSession } from "@auth0/nextjs-auth0";

export default async function Layout({ children }: { children: ReactNode }) {
  const { user } = await getSession();


  if (!user) {
    redirect("/");
  }
  return (
    <div>
      <Navbar />
      <div className=" h-[calc(100vh-89px)]  w-full flex">
        <Sidebar />
        <div className=" w-full h-full bg-secondaryBackground">{children}</div>
      </div>
    </div>
  );
}
