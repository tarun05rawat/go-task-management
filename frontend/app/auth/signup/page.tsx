"use client";

import { useState } from "react";
import { useAuth } from "@/context/AuthContext";

export default function Signup() {
  const { signup } = useAuth();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleSignup = async (e: React.FormEvent<HTMLFormElement>) => {
    // ✅ Fix TypeScript error
    e.preventDefault();
    await signup(name, email, password);
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-900 text-white">
      <h1 className="text-3xl font-bold">Sign Up</h1>
      <form onSubmit={handleSignup} className="flex flex-col gap-4 mt-6">
        <input
          type="text"
          placeholder="Full Name"
          className="p-2 bg-gray-800 rounded-md text-white"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
        />
        <input
          type="email"
          placeholder="Email"
          className="p-2 bg-gray-800 rounded-md text-white"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />
        <input
          type="password"
          placeholder="Password"
          className="p-2 bg-gray-800 rounded-md text-white"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        <button
          type="submit"
          className="bg-green-600 hover:bg-green-500 text-white px-6 py-2 rounded-md"
        >
          Sign Up
        </button>
      </form>
    </div>
  );
}
