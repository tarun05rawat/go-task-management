"use client";

import { createContext, useContext, useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import api from "@/utils/api";
import axios from "axios";

interface User {
  token: string;
}

interface AuthContextType {
  user: User | null;
  login: (email: string, password: string) => Promise<void>;
  signup: (name: string, email: string, password: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);

  const router = useRouter();

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (token) {
      api.defaults.headers.common["Authorization"] = `Bearer ${token}`;
      setUser({ token });
    }
  }, []);

  // 🔑 Login Function
  const login = async (email: string, password: string) => {
    try {
      const res = await api.post("/login", { email, password });
      const token = res.data.token;
      localStorage.setItem("token", token);
      api.defaults.headers.common["Authorization"] = `Bearer ${token}`;
      setUser({ token });
      router.push("/dashboard"); // Redirect to dashboard
    } catch (error) {
      console.error("Login failed", error);
    }
  };

  // 🔑 Signup Function
  const signup = async (name: string, email: string, password: string) => {
    try {
      console.log("Sending signup request to backend...");

      const res = await api.post("/signup", { name, email, password }); // ✅ Ensure this matches your backend

      console.log("Signup successful:", res.data);

      // Auto-login after signup
      await login(email, password);
    } catch (error) {
      console.error("Signup failed:", error);

      if (axios.isAxiosError(error)) {
        console.error("Axios error response:", error.response);
      }
    }
  };

  // 🚪 Logout Function
  const logout = () => {
    localStorage.removeItem("token");
    delete api.defaults.headers.common["Authorization"];
    setUser(null);
    router.push("/login");
  };

  return (
    <AuthContext.Provider value={{ user, login, signup, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) throw new Error("useAuth must be used within an AuthProvider");
  return context;
};
