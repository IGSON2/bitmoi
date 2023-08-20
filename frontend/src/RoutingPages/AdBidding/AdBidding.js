import axiosClient from "../../component/backendConn/axiosClient";
import H_NavBar from "../../component/navbar/H_NavBar";
import styles from "./AdBidding.module.css";
import practice from "../../component/images/preview_practice.png";
import rank from "../../component/images/preview_rank.png";
import symbol from "../../component/images/logo.png";
import previous from "../../component/images/previous.png";
import next from "../../component/images/next.png";
import { useEffect, useState } from "react";
import Countdown from "./Countdown/Countdown";

function AddBidding() {
  const [idx, setIdx] = useState(0);
  const titles = ["연습모드 하단", "랭크 페이지 중간", "무료 토큰 지급 페이지"];
  const locations = [practice, rank];
  const [userID, setUserID] = useState("");
  const [bidAmt, setBidAmt] = useState(0);
  const [nextUnlock, setNextUnlock] = useState();

  const highestBidder = async (loc) => {
    const res = await axiosClient.get(`/highestBidder?location=${loc}`);
    setUserID(res.data.user_id);
    setBidAmt(res.data.amount);
  };

  const clickPrevious = () => {
    setIdx((current) => current - 1);
    const loc = locations[idx - 1];
    highestBidder(loc);
  };

  const clickNext = () => {
    setIdx((current) => current + 1);
    const loc = locations[idx + 1];
    highestBidder(loc);
  };

  useEffect(() => {
    const getNextBidUnlock = async () => {
      const res = await axiosClient.get("/nextBidUnlock");
      setNextUnlock(res.data);
    };

    getNextBidUnlock();
  }, []);

  return (
    <div className={styles.adbidding}>
      <div className={styles.navbar}>
        <H_NavBar />
      </div>
      <div className={styles.title}>{titles[idx]}</div>
      <div className={styles.preview}>
        <img
          className={styles.navbutton}
          src={previous}
          onClick={clickPrevious}
        />
        <img className={styles.previewimage} src={locations[idx]} />
        <div className={styles.highestbidder}>
          <div>최고 입찰자</div>
          <div>{userID}</div>
          <div className={styles.tokenbalance}>
            <img src={symbol} />
            <h3>{bidAmt.toLocaleString()}</h3>
          </div>
        </div>
        <img className={styles.navbutton} src={next} onClick={clickNext} />
      </div>
      <div className={styles.timer}>
        <div>입찰 마감까지</div>
        {nextUnlock ? <Countdown nextUnlock={nextUnlock} /> : null}
      </div>
    </div>
  );
}

export default AddBidding;
