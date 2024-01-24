import H_NavBar from "../../component/navbar/H_NavBar";
import ScoreGraph from "./ScoreGraph/ScoreGraph";
import { useEffect, useState } from "react";
import styles from "./mypage.module.css";
import Header from "./Header/Header";
import { BsCaretLeftFill, BsCaretRightFill } from "react-icons/bs";
import axiosClient from "../../component/backendConn/axiosClient";
import checkAccessTokenValidity from "../../component/backendConn/checkAccessTokenValidity";

function MyPage() {
  const [index, setIndex] = useState(1);
  const [data, setData] = useState([{}]);
  const getUserScore = async (i) => {
    try {
      const response = await axiosClient.get(
        `/myscore?mode=competition&page=${index}`
      );
      if (response.status === 200) {
        setData(response.data);
      }
    } catch (error) {
      setData(null);
      console.error(error);
    }
  };

  const increaseIdx = () => {
    setIndex((current) => current + 1);
  };
  const decreaseIdx = () => {
    if (index <= 1) {
      return;
    }
    setIndex((current) => current - 1);
  };

  useEffect(() => {
    const verifyToken = async () => {
      const userInfo = await checkAccessTokenValidity();
      if (!userInfo) {
        alert("로그인이 필요합니다.");
        window.location.replace("/login");
        return;
      }
    };

    verifyToken();
  }, []);

  useEffect(() => {
    getUserScore(index);
  }, [index]);

  return (
    <div className={styles.scorediv}>
      <div className={styles.navbar}>
        <H_NavBar />
      </div>
      <div className={styles.titlediv}>
        <h1 className={styles.title}>YOUR RECORD</h1>
      </div>
      <div className={styles.graphbody}>
        <Header />
        {data.length >= 1 ? (
          data.map((v, i) => {
            return (
              <ScoreGraph key={i} index={i + 1 + 15 * (index - 1)} obj={v} />
            );
          })
        ) : (
          <h3 className={styles.norecord}>아직 기록이 없어요!</h3>
        )}
      </div>
      <div className={styles.indexnav}>
        <div className={styles.indexbtn} onClick={decreaseIdx}>
          <BsCaretLeftFill />
        </div>
        <div className={styles.indexnum}>{index}</div>
        <div className={styles.indexbtn} onClick={increaseIdx}>
          <BsCaretRightFill />
        </div>
      </div>
      <div className={styles.footer}>
        <div>Copyright &copy; 2023 IGSON All rights reserved.</div>
      </div>
    </div>
  );
}

export default MyPage;
