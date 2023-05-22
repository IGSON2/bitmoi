import { getAuth, onAuthStateChanged } from "firebase/auth";
import H_NavBar from "../../component/navbar/H_NavBar";
import ScoreGraph from "./ScoreGraph/ScoreGraph";
import { useEffect, useState } from "react";
import styles from "./myscore.module.css";
import Header from "./Header/Header";
import { BsCaretLeftFill, BsCaretRightFill } from "react-icons/bs";

function UserScore() {
  const auth = getAuth();
  const [index, setIndex] = useState(1);
  const [userLoaded, setUserLoaded] = useState(false);
  const [data, setData] = useState({
    scorelist: [
      {
        endprice: 0,
        entryprice: 0,
        entrytime: "",
        leverage: 0,
        outtime: 0,
        pairname: "",
        pnl: 0,
        position: "",
        roe: 0,
        scoreid: "",
        stage: 0,
        user: "",
      },
    ],
  });
  const getUserScore = async (i) => {
    // TODO: update userid for firebase
    const result = await fetch(
      "http://www.bitmoi.net/api/myscore/?user=" +
        auth.currentUser.uid +
        "&index=" +
        `${i}` +
        "&scoreid="
    );
    const json = await result.json();
    setData(json);
  };

  useEffect(() => {
    onAuthStateChanged(auth, (user) => {
      if (user) {
        getUserScore(index);
        setUserLoaded(true);
      } else {
        setData({});
      }
    });
  }, []);

  const increaseIdx = () => {
    setIndex((current) => current + 1);
  };
  const decreaseIdx = () => {
    setIndex((current) => current - 1);
  };

  useEffect(() => {
    if (userLoaded) {
      getUserScore(index);
    }
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
        {data.scorelist ? (
          data.scorelist.map((v, i) => {
            return (
              <ScoreGraph key={i} index={i + 1 + 15 * (index - 1)} obj={v} />
            );
          })
        ) : (
          <h3>아직 기록이 없어요!</h3>
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
        <div>Copyright &copy; 2022 IGSON All rights reserved.</div>
      </div>
    </div>
  );
}

export default UserScore;
