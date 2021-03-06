import { useEffect, useState } from "react";
import styles from "./Ranking.module.css";
import RankDiv from "./RankDiv/RankDiv";
import Topbutton from "../../component/Topbutton/Topbutton";
import Footer from "../../component/Footer/Footer";
import H_NavBar from "../../component/navbar/H_NavBar";

function Ranking() {
  const [data, setData] = useState({
    totallist: [
      {
        user: "",
        displayname: "",
        photourl: "",
        scoreid: "",
        balance: 0,
      },
    ],
  });
  const getUserScore = async () => {
    const result = await fetch("http://www.bitmoi.net/api/ranking");
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
        {data.totallist.map((v, i) => {
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

export default Ranking;
