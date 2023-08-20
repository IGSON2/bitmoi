import H_NavBar from "../../component/navbar/H_NavBar";
import styles from "./AdBidding.module.css";
import practice from "../../component/images/preview_practice.png";
import rank from "../../component/images/preview_rank.png";
import previous from "../../component/images/previous.png";
import next from "../../component/images/next.png";
import { useState } from "react";

function AddBidding() {
  const locations = [
    "연습모드 하단",
    "랭크 페이지 중간",
    "무료 토큰 지급 페이지",
  ];
  const previews = [practice, rank];
  const [bidAmt, setBidAmt] = useState(0);

  return (
    <div className={styles.adbidding}>
      <div className={styles.navbar}>
        <H_NavBar />
      </div>
      {locations.map((loc, idx) => {
        return (
          <div>
            <div className={styles.title}>{loc}</div>
            <div className={styles.preview}>
              <img src={previous} />
              <img src={previews[idx]} />
              <div className={styles.highestbidder}>
                <div>최고 입찰자</div>
                <div>그냥적당히기다란이름</div>
                <div>2308272423</div>
              </div>
              <img src={next} />
            </div>
            <div className={styles.timer}></div>
          </div>
        );
      })}
    </div>
  );
}

export default AddBidding;
