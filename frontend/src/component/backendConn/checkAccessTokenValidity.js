import axiosClient from "./axiosClient";

const checkAccessTokenValidity = async () => {
  const accessToken = localStorage.getItem("accessToken");
  const refreshToken = localStorage.getItem("refreshToken");

  if (!accessToken) {
    return null;
  }
  try {
    const response = await axiosClient.post("/verify_token", {
      token: accessToken,
    });
    if (response.status === 200) {
      return response.data;
    } else {
      throw response.data;
    }
  } catch (error) {
    const refResponse = await axiosClient.post("/token/reissue_access", {
      refresh_token: refreshToken,
    });
    if (refResponse.status === 200) {
      localStorage.removeItem("accessToken");
      localStorage.setItem("accessToken", refResponse.data.access_token);
      return refResponse.data.user;
    } else {
      console.error(error);
      return null;
    }
  }
};

export default checkAccessTokenValidity;
