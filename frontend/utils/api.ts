import axios from "axios";

// ✅ Ensure Axios is sending credentials
const api = axios.create({
  baseURL: "http://localhost:8080", // ✅ Make sure this is correct
  withCredentials: true, // ✅ Allow credentials (JWT tokens/cookies)
  headers: {
    "Content-Type": "application/json",
  },
});

export default api;
