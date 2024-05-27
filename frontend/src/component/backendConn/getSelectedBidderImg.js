import axiosClient from "./axiosClient";

const getSelectedBidderImg = async (location) => {
  try {
    const response = await axiosClient.get(
      `/basic/selectedBidder?location=${location}`
    );
    if (response.status === 200) {
      return `https://cdn.bitmoi.co.kr/bidding/${location}/${response.data.user_id}`;
    } else {
      throw response.data;
    }
  } catch (error) {
    console.error("Get selected bidder error. err:", error);
    return "";
  }
};

export default getSelectedBidderImg;
