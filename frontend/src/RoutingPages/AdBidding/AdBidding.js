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
import HorizontalLine from "../../component/lines/HorizontalLine";
import VerticalLine from "../../component/lines/VerticalLine";

function AddBidding() {
  const [idx, setIdx] = useState(0);
  const titles = ["연습모드 하단", "랭크 페이지 중간", "무료 토큰 지급 페이지"];
  const previImages = [practice, rank];
  const pathes = ["practice", "rank"];
  const [userID, setUserID] = useState("");
  const [bidAmt, setBidAmt] = useState(0);
  const [nextUnlock, setNextUnlock] = useState();

  const highestBidder = async (path) => {
    try {
      const res = await axiosClient.get(`/highestBidder?location=${path}`);
      setUserID(res.data.user_id);
      setBidAmt(res.data.amount);
    } catch (error) {
      console.error(error);
      setUserID("아직 입찰자가 없습니다.");
      setBidAmt(0);
    }
  };

  const clickPrevious = () => {
    if (idx <= 0) {
      return;
    }
    setIdx((current) => current - 1);
    const path = pathes[idx - 1];
    highestBidder(path);
  };

  const clickNext = () => {
    if (idx >= titles.length) {
      return;
    }
    setIdx((current) => current + 1);
    const path = pathes[idx + 1];
    highestBidder(path);
  };

  useEffect(() => {
    const getNextBidUnlock = async () => {
      const res = await axiosClient.get("/nextBidUnlock");
      setNextUnlock(res.data.next_unlock);
    };

    getNextBidUnlock();
    highestBidder(pathes[0]);
  }, []);

  return (
    <div className={styles.adbidding}>
      <div className={styles.navbar}>
        <H_NavBar />
      </div>
      <div className={styles.title}>
        {(idx + 1).toString().padStart(2, "0")}
        {". "}
        {titles[idx]}
      </div>
      <div className={styles.preview}>
        <img
          className={styles.navbutton}
          src={previous}
          onClick={clickPrevious}
        />
        <img className={styles.previewimage} src={previImages[idx]} />
        <div className={styles.highestbidder}>
          <h2>최고 입찰자</h2>
          <HorizontalLine />
          <h3>{userID}</h3>
          <div className={styles.tokenbalance}>
            <img src={symbol} />
            <h3>{bidAmt.toLocaleString()}</h3>
          </div>
        </div>
        <img className={styles.navbutton} src={next} onClick={clickNext} />
      </div>
      <HorizontalLine />
      <div className={styles.timer}>
        <h2>입찰 마감까지</h2>
        {nextUnlock ? <Countdown nextUnlock={nextUnlock} /> : null}
      </div>
    </div>
  );
}

export default AddBidding;
