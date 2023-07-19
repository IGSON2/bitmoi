import axiosClient from "./axiosClient";

const checkAccessTokenValidity = async () => {
  const accessToken = localStorage.getItem("accessToken");

  if (!accessToken) {
    return false;
  }

  try {
    const response = await axiosClient.post("/verify_token", {
      token: accessToken,
    });

    if (response.status === 200) {
      return true;
    } else {
      return false;
    }
  } catch (error) {
    console.error("Error checking token validity:", error);
    return false;
  }
};

export default checkAccessTokenValidity;
