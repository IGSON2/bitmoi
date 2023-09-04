import axios from "axios";

const axiosFormClient = axios.create({
  // baseURL: "http://localhost:5000",
  baseURL: "https://api.bitmoi.co.kr",
  headers: {
    "Content-Type": "multipart/form-data",
  },
});

axiosFormClient.interceptors.request.use(
  (config) => {
    const accessToken = localStorage.getItem("accessToken");
    if (accessToken) {
      config.headers["Authorization"] = `Bearer ${accessToken}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

export default axiosFormClient;
