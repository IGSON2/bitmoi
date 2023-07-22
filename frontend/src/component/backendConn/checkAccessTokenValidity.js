import axiosClient from "./axiosClient";

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
    console.log(refResponse.data);
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
