"use client";

import { createContext, useContext, useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import api from "@/utils/api";
import axios from "axios";
import { toast } from "sonner";

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
      console.log("Sending login request to backend...");
      const res = await api.post("/login", { email, password });
      const token = res.data.token;

      // ✅ Store token
      localStorage.setItem("token", token);
      api.defaults.headers.common["Authorization"] = `Bearer ${token}`;
      setUser({ token });

      toast.success("Login successful! 🎉");
      router.push("/dashboard"); // ✅ Redirect to dashboard
    } catch (error: unknown) {
      console.error("Login failed:", error);

      if (axios.isAxiosError(error) && error.response) {
        const errorMessage =
          error.response.data?.message || "Invalid email or password";
        toast.error(errorMessage);
      } else {
        toast.error("An error occurred. Please try again.");
      }
    }
  };

  // 🔑 Signup Function
  const signup = async (name: string, email: string, password: string) => {
    try {
      console.log("Sending signup request to backend...");
      await api.post("/signup", { name, email, password });

      toast.success("Account created successfully! 🎉 Logging in...");
      await login(email, password); // ✅ Auto-login after signup
    } catch (error: unknown) {
      console.error("Signup failed:", error);

      if (axios.isAxiosError(error) && error.response) {
        const status = error.response.status;
        const errorMessage =
          error.response.data?.message || "Signup failed. Please try again.";

        if (status === 400) {
          toast.error("An account with this email already exists.");
        } else {
          toast.error(errorMessage);
        }
      } else {
        toast.error("An error occurred. Please try again.");
      }
    }
  };

  // 🚪 Logout Function
  const logout = () => {
    localStorage.removeItem("token");
    delete api.defaults.headers.common["Authorization"];
    setUser(null);
    toast.info("You have been logged out.");
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
