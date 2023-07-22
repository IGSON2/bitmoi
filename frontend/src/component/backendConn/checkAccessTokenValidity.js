import axiosClient from "./axiosClient";

// SetUserInfo 를 props로 받아서 token과 함께 전달되는 userinfo 초기화하기
const checkAccessTokenValidity = async () => {
  const accessToken = localStorage.getItem("accessToken");
  const refreshToken = localStorage.getItem("refreshToken");

  if (!accessToken) {
    return false;
  }
  try {
    const response = await axiosClient.post("/verify_token", {
      token: accessToken,
    });
    if (response.status === 200) {
      return true;
    }
  } catch {
    const refResponse = await axiosClient.post("/token/reissue_access", {
      refresh_token: refreshToken,
    });
    if (refResponse.status === 200) {
      localStorage.removeItem("accessToken");
      localStorage.setItem("accessToken", refResponse.data.access_token);
      return true;
    } else {
      return false;
    }
  }
};

export default checkAccessTokenValidity;
