import axiosClient from "./axiosClient";

const checkAccessTokenValidity = async () => {
  const accessToken = localStorage.getItem("accessToken");
  const refreshToken = localStorage.getItem("refreshToken");

  if (!accessToken) {
    return null;
  }
  try {
    const response = await axiosClient.post("/verifyToken", {
      token: accessToken,
    });
    if (response.status === 200) {
      return response.data;
    } else {
      throw response.data;
    }
  } catch (error) {
    try {
      const refResponse = await axiosClient.post("/reissueAccess", {
        refresh_token: refreshToken,
      });
      if (refResponse.status === 200) {
        localStorage.removeItem("accessToken");
        localStorage.setItem("accessToken", refResponse.data.access_token);
        console.log("access token updated by refresh token.");
        return refResponse.data.user;
      } else {
        return null;
      }
    } catch (reissueError) {
      console.error("Error while reissuing access token:", reissueError);
      return null;
    }
  }
};

export default checkAccessTokenValidity;
