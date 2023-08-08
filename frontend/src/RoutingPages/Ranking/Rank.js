import { useEffect, useState } from "react";
import styles from "./Rank.module.css";
import RankDiv from "./RankDiv/RankDiv";
import Topbutton from "../../component/Topbutton/Topbutton";
import H_NavBar from "../../component/navbar/H_NavBar";
import axiosClient from "../../component/backendConn/axiosClient";

function Rank() {
  const [pageNum, setPageNum] = useState(1);

  const [data, setData] = useState([{}]);
  const getUserScore = async () => {
    const response = await axiosClient.get(`/rank/${pageNum}`);
    setData(response.data);
  };

  useEffect(() => {
    getUserScore();
  }, [pageNum]);
  return (
    <div className={styles.scorediv}>
      <div className={styles.navbar}>
        <H_NavBar />
      </div>
      <div className={styles.graphbody}>
        <div className={styles.titlediv}>
          <h1 className={styles.title}>RANKING BOARD</h1>
        </div>
        {data.map((v, i) => {
          return <RankDiv key={i} index={i + 1} obj={v} />;
        })}
        <div className={styles.footer}>
          <div>Copyright &copy; 2023 IGSON All rights reserved.</div>
        </div>
      </div>
      <Topbutton />
    </div>
  );
}

export default Rank;
