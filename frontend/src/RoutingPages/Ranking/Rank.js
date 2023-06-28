import { useEffect, useState } from "react";
import styles from "./Rank.module.css";
import RankDiv from "./RankDiv/RankDiv";
import Topbutton from "../../component/Topbutton/Topbutton";
import Footer from "../../component/Footer/Footer";
import H_NavBar from "../../component/navbar/H_NavBar";

function Rank() {
  const [data, setData] = useState({
    rankingBoard: [
      {
        user_id: "",
        photo_url: "",
        display_name: "",
        score_id: "",
        final_balance: 0,
        comment: "",
      },
    ],
  });
  const getUserScore = async () => {
    const result = await fetch("http://bitmoi.co.kr:5000/rank");
    const json = await result.json();
    setData(json);
  };

  useEffect(() => {
    getUserScore();
  }, []);
  return (
    <div className={styles.scorediv}>
      <div className={styles.navbar}>
        <H_NavBar />
      </div>
      <div className={styles.graphbody}>
        <div className={styles.titlediv}>
          <h1 className={styles.title}>RANKING BOARD</h1>
        </div>
        {data.rankingBoard.map((v, i) => {
          return <RankDiv key={i} index={i + 1} obj={v} />;
        })}
        <div className={styles.footer}>
          <div>Copyright &copy; 2022 IGSON All rights reserved.</div>
        </div>
      </div>
      <Topbutton />
    </div>
  );
}

export default Rank;
