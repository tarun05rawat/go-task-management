"use client";

import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Inter } from "next/font/google";

const inter = Inter({
  subsets: ["latin"],
  display: "swap",
});

export default function Home() {
  const router = useRouter();

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
            onClick={() => router.push("/auth/login")}
            variant="outline"
            className="bg-transparent border-slate-500 text-white hover:bg-slate-700 hover:text-white transition-all duration-300 px-6 py-4 text-base h-auto cursor-pointer"
          >
            Login
          </Button>

          <Button
            onClick={() => router.push("/auth/signup")}
            className="bg-teal-600 hover:bg-teal-500 text-white transition-all duration-300 px-6 py-4 text-base h-auto"
          >
            Sign Up
          </Button>
        </div>
      </div>

      {/* Glassmorphism decorative element */}
      <div className="absolute bottom-0 left-0 w-full h-32 bg-gradient-to-t from-slate-800/30 to-transparent backdrop-blur-sm" />

      {/* Subtle decorative circles */}
      <div className="absolute top-1/4 right-1/4 w-64 h-64 rounded-full bg-teal-500/5 blur-3xl" />
      <div className="absolute bottom-1/4 left-1/3 w-72 h-72 rounded-full bg-indigo-500/5 blur-3xl" />
    </main>
  );
}
