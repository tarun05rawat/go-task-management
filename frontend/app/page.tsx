"use client";

import { useRouter } from "next/navigation"; // ✅ Import useRouter
import { Button } from "@/components/ui/button";
import { Inter } from "next/font/google";

const inter = Inter({
  subsets: ["latin"],
  display: "swap",
});

export default function Home() {
  const router = useRouter(); // ✅ Initialize router

  return (
    <main className="flex min-h-screen flex-col items-center justify-center bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 text-white p-4">
      <div className="w-full max-w-md mx-auto text-center space-y-8">
        <div className="space-y-2">
          <h1
            className={`text-5xl md:text-6xl font-bold tracking-tight ${inter.className}`}
          >
            Ta-da List
          </h1>
          <p className="text-xl text-slate-300">Organize. Simplify. Achieve.</p>
        </div>

        <div className="flex flex-col sm:flex-row gap-4 justify-center mt-8">
          <Button
            onClick={() => router.push("/auth/login")} // ✅ Fix navigation
            className="bg-transparent border border-slate-500 text-white px-6 py-4 rounded-md hover:bg-slate-700 transition-all duration-300"
          >
            Login
          </Button>

          <Button
            onClick={() => router.push("/auth/signup")} // ✅ Fix navigation
            className="bg-teal-600 hover:bg-teal-500 text-white px-6 py-4 rounded-md transition-all duration-300"
          >
            Sign Up
          </Button>
        </div>
      </div>
    </main>
  );
}
