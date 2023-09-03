import axios from "axios";

const axiosClient = axios.create({
  baseURL: "https://api.bitmoi.co.kr",
  headers: {
    "Content-Type": "application/json",
  },
});

const getSelectedBidderImg = async (location) => {
  try {
    const response = await axiosClient.get(
      `/selectedBidder?location=${location}`
    );
    if (response.status === 200) {
      return `https://cnd.bitmoi.co.kr/bidding/${location}/${response.data}`;
    } else {
      throw response.data;
    }
  } catch (error) {
    console.error("Get selected bidder error. err:", error);
    return "";
  }
};

export default getSelectedBidderImg;
